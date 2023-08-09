package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
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

func (s *PostgresStorage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	query := fmt.Sprintf(`
			SELECT id, email, password, created_at, updated_at, deleted
			FROM %s 
			WHERE email = $1
	`, usersTable)

	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Deleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, custom_error.CustomError{Field: "email", Message: "invalid email"}
		}
		return user, custom_error.CustomError{Field: "", Message: err.Error()}
	}

	return user, nil
}
