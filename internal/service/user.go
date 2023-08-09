package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/logger"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/romandnk/advertisement/internal/storage"
	"net/mail"
	"time"
)

var secret = []byte("secret")

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

func (u *UserService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := u.user.GetUserByEmail(ctx, email)
	if err != nil {
		var customError custom_error.CustomError
		if errors.As(err, &customError) {
			return "", err
		}
		return "", custom_error.CustomError{Field: "", Message: err.Error()}
	}

	if user.Deleted == true {
		return "", nil
	}

	if !comparePassword(password, user.Password) {
		return "", custom_error.CustomError{Field: "password", Message: "invalid password"}
	}

	token, err := createJWT(secret, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
