package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitSchema(ctx context.Context, pool *pgxpool.Pool) {

	tableDefinitions := []string{
		`
        CREATE TABLE IF NOT EXISTS project (
            id SERIAL PRIMARY KEY,
            name VARCHAR(64) NOT NULL,
            url VARCHAR(256)
        );`,
		`
        CREATE TABLE IF NOT EXISTS task (
            id SERIAL PRIMARY KEY,
            project_id INTEGER NOT NULL REFERENCES project(id),
            name VARCHAR(64) NOT NULL,
            description TEXT NOT NULL
        );`,
		`
        CREATE TABLE IF NOT EXISTS solution (
            id SERIAL PRIMARY KEY,
            task_id INTEGER NOT NULL REFERENCES task(id),
            code TEXT NOT NULL,
            correct BOOLEAN,
            error TEXT,
            updated_at TIMESTAMP DEFAULT NOW()
        );`,
	}

	for _, tableDef := range tableDefinitions {
		rows, err := pool.Query(ctx, tableDef)
		if err != nil {
			panic("Schema initiation failed: " + err.Error())
		}
        rows.Close()
	}

	fmt.Println("Schema inited successfully!")
}