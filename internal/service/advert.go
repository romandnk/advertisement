package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/romandnk/advertisement/internal/storage"
	"os"
	"strings"
	"time"
)

var pathToImages = "static/images/"

type AdvertService struct {
	advert storage.AdvertStorage
}

func NewAdvertService(advert storage.AdvertStorage) *AdvertService {
	return &AdvertService{advert: advert}
}

func (a *AdvertService) CreateAdvert(ctx context.Context, advert models.Advert) (string, error) {
	advert.ID = uuid.New().String()

	advert.Title = strings.TrimSpace(advert.Title)
	if advert.Title == "" {
		return "", custom_error.CustomError{Field: "title", Message: "empty title"}
	}

	advert.Description = strings.TrimSpace(advert.Description)

	if advert.Price.IsNegative() {
		return "", custom_error.CustomError{Field: "price", Message: "negative price"}
	}

	now := time.Now()
	advert.CreatedAt = now
	advert.UpdatedAt = now
	if len(advert.Images) == 0 {
		return "", custom_error.CustomError{Field: "images", Message: "no images"}
	}
	if len(advert.Images) > 7 {
		return "", custom_error.CustomError{Field: "images", Message: "max number of images is 7"}
	}

	//TODO: authorization
	advert.UserID = uuid.New().String()

	for _, image := range advert.Images {
		image.ID = uuid.New().String()
		image.CreatedAt = now
		err := saveImage(image, pathToImages)
		if err != nil {
			return "", custom_error.CustomError{Field: "images", Message: err.Error()}
		}
	}

	id, err := a.advert.CreateAdvert(ctx, advert)
	if err != nil {
		for _, image := range advert.Images {
			err := deleteImage(image.ID, pathToImages)
			if err != nil {
				return "", err
			}
		}
		return "", err
	}

	return id, nil
}

func saveImage(image *models.Image, path string) error {
	err := os.WriteFile(path+image.ID+".jpg", image.Data, 0o644)
	return err
}

func (a *AdvertService) DeleteAdvert(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return custom_error.CustomError{Field: "id", Message: err.Error()}
	}

	imageIDs, err := a.advert.DeleteAdvert(ctx, parsedID.String())
	if err != nil {
		return err
	}

	for _, imageID := range imageIDs {
		err := deleteImage(imageID, pathToImages)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteImage(imageID string, path string) error {
	return os.Remove(path + imageID + ".jpg")
}
