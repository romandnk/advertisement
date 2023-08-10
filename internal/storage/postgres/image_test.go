package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestPostgresStorageGetImageByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	query := fmt.Sprintf(`
				SELECT id, advert_id, created_at, deleted
				FROM %s
				WHERE id = $1
	`, imagesTable)

	expectedID := "test id 1"
	createdAt := time.Now()

	columns := []string{"id", "advert_id", "created_at", "deleted"}
	rows := pgxmock.NewRows(columns).
		AddRow("test id 1", "advert id 1", createdAt, false)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expectedID).WillReturnRows(rows)

	storage := NewPostgresStorage(mock)

	image, err := storage.GetImageByID(context.Background(), expectedID)
	require.NoError(t, err)

	expectedImage := models.Image{
		ID:        expectedID,
		Data:      nil,
		AdvertID:  "advert id 1",
		CreatedAt: createdAt,
		Deleted:   false,
	}

	require.Equal(t, expectedImage, image)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestPostgresStorageGetImageByIDError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	query := fmt.Sprintf(`
				SELECT id, advert_id, created_at, deleted
				FROM %s
				WHERE id = $1
	`, imagesTable)

	expectedError := custom_error.CustomError{
		Field:   "id",
		Message: ErrImageNotFound.Error(),
	}
	expectedID := uuid.New().String()

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expectedID).WillReturnError(expectedError)

	storage := NewPostgresStorage(mock)

	image, err := storage.GetImageByID(context.Background(), expectedID)
	require.ErrorIs(t, err, expectedError)
	require.Equal(t, models.Image{}, image)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}
