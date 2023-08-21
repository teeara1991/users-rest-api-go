package storage

//repository interface abstraction
import (
	"context"
	"rest-api-go/internal/entities/user"
)

type UserRepository interface {
	Create(ctx context.Context, user user.User) (string, error)
	FindOne(ctx context.Context, id string) (user.User, error)
	FindAll(ctx context.Context) ([]user.User, error)
	Update(ctx context.Context, user user.User) error
	Delete(ctx context.Context, id string) error
}

// add other repositories interfaces here
type Repository struct {
	User UserRepository
	//add other repositories here
}
