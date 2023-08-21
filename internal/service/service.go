package service

import (
	"context"
	"rest-api-go/internal/entities/user"
)

//abstraction of the service layer

type UserService interface {
	Create(ctx context.Context, dto user.CreateUserDTO) (userUUID string, err error)
	FindOne(ctx context.Context, id string) (user.User, error)
	FindAll(ctx context.Context) ([]user.User, error)
	Update(ctx context.Context, dto user.UpdateUserDTO) error
	Delete(ctx context.Context, id string) error
}

type Service struct {
	UserService UserService
}
