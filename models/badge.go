package models

import "time"

type Badge struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	MinXP       int       `gorm:"default:0" json:"min_xp"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserBadge struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     string    `gorm:"type:uuid;not null;index" json:"user_id"`
	BadgeID    string    `gorm:"type:uuid;not null" json:"badge_id"`
	IsEquipped bool      `gorm:"default:false" json:"is_equipped"`
	AchievedAt time.Time `json:"achieved_at"`

	Badge Badge `gorm:"foreignKey:BadgeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"badge_details"`
}
