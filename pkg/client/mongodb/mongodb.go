package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDb string) (db *mongo.Database, err error) {
	var mongoDBURL string
	var isAuth bool
	if username == "" && password == "" {
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		isAuth = true
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}

	//Connect

	clientOptions := options.Client().ApplyURI(mongoDBURL)
	if isAuth {
		if authDb == "" {
			authDb = database
		}
		clientOptions.SetAuth(options.Credential{
			Username:   username,
			Password:   password,
			AuthSource: authDb,
		})
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.New("MongoDB Connect Error: " + err.Error())
	}
	//Ping
	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.New("MongoDB Ping Error: " + err.Error())
	}
	return client.Database(database), nil

}
