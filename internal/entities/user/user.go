package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `bson:"_id,omitempty" json:"id"`
	Username     string `bson:"username" json:"username"`
	PasswordHash string `bson:"password" json:"-"`
	Email        string `bson:"email" json:"email"`
}

type CreateUserDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}
type UpdateUserDTO struct {
	ID          string `json:"uuid,omitempty" bson:"_id,omitempty"`
	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	Password    string `json:"password,omitempty" bson:"password,omitempty"`
	Username    string `json:"username,omitempty" bson:"username,omitempty"`
	OldPassword string `json:"old_password,omitempty" bson:"-"`
	NewPassword string `json:"new_password,omitempty" bson:"-"`
}

func NewUser(dto CreateUserDTO) *User {
	return &User{
		Email:    dto.Email,
		Username: dto.Username,
	}
}
func UpdatedUser(dto UpdateUserDTO) *User {
	return &User{
		ID:           dto.ID,
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: dto.Password,
	}
}

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password due to error %w", err)
	}
	return string(hash), nil
}
