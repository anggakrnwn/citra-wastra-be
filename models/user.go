package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	XP        int            `gorm:"default:0" json:"xp"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Profile UserProfile `json:"profile" gorm:"foreignKey:UserID"`
	Shop    *Shop       `json:"shop,omitempty" gorm:"foreignKey:UserID"`
	Badges  []UserBadge `json:"badges" gorm:"foreignKey:UserID"`
}

type UserProfile struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null" json:"user_id"`
	FullName    string    `json:"full_name"`
	AvatarURL   string    `json:"avatar_url"`
	Bio         string    `json:"bio"`
	PhoneNumber string    `json:"phone_number"`
	UpdatedAt   time.Time `json:"updated_at"`
}
