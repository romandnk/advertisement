package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/logger"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/romandnk/advertisement/internal/storage"
	"net/mail"
	"time"
)

type UserService struct {
	user   storage.UserStorage
	logger logger.Logger
}

func NewUserService(user storage.UserStorage, logger logger.Logger) *UserService {
	return &UserService{
		user:   user,
		logger: logger,
	}
}

func (u *UserService) SignUp(ctx context.Context, user models.User) (string, error) {
	user.ID = uuid.New().String()

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return "", custom_error.CustomError{Field: "email", Message: err.Error()}
	}

	if err := validatePassword(user.Password); err != nil {
		return "", err
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return "", custom_error.CustomError{Field: "password", Message: err.Error()}
	}

	user.Password = hashedPassword

	return u.user.CreateUser(ctx, user)
}
