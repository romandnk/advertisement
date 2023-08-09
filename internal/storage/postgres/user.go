package postgres

import (
	"context"
	"fmt"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
)

func (s *PostgresStorage) CreateUser(ctx context.Context, user models.User) (string, error) {
	query := fmt.Sprintf(`
			INSERT INTO %s (id, email, password, created_at, updated_at, deleted)
			VALUES ($1, $2, $3, $4, $5, $6)
	`, usersTable)

	ct, err := s.db.Exec(ctx, query, user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.Deleted)
	if err != nil {
		return "", custom_error.CustomError{Field: "", Message: err.Error()}
	}

	if ct.RowsAffected() == 0 {
		return "", custom_error.CustomError{Field: "", Message: "user was not created"}
	}

	return user.ID, nil
}
