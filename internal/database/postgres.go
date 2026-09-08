package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// note: use pgx instead of standard database/sql package because pgx is faster and manages concurrent connection pooling more efficiently under massive load

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

/*

func NewPostgres(ctx context.Context, cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	db := sql.OpenDB(stdlib.GetConnector(dsn))


	// Pool configuration
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

*/

// initializes connection pool to postgresql
func InitPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {

	// Parse the raw connection string into a configuration struct
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute

	// Create the actual connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// verify database is awake and healthy
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when pinging: %w", err)
	}

	slog.Info("Database connection pool successfully initialized and pinged")

	return pool, nil
}
