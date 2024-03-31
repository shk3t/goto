package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func InitSchema(conn *pgx.Conn) {
	conn.Query(
		context.Background(),
		`
        CREATE TABLE IF NOT EXISTS project (
            id SERIAL PRIMARY KEY,
            name VARCHAR(64) NOT NULL,
            path VARCHAR(64) NOT NULL
        );

        CREATE TABLE IF NOT EXISTS task (
            id SERIAL PRIMARY KEY,
            project_id INTEGER NOT NULL REFERENCES project(id),
            name VARCHAR(64) NOT NULL,
            description TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS solution (
            id SERIAL PRIMARY KEY,
            task_id INTEGER NOT NULL REFERENCES task(id),
            code TEXT NOT NULL,
            correct BOOLEAN,
            error TEXT,
            updated_at TIMESTAMP DEFAULT NOW();
        );
        `,
	)

	fmt.Println("Schema inited successfully!")
}