package models

import "time"

type Shop struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     string    `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	ShopName   string    `gorm:"not null" json:"shop_name"`
	Address    string    `json:"address"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
