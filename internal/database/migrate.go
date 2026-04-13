package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunProductionMigrations() error {
	if os.Getenv("APP_ENV") != "production" {
		return nil
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	// This path matches your Dockerfile.render copy destination
	sourceURL := "file:///app/infra/migrations"

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			fmt.Printf("migration source close error: %v\n", sourceErr)
		}
		if dbErr != nil {
			fmt.Printf("migration database close error: %v\n", dbErr)
		}
	}()

	// Optional graceful stop channel
	stopChan := make(chan bool, 1)
	m.GracefulStop = stopChan

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}
