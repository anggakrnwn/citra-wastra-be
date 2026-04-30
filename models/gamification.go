package models

import "time"

type Island struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string `gorm:"not null"`
	Description string
	ImageURL    string
	Modules     []Module `gorm:"foreignKey:IslandID"`
}

type Module struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	IslandID    string `gorm:"index"`
	Name        string `gorm:"not null"`
	Description string
	Levels      []Level `gorm:"foreignKey:ModuleID"`
}

type Level struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ModuleID string `gorm:"index"`
	Order    int    `gorm:"not null"`
	Title    string `gorm:"not null"`
	Content  string `gorm:"type:text"`
	XPReward int    `gorm:"default:50"`
}

type UserLevelProgress struct {
	UserID      string    `gorm:"primaryKey"`
	LevelID     string    `gorm:"primaryKey"`
	CompletedAt time.Time `gorm:"autoCreateTime"`
}

type Question struct {
	ID       string `gorm:"primaryKey;type:uuid"`
	LevelID  string `gorm:"type:uuid;index"`
	Question string `gorm:"type:text"`
	OptionA  string `gorm:"type:text"`
	OptionB  string `gorm:"type:text"`
	OptionC  string `gorm:"type:text"`
	OptionD  string `gorm:"type:text"`
	Correct  string `gorm:"size:1"`
}
