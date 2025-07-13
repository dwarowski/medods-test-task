package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connection to database via environment variable
func ConnectToDatabase() (*gorm.DB, error) {

	// Try getting envs
	dbc := ""
	vars := map[string]string{"host": "POSTGRES_HOST", "user": "POSTGRES_USER", "password": "POSTGRES_PASSWORD", "dbname": "POSTGRES_DB", "port": "POSTGRES_PORT"}
	for key, value := range vars {
		env_var := os.Getenv(value)
		if env_var == "" {
			return nil, fmt.Errorf("%v environment variable not set", value)
		}
		dbc = dbc + " " + key + "=" + env_var
	}

	if dbc == "" {
		return nil, fmt.Errorf("env variables error")
	}

	// Try connecting to db
	db, err := gorm.Open(postgres.Open(dbc), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, err
}
