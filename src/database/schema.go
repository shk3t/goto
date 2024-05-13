package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var TABLE_DEFINITIONS = [...]string{
	`
    CREATE TABLE IF NOT EXISTS "user" (
        id SERIAL PRIMARY KEY,
        login VARCHAR(64) NOT NULL UNIQUE,
        password VARCHAR(128) NOT NULL,
        is_admin BOOLEAN NOT NULL DEFAULT false
    )`,
	`
    CREATE TABLE IF NOT EXISTS project (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
        dir VARCHAR(128) NOT NULL UNIQUE,
        name VARCHAR(64) NOT NULL,
        language VARCHAR(64) NOT NULL,
        containerization VARCHAR(64) NOT NULL DEFAULT 'docker',
        srcdir VARCHAR(64) NOT NULL DEFAULT 'src',
        stubdir VARCHAR(64) NOT NULL DEFAULT 'stubs',
        updated_at TIMESTAMP DEFAULT NOW()
    )`,
	`
    CREATE TABLE IF NOT EXISTS module (
        id SERIAL PRIMARY KEY,
        project_id INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
        name VARCHAR(64) NOT NULL
    )`,

	`
    CREATE TABLE IF NOT EXISTS task (
        id SERIAL PRIMARY KEY,
        project_id INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
        name VARCHAR(64) NOT NULL,
        description TEXT NOT NULL,
        runtarget VARCHAR(256) NOT NULL,
        UNIQUE(project_id, name)
    )`,
	`
    CREATE TABLE IF NOT EXISTS task_file (
        id SERIAL PRIMARY KEY,
        task_id INTEGER NOT NULL REFERENCES task(id) ON DELETE CASCADE,
        name VARCHAR(64) NOT NULL,
        path VARCHAR(256) NOT NULL,
        stub TEXT NOT NULL,
        UNIQUE(task_id, name)
    )`,

	`
    CREATE TABLE IF NOT EXISTS solution (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
        task_id INTEGER NOT NULL REFERENCES task(id) ON DELETE CASCADE,
        status VARCHAR(64) NOT NULL DEFAULT 'check',
        result TEXT NOT NULL DEFAULT '',
        updated_at TIMESTAMP DEFAULT NOW(),
        UNIQUE(user_id, task_id)
    )`,
	`
    CREATE TABLE IF NOT EXISTS solution_file (
        id SERIAL PRIMARY KEY,
        solution_id INTEGER NOT NULL REFERENCES solution(id) ON DELETE CASCADE,
        name VARCHAR(64) NOT NULL,
        code TEXT NOT NULL,
        UNIQUE(solution_id, name)
    )`,

	`
    CREATE TABLE IF NOT EXISTS delayed_task (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
        target_id INTEGER,
        action VARCHAR(256) NOT NULL,
        target VARCHAR(256) NOT NULL,
        status VARCHAR(64) NOT NULL DEFAULT 'processing',
        details TEXT NOT NULL DEFAULT '',
        updated_at TIMESTAMP DEFAULT NOW()
    )`,
}

func InitSchema(ctx context.Context) {
	tx, _ := ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	for _, tableDef := range TABLE_DEFINITIONS {
		_, err := tx.Exec(ctx, tableDef)
		if err != nil {
			panic("Schema initiation failed: " + err.Error())
		}
	}

	tx.Commit(ctx)
	log.Println("Schema inited successfully!")
}