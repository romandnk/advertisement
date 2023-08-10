package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/logger"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/romandnk/advertisement/internal/storage"
	"go.uber.org/zap"
	"strings"
	"time"
)

var (
	ErrAdvertServiceEmptyTitle    = errors.New("empty title")
	ErrAdvertServiceNegativePrice = errors.New("negative price")
	ErrAdvertServiceNoImages      = errors.New("no images")
	ErrAdvertServiceManyImages    = errors.New("max number of images is 7")
	ErrAdvertServiceNoUserID      = errors.New("no user id")
)

type AdvertService struct {
	advert       storage.AdvertStorage
	logger       logger.Logger
	pathToImages string
}

func NewAdvertService(advert storage.AdvertStorage, logger logger.Logger, pathToImages string) *AdvertService {
	return &AdvertService{
		advert:       advert,
		logger:       logger,
		pathToImages: pathToImages,
	}
}

func (a *AdvertService) CreateAdvert(ctx context.Context, advert models.Advert) (string, error) {
	advert.ID = uuid.New().String()

	advert.Title = strings.TrimSpace(advert.Title)
	if advert.Title == "" {
		return "", custom_error.CustomError{Field: "title", Message: ErrAdvertServiceEmptyTitle.Error()}
	}

	advert.Description = strings.TrimSpace(advert.Description)

	if advert.Price.IsNegative() {
		return "", custom_error.CustomError{Field: "price", Message: ErrAdvertServiceNegativePrice.Error()}
	}

	now := time.Now()
	advert.CreatedAt = now
	advert.UpdatedAt = now
	if len(advert.Images) == 0 {
		return "", custom_error.CustomError{Field: "images", Message: ErrAdvertServiceNoImages.Error()}
	}
	if len(advert.Images) > 7 {
		return "", custom_error.CustomError{Field: "images", Message: ErrAdvertServiceManyImages.Error()}
	}

	for _, image := range advert.Images {
		image.ID = uuid.New().String()
		image.AdvertID = advert.ID
		image.CreatedAt = now
		err := saveImage(image, a.pathToImages)
		if err != nil {
			return "", custom_error.CustomError{Field: "images", Message: err.Error()}
		}
	}

	id, err := a.advert.CreateAdvert(ctx, advert)
	if err != nil {
		for _, image := range advert.Images {
			err := deleteImage(image.ID, a.pathToImages)
			if err != nil {
				a.logger.Error("error deleting image while creating advert", zap.String("error", err.Error()))
			}
		}
		return "", err
	}

	return id, nil
}

func (a *AdvertService) DeleteAdvert(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return custom_error.CustomError{Field: "id", Message: err.Error()}
	}

	var userID string
	switch ctx.Value("user_id").(type) {
	case string:
		userID = ctx.Value("user_id").(string)
	default:
		return custom_error.CustomError{Field: "user_id", Message: ErrAdvertServiceNoUserID.Error()}
	}

	imageIDs, err := a.advert.DeleteAdvert(ctx, parsedID.String(), userID)
	if err != nil {
		return err
	}

	for _, imageID := range imageIDs {
		err := deleteImage(imageID, a.pathToImages)
		if err != nil {
			a.logger.Error("error deleting image while deleting advert", zap.String("error", err.Error()))
		}
	}

	return nil
}

func (a *AdvertService) GetAdvertByID(ctx context.Context, id string) (models.Advert, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Advert{}, custom_error.CustomError{Field: "id", Message: err.Error()}
	}
	return a.advert.GetAdvertByID(ctx, parsedID.String())
}
