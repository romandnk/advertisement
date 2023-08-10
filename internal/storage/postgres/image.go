package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
)

var ErrImageNotFound = errors.New("image not found")

func (s *PostgresStorage) GetImageByID(ctx context.Context, id string) (models.Image, error) {
	var image models.Image

	query := fmt.Sprintf(`
				SELECT id, advert_id, created_at, deleted
				FROM %s
				WHERE id = $1
	`, imagesTable)

	err := s.db.QueryRow(ctx, query, id).Scan(
		&image.ID,
		&image.AdvertID,
		&image.CreatedAt,
		&image.Deleted,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return image, custom_error.CustomError{
				Field:   "id",
				Message: ErrImageNotFound.Error(),
			}
		}
		return image, err
	}

	return image, nil
}
