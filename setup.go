package posmigrator

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectPostgres(config Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		config.DbUser,
		config.DbName,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	fmt.Println("Successfully connected to PostgreSQL database")

	return db, nil
}

func ensureMigrationsTable(conn *sql.DB) error {
	var exists bool
	err := conn.QueryRow(
		`SELECT EXISTS (
			SELECT 1
			FROM pg_tables
			WHERE schemaname = 'public'
			AND tablename = 'migrations'
		);`,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking migrations table: %w", err)
	}

	if !exists {
		fmt.Println("Creating migrations table")
		_, err := conn.Exec(
			`CREATE TABLE "migrations" (
				"id" INT PRIMARY KEY NOT NULL,
				"description" VARCHAR NOT NULL,
				"migrated_at" TIMESTAMPTZ NOT NULL
			);
			`,
		)
		if err != nil {
			return fmt.Errorf("error creating migrations table: %w", err)
		}
		green.Println("Migrations table created")
	}

	return nil
}
