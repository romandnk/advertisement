package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
)

func (s *PostgresStorage) CreateAdvert(ctx context.Context, advert models.Advert) (string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", custom_error.CustomError{Field: "", Message: err.Error()}
	}
	defer tx.Rollback(ctx)

	insertAdvert := fmt.Sprintf(`
				INSERT INTO %s (id, title, description, price, created_at, updated_at, user_id, deleted)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, advertsTable)

	ct, err := tx.Exec(ctx, insertAdvert,
		advert.ID,
		advert.Title,
		advert.Description,
		advert.Price,
		advert.CreatedAt,
		advert.UpdatedAt,
		advert.UserID,
		advert.Deleted,
	)
	if err != nil {
		return "", custom_error.CustomError{Field: "", Message: err.Error()}
	}

	if ct.RowsAffected() == 0 {
		return "", custom_error.CustomError{Field: "", Message: "advert was not inserted"}
	}

	insertImage := fmt.Sprintf(`
				INSERT INTO %s (id, advert_id, created_at, deleted)
				VALUES ($1, $2, $3, $4)`, imagesTable)

	for _, image := range advert.Images {
		ct, err := tx.Exec(ctx, insertImage, image.ID, advert.ID, image.CreatedAt, image.Deleted)
		if err != nil {
			return "", custom_error.CustomError{Field: "", Message: err.Error()}
		}
		if ct.RowsAffected() == 0 {
			return "", custom_error.CustomError{Field: "images", Message: "image was not inserted"}
		}

	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", custom_error.CustomError{Field: "", Message: err.Error()}
	}

	return advert.ID, nil
}

func (s *PostgresStorage) DeleteAdvert(ctx context.Context, id string) ([]string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, custom_error.CustomError{Field: "", Message: err.Error()}
	}
	defer tx.Rollback(ctx)

	updateAdverts := fmt.Sprintf(`
				UPDATE %s
				SET deleted = TRUE
				WHERE id = $1`, advertsTable)

	ct, err := tx.Exec(ctx, updateAdverts, id)
	if err != nil {
		return nil, custom_error.CustomError{Field: "", Message: err.Error()}
	}

	if ct.RowsAffected() == 0 {
		return nil, custom_error.CustomError{Field: "id", Message: pgx.ErrNoRows.Error()}
	}

	updateImages := fmt.Sprintf(`
				UPDATE %s
				SET deleted = TRUE
				WHERE advert_id = $1 RETURNING id`, imagesTable)

	rows, err := tx.Query(ctx, updateImages, id)
	if err != nil {
		return nil, custom_error.CustomError{Field: "", Message: err.Error()}
	}

	var imageIDs []string

	for rows.Next() {
		var imageID string

		err = rows.Scan(&imageID)
		if err != nil {
			return nil, custom_error.CustomError{Field: "", Message: err.Error()}
		}

		imageIDs = append(imageIDs, imageID)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, custom_error.CustomError{Field: "", Message: err.Error()}
	}

	return imageIDs, nil
}
