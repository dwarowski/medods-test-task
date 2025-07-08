package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase() (*gorm.DB, error) {
	dbc := os.Getenv("DATABASE_CONFIG")
	if dbc == "" {
		return nil, fmt.Errorf("DATABASE_CONFIG environment variable not set")
	}

	db, err := gorm.Open(postgres.Open(dbc), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, err
}
