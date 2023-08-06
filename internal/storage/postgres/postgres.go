package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romandnk/advertisement/configs"
)

var (
	usersTable   = "users"
	advertsTable = "adverts"
	imagesTable  = "images"
)

func NewPostgresDB(ctx context.Context, cfg configs.PostgresConf) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	connConfig.MaxConns = int32(cfg.MaxConns)
	connConfig.MinConns = int32(cfg.MinConns)
	connConfig.MaxConnLifetime = cfg.MaxConnLifetime
	connConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	db, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
