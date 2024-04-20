package database

import (
	"context"
	"log"
)

func InitSchema(ctx context.Context) {
	tableDefinitions := []string{
		`
        CREATE TABLE IF NOT EXISTS project (
            id SERIAL PRIMARY KEY,
            url VARCHAR(256),
            container VARCHAR(128) NOT NULL UNIQUE,
            name VARCHAR(64) NOT NULL,
            language VARCHAR(64) NOT NULL,
            containerization VARCHAR(64) NOT NULL DEFAULT 'docker',
            srcdir VARCHAR(64) NOT NULL DEFAULT 'src',
            stubdir VARCHAR(64) NOT NULL DEFAULT 'stubs'
        );`,
		`
        CREATE TABLE IF NOT EXISTS project_module (
            id SERIAL PRIMARY KEY,
            project_id INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
            name VARCHAR(64) NOT NULL
        );`,

		`
        CREATE TABLE IF NOT EXISTS task (
            id SERIAL PRIMARY KEY,
            project_id INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
            name VARCHAR(64) NOT NULL,
            description TEXT NOT NULL,
            runtarget VARCHAR(256) NOT NULL,
            UNIQUE(project_id, name)
        );`,
		`
        CREATE TABLE IF NOT EXISTS injectfile (
            id SERIAL PRIMARY KEY,
            task_id INTEGER NOT NULL REFERENCES task(id) ON DELETE CASCADE,
            name VARCHAR(64) NOT NULL,
            path VARCHAR(256) NOT NULL
            UNIQUE(task_id, name)
        );`,

		`
        CREATE TABLE IF NOT EXISTS solution (
            id SERIAL PRIMARY KEY,
            task_id INTEGER NOT NULL REFERENCES task(id),
            status VARCHAR(64) NOT NULL,
            code TEXT NOT NULL,
            respone TEXT,
            updated_at TIMESTAMP DEFAULT NOW()
        );`,
	}

	for _, tableDef := range tableDefinitions {
		rows, err := ConnPool.Query(ctx, tableDef)
		if err != nil {
			panic("Schema initiation failed: " + err.Error())
		}
		rows.Close()
	}

	log.Println("Schema inited successfully!")
}