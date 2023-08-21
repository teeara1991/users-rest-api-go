package handlers

import (
	"rest-api-go/internal/handlers/user"
	"rest-api-go/internal/service"
	"rest-api-go/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(router *httprouter.Router, service *service.Service, logger *logging.Logger) {
	//register handlers here
	handler := user.NewUserHandler(logger, service.UserService)
	handler.Register(router)

}
