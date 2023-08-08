package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestPostgresStorageCreateAdvert(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	advert := models.Advert{
		ID:          uuid.New().String(),
		Title:       "test title",
		Description: "test description",
		Price:       decimal.New(1200, 0),
		CreatedAt:   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
		UserID:      uuid.New().String(),
		Deleted:     false,
		Images: []*models.Image{{
			ID:        uuid.New().String(),
			Data:      []byte("test data"),
			CreatedAt: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
			Deleted:   false,
		}},
	}

	queryAdvert := fmt.Sprintf(`
				INSERT INTO %s (id, title, description, price, created_at, updated_at, user_id, deleted)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, advertsTable)

	queryImage := fmt.Sprintf(`
				INSERT INTO %s (id, advert_id, created_at, deleted)
				VALUES ($1, $2, $3, $4)
	`, imagesTable)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryAdvert)).WithArgs(
		advert.ID,
		advert.Title,
		advert.Description,
		advert.Price,
		advert.CreatedAt,
		advert.UpdatedAt,
		advert.UserID,
		advert.Deleted,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryImage)).WithArgs(
		advert.Images[0].ID,
		advert.ID,
		advert.Images[0].CreatedAt,
		advert.Images[0].Deleted,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	storage := NewPostgresStorage(mock)

	id, err := storage.CreateAdvert(context.Background(), advert)
	require.NoError(t, err)
	require.Equal(t, advert.ID, id)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestPostgresStorageDeleteAdvert(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	id := uuid.New().String()

	updateAdverts := fmt.Sprintf(`
				UPDATE %s
				SET deleted = TRUE
				WHERE id = $1`, advertsTable)

	updateImages := fmt.Sprintf(`
				UPDATE %s
				SET deleted = TRUE
				WHERE advert_id = $1 RETURNING id`, imagesTable)

	columns := []string{"id"}
	rows := pgxmock.NewRows(columns).
		AddRow("test id 1").
		AddRow("test id 2").
		AddRow("test id 3")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateAdverts)).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectQuery(regexp.QuoteMeta(updateImages)).WithArgs(id).WillReturnRows(rows)
	mock.ExpectCommit()

	storage := NewPostgresStorage(mock)

	images, err := storage.DeleteAdvert(context.Background(), id)
	require.NoError(t, err)

	expectedImages := []string{"test id 1", "test id 2", "test id 3"}
	require.ElementsMatch(t, expectedImages, images)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestPostgresStorageDeleteAdvertError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	id := uuid.New().String()

	updateAdverts := fmt.Sprintf(`
				UPDATE %s
				SET deleted = TRUE
				WHERE id = $1`, advertsTable)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateAdverts)).WithArgs(id).WillReturnError(pgx.ErrNoRows)
	mock.ExpectRollback()

	storage := NewPostgresStorage(mock)

	images, err := storage.DeleteAdvert(context.Background(), id)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Len(t, images, 0)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}
