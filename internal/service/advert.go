package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
	"strings"
	"time"
)

func (s *Service) CreateAdvert(ctx context.Context, advert models.Advert, path string) (string, error) {
	id := uuid.New().String()
	advert.ID = id

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

	//TODO: authorization
	advert.UserID = uuid.New().String()

	return s.Advert.CreateAdvert(ctx, advert, path)
}

func (s *Service) DeleteAdvert(ctx context.Context, id string, path string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return custom_error.CustomError{Field: "id", Message: "invalid id"}
	}
	return s.Advert.DeleteAdvert(ctx, parsedID.String(), path)
}
