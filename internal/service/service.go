package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"context"
	"github.com/romandnk/advertisement/internal/logger"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/romandnk/advertisement/internal/storage"
)

type User interface {
	SignUp(ctx context.Context, user models.User) (string, error)
	SignIn(ctx context.Context, email, password string) (string, error)
}

type Advert interface {
	CreateAdvert(ctx context.Context, advert models.Advert) (string, error)
	DeleteAdvert(ctx context.Context, id string) error
	GetAdvertByID(ctx context.Context, id string) (models.Advert, error)
}

type Image interface {
	GetImageByID(ctx context.Context, id string) (models.Image, error)
}

type Services interface {
	User
	Advert
	Image
}

type Service struct {
	User
	Advert
	Image
}

func NewService(storage storage.Storage, logger logger.Logger, secretKey, pathToImages string) *Service {
	return &Service{
		NewUserService(storage, logger, secretKey),
		NewAdvertService(storage, logger, pathToImages),
		NewImageService(storage, logger, pathToImages),
	}
}
