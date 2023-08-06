package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
	"strings"
	"time"
)

func (s *Service) CreateAdvert(ctx context.Context, advert models.Advert) (string, error) {
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

	return s.Advert.CreateAdvert(ctx, advert)
}
