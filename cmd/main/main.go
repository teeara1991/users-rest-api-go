package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest-api-go/internal/config"
	"rest-api-go/internal/handlers"
	service "rest-api-go/internal/service/domain"
	storage "rest-api-go/internal/storage/mongodb"
	"rest-api-go/pkg/client/mongodb"
	"rest-api-go/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()
	cfgMongo := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(context.Background(),
		cfgMongo.Host, cfgMongo.Port, cfgMongo.Username, cfgMongo.Password,
		cfgMongo.Database, cfgMongo.AuthDB)

	if err != nil {
		panic(err)
	}
	storage := storage.NewRepository(mongoDBClient, cfg.MongoDB.Collection, logger)
	services := service.NewService(storage, logger)
	handlers.RegisterHandlers(router, services, logger)
	logger.Info("register handlers")

	run(router, cfg)

}
func run(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("run server")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Debug("socket path: ", socketPath)
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Info("Server is listening unix socket")

	} else {
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Infof("Server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	}
	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}
