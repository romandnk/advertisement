package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/logger"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/romandnk/advertisement/internal/storage"
)

var ErrImageServiceImageNotFound = errors.New("image not found")

type ImageService struct {
	image        storage.ImageStorage
	logger       logger.Logger
	pathToImages string
}

func NewImageService(image storage.ImageStorage, logger logger.Logger, pathToImages string) *ImageService {
	return &ImageService{
		image:        image,
		logger:       logger,
		pathToImages: pathToImages,
	}
}

func (i *ImageService) GetImageByID(ctx context.Context, id string) (models.Image, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Image{}, err
	}

	image, err := i.image.GetImageByID(ctx, parsedID.String())
	if err != nil {
		return image, err
	}

	if image.Deleted == true {
		return models.Image{}, custom_error.CustomError{
			Field:   "id",
			Message: ErrImageServiceImageNotFound.Error(),
		}
	}

	data, err := findImageByID(i.pathToImages, parsedID.String())
	if err != nil {
		return models.Image{}, err
	}

	image.Data = data

	return image, nil
}
