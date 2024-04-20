package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ConnPool *pgxpool.Pool

func Connect(ctx context.Context) {
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	ConnPool, err = pgxpool.New(ctx, databaseUrl)
	defer ConnPool.Close()
	if err != nil {
		panic(err)
	}

	InitSchema(ctx)
}