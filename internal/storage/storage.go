package storage

import "github.com/romandnk/advertisement/internal/service"

type Storage interface {
	service.Advert
	//service.Image
}
