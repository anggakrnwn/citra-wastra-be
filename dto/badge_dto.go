package dto

import "time"

type BadgeResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ImageURL    string     `json:"image_url"`
	MinXP       int        `json:"min_xp"`
	IsUnlocked  bool       `json:"is_unlocked"`
	IsEquipped  bool       `json:"is_equipped"`
	AchievedAt  *time.Time `json:"achieved_at,omitempty"`
}

type UserProfileBadgeResponse struct {
	TotalXP int             `json:"total_xp"`
	Badges  []BadgeResponse `json:"badges"`
}

type EquipBadgeRequest struct {
	BadgeID string `json:"badge_id" binding:"required"`
}

type XPJobPayload struct {
	UserID string `json:"user_id"`
	XPGain int    `json:"xp_gain"`
	Retry  int    `json:"retry"`
	JobID  string `json:"job_id"`
}
