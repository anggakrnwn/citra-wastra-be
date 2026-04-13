package models

import "time"

type BatikCatalog struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	History   string    `gorm:"type:text;not null" json:"history"`
	Origin    string    `json:"origin"`
	ImageURL  string    `json:"image_url"`
	CreatedBy string    `gorm:"type:uuid" json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
