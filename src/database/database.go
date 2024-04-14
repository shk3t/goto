package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) *pgxpool.Pool {
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(ctx, databaseUrl)
    defer pool.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create a pool for DB: %v\n", err)
		os.Exit(1)
	}

	InitSchema(ctx, pool)

	return pool
}