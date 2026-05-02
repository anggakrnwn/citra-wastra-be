package dto

import "time"

type AdminLevelRequest struct {
	ModuleID string `json:"module_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Order    int    `json:"order"`
	XPReward int    `json:"xp_reward"`
}

type AdminModuleRequest struct {
	IslandID    string `json:"island_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type AdminQuestionRequest struct {
	LevelID  string `json:"level_id" binding:"required"`
	Question string `json:"question" binding:"required"`
	OptionA  string `json:"option_a" binding:"required"`
	OptionB  string `json:"option_b" binding:"required"`
	OptionC  string `json:"option_c" binding:"required"`
	OptionD  string `json:"option_d" binding:"required"`
	Correct  string `json:"correct" binding:"required"`
}

type AdminCatalogRequest struct {
	Name     string `json:"name" binding:"required"`
	History  string `json:"history" binding:"required"`
	Origin   string `json:"origin"`
	ImageURL string `json:"image_url"`
}

type AdminBadgeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	MinXP       int    `json:"min_xp"`
	ImageURL    string `json:"image_url"`
}

type AdminUserResponse struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	XP         int       `json:"xp"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	IsBanned   bool      `json:"is_banned"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdminSystemHealthResponse struct {
	WorkerStatus string            `json:"worker_status"`
	QueueLength  int64             `json:"queue_length"`
	Cloudinary   string            `json:"cloudinary_status"`
	APIStats     map[string]string `json:"api_stats"`
}
