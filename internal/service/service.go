package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"context"
	"github.com/romandnk/advertisement/internal/models"
)

//type Image interface {
//	GetImage(ctx context.Context)
//}

type Advert interface {
	CreateAdvert(ctx context.Context, advert models.Advert, path string) (string, error)
	DeleteAdvert(ctx context.Context, id string, path string) error
}

type Services interface {
	Advert
	//Image
}

type Service struct {
	Advert Advert
	//Image  Image
}

func NewService(advert Advert) *Service {
	return &Service{
		Advert: advert,
		//Image:  image,
	}
}
