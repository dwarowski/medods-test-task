package models

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Username     string
	Email        string `gorm:"index:idx_name,unique"`
	Password     string `json:"-"`
	RefreshToken string
}
