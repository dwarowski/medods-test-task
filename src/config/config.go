package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connection to database via environment variable
func ConnectToDatabase() (*gorm.DB, error) {

	// Try getting env
	dbc := os.Getenv("DATABASE_CONFIG")
	if dbc == "" {
		return nil, fmt.Errorf("DATABASE_CONFIG environment variable not set")
	}

	// Try connecting to db
	db, err := gorm.Open(postgres.Open(dbc), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, err
}
