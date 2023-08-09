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
}

type Advert interface {
	CreateAdvert(ctx context.Context, advert models.Advert) (string, error)
	DeleteAdvert(ctx context.Context, id string) error
}

type Services interface {
	User
	Advert
}

type Service struct {
	User
	Advert
}

func NewService(storage storage.Storage, logger logger.Logger) *Service {
	return &Service{
		NewUserService(storage, logger),
		NewAdvertService(storage, logger),
	}
}
