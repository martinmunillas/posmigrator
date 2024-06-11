package posmigrator

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
)

func RunMigrations(conn *sql.DB, config Config) error {
	migrationFiles, err := readMigrationFiles(config)
	if err != nil {
		return fmt.Errorf("error reading migration files: %w", err)
	}
	views := readViews(config)
	err = ensureMigrationsTable(conn)
	if err != nil {
		return fmt.Errorf("error ensuring migrations table: %w", err)
	}
	migrations, err := getMigrations(conn)
	fmt.Println(migrations)
	if err != nil {
		return fmt.Errorf("error getting migrations: %w", err)
	}
	err = ensureRanMigrationsAreValid(migrationFiles, migrations)
	if err != nil {
		return fmt.Errorf("error ensuring ran migrations are valid: %w", err)
	}

	err = dropAllViews(conn, views)
	if err != nil {
		return fmt.Errorf("error dropping views: %w", err)
	}
	err = migrateMigrations(conn, config, migrationFiles)
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}
	err = migrateViews(conn, config, views)
	if err != nil {
		return fmt.Errorf("error running views: %w", err)
	}

	return nil
}

func EnsureAllMigrationsRanAndAreValid(conn *sql.DB, config Config) error {
	migrationFiles, err := readMigrationFiles(config)
	if err != nil {
		return err
	}

	migrations, err := getMigrations(conn)
	if err != nil {
		return err
	}

	if err := ensureRanMigrationsAreValid(migrationFiles, migrations); err != nil {
		return err
	}

	if len(migrations) != len(migrationFiles) {
		return fmt.Errorf("there are %d migrations left to run", len(migrationFiles)-len(migrations))
	}

	green.Println("Migrations are up to date")
	return nil
}

func getMigrations(conn *sql.DB) ([]Migration, error) {
	var migrations []Migration
	rows, err := conn.Query("SELECT id, description, migrated_at FROM migrations ORDER BY id")
	for rows.Next() {
		var migration Migration
		err = rows.Scan(&migration.ID, &migration.Description, &migration.MigratedAt)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}
	if err != nil {
		return nil, err
	}
	return migrations, nil
}

func ensureMigrationsAreRelated(migrationFile migrationFile, migration Migration) error {
	if migration.ID != migrationFile.id {
		return fmt.Errorf("inconsistent migration state, expected migration with id %d, got %d. this requires manual solution", migrationFile.id, migration.ID)
	}
	if migration.Description != migrationFile.name {
		return fmt.Errorf("inconsistent migration state, expected migration with description %s, got %s. this requires manual solution", migrationFile.name, migration.Description)
	}
	return nil
}

func ensureRanMigrationsAreValid(migrationFiles []migrationFile, migrations []Migration) error {
	for i, migration := range migrations {
		migrationFile := migrationFiles[i]
		if err := ensureMigrationsAreRelated(migrationFile, migration); err != nil {
			return err
		}
	}

	return nil
}

func migrateMigrations(conn *sql.DB, config Config, migrationFiles []migrationFile) error {
	migrations, err := getMigrations(conn)
	if err != nil {
		return err
	}

	if len(migrations) == len(migrationFiles) {
		green.Println("All migrations already ran")
		return nil
	} else {
		blue.Printf("%d/%d migrations already ran\n", len(migrations), len(migrationFiles))
	}
	blue.Println("Running migrations")

withNextMigration:
	for i, migrationFile := range migrationFiles {
		// Skip migrations that have already ran
		if i < len(migrations) {
			if err := ensureMigrationsAreRelated(migrationFile, migrations[i]); err != nil {
				return err
			}
			continue withNextMigration
		}

		sql, err := migrationFile.sql(config)
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %w", migrationFile.fileName, err)
		}
		if len(sql) == 0 {
			return fmt.Errorf("empty migration file %s", migrationFile.fileName)
		}

		fmt.Printf("Running migration %s\n", migrationFile.fileName)

		tx, err := conn.BeginTx(context.Background(), nil)
		if err != nil {
			return err
		}
		_, err = tx.Exec(string(sql))
		if err != nil {
			return err
		}

		_, err = tx.Exec("INSERT INTO migrations (id, description, migrated_at) VALUES ($1, $2, $3)", migrationFile.id, migrationFile.name, time.Now())

		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		green.Printf("Migration %s ran successfully\n", migrationFile.fileName)
	}

	green.Println("All migrations ran successfully")
	return nil
}

func readViews(config Config) []string {
	viewDirs, err := os.ReadDir(config.MigrationsPath + "/views")
	if err != nil {
		red.Println("Error reading views directory")
		fmt.Println(err)
		os.Exit(1)
	}

	var viewNames []string
	for _, dir := range viewDirs {
		if dir.IsDir() {
			continue
		}
		viewNames = append(viewNames, strings.Replace(dir.Name(), ".sql", "", -1))
	}

	return viewNames
}

func migrateViews(conn *sql.DB, config Config, views []string) error {
	blue.Println("Setting up views")

	for _, view := range views {

		sql, err := os.ReadFile(config.MigrationsPath + "/views/" + view + ".sql")
		if err != nil {
			return fmt.Errorf("error reading view file %s: %w", view, err)
		}

		if len(sql) == 0 {
			return fmt.Errorf("empty view file %s", view)
		}

		_, err = conn.Exec(string(sql))

		if err != nil {
			return fmt.Errorf("error running view %s: %w", view, err)
		}

	}

	green.Println("All views are set up")

	return nil
}

func dropAllViews(conn *sql.DB, views []string) error {
	for _, view := range views {
		conn.Exec(fmt.Sprintf("DROP VIEW IF EXISTS %s;", view))
	}
	return nil
}
