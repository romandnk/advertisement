package postgres

import (
	"context"
	"fmt"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
)

func (s *PostgresStorage) CreateAdvert(ctx context.Context, advert models.Advert) (string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return "", custom_error.CustomError{Field: "", Message: err.Error()}
	}

	queryAdvert := fmt.Sprintf(`
				INSERT INTO %s (id, title, description, price, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6)`, avertsTable)

	ct, err := tx.Exec(ctx, queryAdvert,
		advert.ID,
		advert.Title,
		advert.Description,
		advert.Price,
		advert.CreatedAt,
		advert.UpdatedAt,
	)
	if err != nil {
		tx.Rollback(ctx)
		return "", custom_error.CustomError{Field: "", Message: err.Error()}
	}

	if ct.RowsAffected() == 0 {
		tx.Rollback(ctx)
		return "", custom_error.CustomError{Field: "", Message: "advert was not inserted"}
	}

	queryImage := fmt.Sprintf(`
				INSERT INTO %s (id, advert_id, url)
				VALUES ($1, $2, $3)`, imagesTable)

	for _, image := range advert.Images {
		ct, err := tx.Exec(ctx, queryImage, image.ID, advert.ID, image.Url)
		if err != nil {
			tx.Rollback(ctx)
			return "", custom_error.CustomError{Field: "", Message: err.Error()}
		}
		if ct.RowsAffected() == 0 {
			tx.Rollback(ctx)
			return "", custom_error.CustomError{Field: "", Message: "image was not inserted"}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return "", custom_error.CustomError{Field: "", Message: err.Error()}
	}

	return advert.ID, nil
}
