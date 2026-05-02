package models

import "time"

type Island struct {
	ID          string   `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string   `gorm:"not null" json:"name"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image_url"`
	Modules     []Module `gorm:"foreignKey:IslandID" json:"modules,omitempty"`
}

type Module struct {
	ID          string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	IslandID    string  `gorm:"index" json:"island_id"`
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Levels      []Level `gorm:"foreignKey:ModuleID" json:"levels,omitempty"`
}

type Level struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ModuleID string `gorm:"index" json:"module_id"`
	Order    int    `gorm:"not null" json:"order"`
	Title    string `gorm:"not null" json:"title"`
	Content  string `gorm:"type:text" json:"content"`
	XPReward int    `gorm:"default:50" json:"xp_reward"`
}

type UserLevelProgress struct {
	UserID      string    `gorm:"primaryKey" json:"user_id"`
	LevelID     string    `gorm:"primaryKey" json:"level_id"`
	CompletedAt time.Time `gorm:"autoCreateTime" json:"completed_at"`
}

type Question struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	LevelID  string `gorm:"type:uuid;index" json:"level_id"`
	Question string `gorm:"type:text" json:"question"`
	OptionA  string `gorm:"type:text" json:"option_a"`
	OptionB  string `gorm:"type:text" json:"option_b"`
	OptionC  string `gorm:"type:text" json:"option_c"`
	OptionD  string `gorm:"type:text" json:"option_d"`
	Correct  string `gorm:"size:1" json:"correct"`
}
