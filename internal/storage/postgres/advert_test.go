package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

	advertID := uuid.New().String()
	advert := models.Advert{
		ID:          advertID,
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
			AdvertID:  advertID,
			CreatedAt: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
			Deleted:   false,
		}},
	}

	insertAdvert := fmt.Sprintf(`
				INSERT INTO %s (id, title, description, price, created_at, updated_at, user_id, deleted)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, advertsTable)

	insertImage := fmt.Sprintf(`
				INSERT INTO %s (id, advert_id, created_at, deleted)
				VALUES ($1, $2, $3, $4)
	`, imagesTable)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(insertAdvert)).WithArgs(
		advert.ID,
		advert.Title,
		advert.Description,
		advert.Price,
		advert.CreatedAt,
		advert.UpdatedAt,
		advert.UserID,
		advert.Deleted,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec(regexp.QuoteMeta(insertImage)).WithArgs(
		advert.Images[0].ID,
		advert.Images[0].AdvertID,
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

	advertID := uuid.New().String()
	userID := uuid.New().String()

	updateAdverts := fmt.Sprintf(`
				UPDATE %s
				SET deleted = TRUE
				WHERE id = $1 AND user_id = $2
	`, advertsTable)

	updateImages := fmt.Sprintf(`
				UPDATE %s
				SET deleted = TRUE
				WHERE advert_id = $1 RETURNING id
	`, imagesTable)

	columns := []string{"id"}
	rows := pgxmock.NewRows(columns).
		AddRow("test id 1").
		AddRow("test id 2").
		AddRow("test id 3")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateAdverts)).WithArgs(advertID, userID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectQuery(regexp.QuoteMeta(updateImages)).WithArgs(advertID).WillReturnRows(rows)
	mock.ExpectCommit()

	storage := NewPostgresStorage(mock)

	images, err := storage.DeleteAdvert(context.Background(), advertID, userID)
	require.NoError(t, err)

	expectedImages := []string{"test id 1", "test id 2", "test id 3"}
	require.ElementsMatch(t, expectedImages, images)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestPostgresStorageDeleteAdvertError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	advertID := uuid.New().String()
	userID := uuid.New().String()

	updateAdverts := fmt.Sprintf(`
				UPDATE %s
				SET deleted = TRUE
				WHERE id = $1 AND user_id = $2
	`, advertsTable)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateAdverts)).WithArgs(advertID, userID).WillReturnError(ErrAdvertNotFound)
	mock.ExpectRollback()

	storage := NewPostgresStorage(mock)

	images, err := storage.DeleteAdvert(context.Background(), advertID, userID)
	require.ErrorIs(t, err, ErrAdvertNotFound)
	require.Len(t, images, 0)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestPostgresStorageGetAdvertByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

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

	expectedID := uuid.New().String()
	expectedImageIDs := []string{"id1", "id2"}

	expectedAdvert := models.Advert{
		ID:          expectedID,
		Title:       "test",
		Description: "test",
		Price:       decimal.New(1200, 0),
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		UserID:      uuid.New().String(),
		Images: []*models.Image{
			{
				ID: "id1",
			},
			{
				ID: "id2",
			},
		},
	}

	columns := []string{"id", "title", "desctiption", "price", "created_at", "updated_at", "user_id", "images"}
	rows := pgxmock.NewRows(columns).
		AddRow(expectedID,
			expectedAdvert.Title,
			expectedAdvert.Description,
			expectedAdvert.Price,
			expectedAdvert.CreatedAt,
			expectedAdvert.UpdatedAt,
			expectedAdvert.UserID,
			expectedImageIDs)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expectedID).WillReturnRows(rows)

	storage := NewPostgresStorage(mock)

	advert, err := storage.GetAdvertByID(context.Background(), expectedID)
	require.NoError(t, err)
	require.Equal(t, expectedAdvert, advert)
	for i := 0; i < 2; i++ {
		require.Equal(t, expectedImageIDs[i], advert.Images[i].ID)
	}

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestPostgresStorageGetAdvertByIDError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

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

	expectedID := uuid.New().String()

	expectedAdvert := models.Advert{}

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expectedID).WillReturnError(ErrAdvertNotFound)

	storage := NewPostgresStorage(mock)

	advert, err := storage.GetAdvertByID(context.Background(), expectedID)
	require.ErrorIs(t, err, ErrAdvertNotFound)
	require.Equal(t, expectedAdvert, advert)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}
