package mongodb

import (
	"rest-api-go/internal/storage"
	"rest-api-go/internal/storage/mongodb/user"
	"rest-api-go/pkg/logging"

	"go.mongodb.org/mongo-driver/mongo"
)

// NewRepository implementation for storage of all repositories.
func NewRepository(database *mongo.Database, collection string, logger *logging.Logger) *storage.Repository {
	return &storage.Repository{
		User: user.NewUserRepository(database, collection, logger),
		//add other repositories here
	}
}
