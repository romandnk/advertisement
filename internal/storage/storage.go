package storage

import (
	"context"
	"github.com/romandnk/advertisement/internal/models"
)

type ImageStorage interface {
	GetImageByID(ctx context.Context, id string) (models.Image, error)
}

type UserStorage interface {
	CreateUser(ctx context.Context, user models.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type AdvertStorage interface {
	CreateAdvert(ctx context.Context, advert models.Advert) (string, error)
	GetAdvertByID(ctx context.Context, id string) (models.Advert, error)
	DeleteAdvert(ctx context.Context, advertID, userID string) ([]string, error)
}

type Storage interface {
	AdvertStorage
	UserStorage
	ImageStorage
}
