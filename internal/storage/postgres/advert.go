package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
)

var (
	ErrAdvertNotCreated      = errors.New("advert was not created")
	ErrAdvertImageNotCreated = errors.New("image was not created")
	ErrAdvertNotFound        = errors.New("advert not found")
)

func (s *PostgresStorage) CreateAdvert(ctx context.Context, advert models.Advert) (string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	insertAdvert := fmt.Sprintf(`
				INSERT INTO %s (id, title, description, price, created_at, updated_at, user_id, deleted)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, advertsTable)

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
		return "", err
	}

	if ct.RowsAffected() == 0 {
		return "", custom_error.CustomError{Field: "", Message: ErrAdvertNotCreated.Error()}
	}

	insertImage := fmt.Sprintf(`
				INSERT INTO %s (id, advert_id, created_at, deleted)
				VALUES ($1, $2, $3, $4)
	`, imagesTable)

	for _, image := range advert.Images {
		ct, err := tx.Exec(ctx, insertImage, image.ID, image.AdvertID, image.CreatedAt, image.Deleted)
		if err != nil {
			return "", err
		}
		if ct.RowsAffected() == 0 {
			return "", custom_error.CustomError{Field: "images", Message: ErrAdvertImageNotCreated.Error()}
		}

	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return advert.ID, nil
}

func (s *PostgresStorage) DeleteAdvert(ctx context.Context, advertID, userID string) ([]string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	updateAdverts := fmt.Sprintf(`
				UPDATE %s
				SET deleted = TRUE
				WHERE id = $1 AND user_id = $2
	`, advertsTable)

	ct, err := tx.Exec(ctx, updateAdverts, advertID, userID)
	if err != nil {
		return nil, err
	}

	if ct.RowsAffected() == 0 {
		return nil, custom_error.CustomError{Field: "id", Message: ErrAdvertNotFound.Error()}
	}

	updateImages := fmt.Sprintf(`
				UPDATE %s
				SET deleted = TRUE
				WHERE advert_id = $1 RETURNING id
	`, imagesTable)

	rows, err := tx.Query(ctx, updateImages, advertID)
	if err != nil {
		return nil, err
	}

	var imageIDs []string

	for rows.Next() {
		var imageID string

		err = rows.Scan(&imageID)
		if err != nil {
			return nil, err
		}

		imageIDs = append(imageIDs, imageID)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return imageIDs, nil
}

func (s *PostgresStorage) GetAdvertByID(ctx context.Context, id string) (models.Advert, error) {
	var advert models.Advert
	var imageIDs []string

	query := fmt.Sprintf(`
				SELECT
    			a.id,
    			a.title,
    			a.description,
    			a.price,
    			a.created_at,
    			a.updated_at,
    			a.user_id,
    			ARRAY_AGG(i.id) as images
				FROM %s a
				JOIN %s i ON a.id = i.advert_id
				WHERE a.id = $1 AND a.deleted = false AND i.deleted = false
				GROUP BY a.id
	`, advertsTable, imagesTable)

	err := s.db.QueryRow(ctx, query, id).Scan(
		&advert.ID,
		&advert.Title,
		&advert.Description,
		&advert.Price,
		&advert.CreatedAt,
		&advert.UpdatedAt,
		&advert.UserID,
		&imageIDs)

	for _, imageID := range imageIDs {
		advert.Images = append(advert.Images, &models.Image{ID: imageID})
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return advert, custom_error.CustomError{Field: "id", Message: ErrAdvertNotFound.Error()}
		}
		return advert, err
	}

	return advert, nil
}
