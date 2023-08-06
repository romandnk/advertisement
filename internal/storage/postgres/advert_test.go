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
	"os"
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
		}},
	}

	queryAdvert := fmt.Sprintf(`
				INSERT INTO %s (id, title, description, price, created_at, updated_at, user_id, deleted)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, advertsTable)

	queryImage := fmt.Sprintf(`
				INSERT INTO %s (id, advert_id, created_at)
				VALUES ($1, $2, $3)
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
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	storage := NewPostgresStorage(mock)

	dir := t.TempDir()

	id, err := storage.CreateAdvert(context.Background(), advert, dir)
	require.NoError(t, err)
	require.Equal(t, advert.ID, id)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestSaveImage(t *testing.T) {
	dir := t.TempDir()

	image := &models.Image{
		ID:        uuid.New().String(),
		Data:      []byte("test"),
		CreatedAt: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	err := saveImage(image, dir)
	require.NoError(t, err)

	_, err = os.Stat(dir + image.ID + ".jpg")
	require.NoError(t, err)
}

func TestPostgresStorageDeleteAdvert(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	id := uuid.New().String()

	querySelect := fmt.Sprintf(`
				SELECT id FROM %s
				WHERE advert_id = $1`, imagesTable)

	queryAdvert := fmt.Sprintf(`
				DELETE FROM %s
				WHERE id = $1`, advertsTable)

	columns := []string{"id"}
	rows := pgxmock.NewRows(columns).AddRow("test image id")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelect)).WithArgs(id).WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(queryAdvert)).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	storage := NewPostgresStorage(mock)

	dir := t.TempDir()
	err = saveImage(&models.Image{ID: "test image id"}, dir)
	require.NoError(t, err)

	err = storage.DeleteAdvert(context.Background(), id, dir)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestPostgresStorageDeleteAdvertError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	id := uuid.New().String()

	querySelect := fmt.Sprintf(`
				SELECT id FROM %s
				WHERE advert_id = $1`, imagesTable)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelect)).WithArgs(id).WillReturnError(pgx.ErrNoRows)
	mock.ExpectRollback()

	storage := NewPostgresStorage(mock)

	dir := t.TempDir()

	err = storage.DeleteAdvert(context.Background(), id, dir)
	require.EqualError(t, err, pgx.ErrNoRows.Error())

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestDeleteImage(t *testing.T) {
	dir := t.TempDir()

	image := &models.Image{
		ID:        uuid.New().String(),
		Data:      []byte("test"),
		CreatedAt: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	err := saveImage(image, dir)
	require.NoError(t, err)

	err = deleteImage([]string{image.ID}, dir)
	require.NoError(t, err)
}
