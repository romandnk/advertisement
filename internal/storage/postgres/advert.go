package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
	"os"
)

func (s *PostgresStorage) CreateAdvert(ctx context.Context, advert models.Advert, path string) (string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", custom_error.CustomError{Field: "", Message: err.Error()}
	}
	defer tx.Rollback(ctx)

	queryAdvert := fmt.Sprintf(`
				INSERT INTO %s (id, title, description, price, created_at, updated_at, user_id, deleted)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, advertsTable)

	ct, err := tx.Exec(ctx, queryAdvert,
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

	queryImage := fmt.Sprintf(`
				INSERT INTO %s (id, advert_id, created_at)
				VALUES ($1, $2, $3)`, imagesTable)

	for _, image := range advert.Images {
		ct, err := tx.Exec(ctx, queryImage, image.ID, advert.ID, image.CreatedAt)
		if err != nil {
			return "", custom_error.CustomError{Field: "", Message: err.Error()}
		}
		if ct.RowsAffected() == 0 {
			return "", custom_error.CustomError{Field: "", Message: "image was not inserted"}
		}

		err = saveImage(image, path)
		if err != nil {
			return "", custom_error.CustomError{Field: "", Message: err.Error()}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", custom_error.CustomError{Field: "", Message: err.Error()}
	}

	return advert.ID, nil
}

func saveImage(image *models.Image, path string) error {
	err := os.WriteFile(path+image.ID+".jpg", image.Data, 0o644)
	return err
}

func (s *PostgresStorage) DeleteAdvert(ctx context.Context, id string, path string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return custom_error.CustomError{Field: "", Message: err.Error()}
	}
	defer tx.Rollback(ctx)

	querySelect := fmt.Sprintf(`
				SELECT id FROM %s
				WHERE advert_id = $1`, imagesTable)

	rows, err := tx.Query(ctx, querySelect, id)
	if err != nil {
		return custom_error.CustomError{Field: "", Message: err.Error()}
	}
	defer rows.Close()
	var imagesIDs []string
	for rows.Next() {
		var imagesID string
		err = rows.Scan(&imagesID)
		if err != nil {
			return custom_error.CustomError{Field: "", Message: err.Error()}
		}
		imagesIDs = append(imagesIDs, imagesID)
	}
	err = rows.Err()
	if err != nil {
		return custom_error.CustomError{Field: "", Message: err.Error()}
	}
	if len(imagesIDs) == 0 {
		return custom_error.CustomError{Field: "", Message: pgx.ErrNoRows.Error()}
	}

	queryAdvert := fmt.Sprintf(`
				DELETE FROM %s
				WHERE id = $1`, advertsTable)

	_, err = tx.Exec(ctx, queryAdvert, id)
	if err != nil {
		return custom_error.CustomError{Field: "", Message: err.Error()}
	}

	err = deleteImage(imagesIDs, path)
	if err != nil {
		return custom_error.CustomError{Field: "", Message: err.Error()}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return custom_error.CustomError{Field: "", Message: err.Error()}
	}

	return nil
}

func deleteImage(imageIDs []string, path string) error {
	for _, imageID := range imageIDs {
		err := os.Remove(path + imageID + ".jpg")
		if err != nil {
			return err
		}
	}

	return nil
}
