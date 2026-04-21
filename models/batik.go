package models

import "time"

type Batik struct {
	ID                 string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID             string    `gorm:"type:uuid;not null;index" json:"user_id"`
	ImageURL           string    `gorm:"type:text;not null" json:"image_url"`
	Label              string    `gorm:"type:varchar(100);not null" json:"label"`
	Confidence         float64   `gorm:"not null" json:"confidence"`
	PhilosophyNarrator string    `gorm:"type:text" json:"philosophy_narrator"`
	Philosophy         string    `gorm:"type:text" json:"philosophy"`
	CreatedAt          time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
