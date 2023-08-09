package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestPostgresStorageCreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	query := fmt.Sprintf(`
			INSERT INTO %s (id, email, password, created_at, updated_at, deleted)
			VALUES ($1, $2, $3, $4, $5, $6)
	`, usersTable)

	expectedUser := models.User{
		ID:        uuid.New().String(),
		Email:     "test@mail.ru",
		Password:  "test_password",
		CreatedAt: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
		Deleted:   false,
	}

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		expectedUser.ID,
		expectedUser.Email,
		expectedUser.Password,
		expectedUser.CreatedAt,
		expectedUser.UpdatedAt,
		expectedUser.Deleted,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	storage := NewPostgresStorage(mock)
	ctx := context.Background()

	id, err := storage.CreateUser(ctx, expectedUser)
	require.NoError(t, err)
	require.Equal(t, expectedUser.ID, id)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}
