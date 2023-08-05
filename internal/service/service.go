package service

import "context"

type Image interface {
	UploadImage(ctx context.Context)
	GetImage(ctx context.Context)
}

type Services interface {
	Image
}

type Service struct {
	image Image
}

func NewService(image Image) *Service {
	return &Service{
		image: image,
	}
}
