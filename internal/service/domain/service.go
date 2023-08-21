package domain

import (
	"rest-api-go/internal/service"
	"rest-api-go/internal/service/domain/user"
	"rest-api-go/internal/storage"
	"rest-api-go/pkg/logging"
)

// all implementations of service in one
func NewService(
	repositories *storage.Repository,
	logger *logging.Logger,
) *service.Service {
	return &service.Service{
		UserService: user.NewUserService(logger, repositories.User),
		//add other services here
	}
}
