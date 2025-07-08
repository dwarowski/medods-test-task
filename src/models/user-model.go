package models

type User struct {
	ID           uint `gorm:"primaryKey"`
	Username     string
	Email        string
	Password     string `json:"-"`
	RefreshToken string
}
