package storages

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type StorageConfig struct {
	Host          string
	Port          string
	Username      string
	Password      string
	DBName        string
	RetryAttempts int
}

func NewClient(ctx context.Context, cfg StorageConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err = errors.New("")
	var pool *pgxpool.Pool
	for i := 0; err != nil && i < cfg.RetryAttempts; i++ {
		func() {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			pool, err = pgxpool.New(ctx, dsn)
		}()
	}
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	return pool, nil
}
