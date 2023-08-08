package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"context"
	"github.com/romandnk/advertisement/internal/logger"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/romandnk/advertisement/internal/storage"
)

type Advert interface {
	CreateAdvert(ctx context.Context, advert models.Advert) (string, error)
	DeleteAdvert(ctx context.Context, id string) error
}

type Services interface {
	Advert
}

type Service struct {
	Advert
}

func NewService(storage storage.Storage, logger logger.Logger) *Service {
	return &Service{
		NewAdvertService(storage, logger),
	}
}
