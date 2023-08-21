package user

import (
	"context"
	"errors"
	"fmt"
	"rest-api-go/internal/apperrors"
	"rest-api-go/internal/entities/user"
	"rest-api-go/internal/storage"
	"rest-api-go/pkg/logging"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	logger         *logging.Logger
	UserRepository storage.UserRepository
}

func (s *UserService) Create(ctx context.Context, dto user.CreateUserDTO) (userUUID string, err error) {
	s.logger.Debug("check password and repeat password")
	if dto.Password == "" {
		return userUUID, apperrors.BadRequestError("password is empty")
	}
	if dto.Email == "" {
		return userUUID, apperrors.BadRequestError("email is empty")
	}
	if dto.Username == "" {
		return userUUID, apperrors.BadRequestError("username is empty")
	}

	newUser := user.NewUser(dto)

	s.logger.Debug("generate password hash")
	hash, err := user.GeneratePasswordHash(dto.Password)
	if err != nil {
		s.logger.Errorf("failed to create user due to error %v", err)
		return
	}
	newUser.PasswordHash = hash

	userUUID, err = s.UserRepository.Create(ctx, *newUser)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return userUUID, err
		}
		return userUUID, fmt.Errorf("failed to create user. error: %w", err)
	}

	return userUUID, nil
}
func (s *UserService) FindOne(ctx context.Context, uuid string) (user.User, error) {
	user, err := s.UserRepository.FindOne(ctx, uuid)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return user, err
		}
		return user, fmt.Errorf("failed to find user by uuid. error: %w", err)
	}
	return user, nil
}
func (s *UserService) FindAll(ctx context.Context) ([]user.User, error) {
	users, err := s.UserRepository.FindAll(ctx)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return users, err
		}
		return users, fmt.Errorf("failed to find users. error: %w", err)
	}
	return users, nil
}
func (s *UserService) Update(ctx context.Context, dto user.UpdateUserDTO) error {
	s.logger.Debug("compare old and new passwords")
	if dto.OldPassword != dto.NewPassword || dto.NewPassword == "" {
		s.logger.Debug("get user by uuid")
		findedUser, err := s.FindOne(ctx, dto.ID)
		if err != nil {
			return err
		}

		s.logger.Debug("compare hash current password and old password")
		err = bcrypt.CompareHashAndPassword([]byte(findedUser.PasswordHash), []byte(dto.OldPassword))
		if err != nil {
			return apperrors.BadRequestError("old password does not match current password")
		}

		dto.Password = dto.NewPassword
	}

	updatedUser := user.UpdatedUser(dto)
	s.logger.Debug("generate password hash")
	hash, err := user.GeneratePasswordHash(dto.Password)
	if err != nil {
		return fmt.Errorf("failed to generate hash. error %w", err)
	}
	updatedUser.PasswordHash = hash

	s.logger.Printf("update user with uuid: %+v", updatedUser)

	err = s.UserRepository.Update(ctx, *updatedUser)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update user. error: %w", err)
	}
	return nil
}
func (s *UserService) Delete(ctx context.Context, id string) (err error) {
	err = s.UserRepository.Delete(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return err
}

func NewUserService(
	logger *logging.Logger,
	UserRepository storage.UserRepository,
) *UserService {
	return &UserService{
		logger:         logger,
		UserRepository: UserRepository,
	}
}
