package storage

import (
	"context"
	"github.com/romandnk/advertisement/internal/models"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user models.User) (string, error)
}

type AdvertStorage interface {
	CreateAdvert(ctx context.Context, advert models.Advert) (string, error)
	DeleteAdvert(ctx context.Context, id string) ([]string, error)
}

type Storage interface {
	AdvertStorage
	UserStorage
}
