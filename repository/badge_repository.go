package repository

import (
	"citra-wastra-be/models"
	"time"

	"gorm.io/gorm"
)

type BadgeRepository interface {
	GetAllBadges() ([]models.Badge, error)
	GetUnlockedBadgeIDs(userID string) ([]string, error)
	UnlockBadge(userID string, badgeID string) error
	GetEquippedBadge(userID string) (models.UserBadge, error)
	EquipBadge(userID string, badgeID string) error
}

type badgeRepository struct {
	db *gorm.DB
}

func NewBadgeRepository(db *gorm.DB) BadgeRepository {
	return &badgeRepository{db}
}

func (r *badgeRepository) GetAllBadges() ([]models.Badge, error) {
	var badges []models.Badge
	err := r.db.Order("min_xp asc").Find(&badges).Error
	return badges, err
}

func (r *badgeRepository) GetUnlockedBadgeIDs(userID string) ([]string, error) {
	var badgeIDs []string
	err := r.db.Model(&models.UserBadge{}).Where("user_id = ?", userID).Pluck("badge_id", &badgeIDs).Error
	return badgeIDs, err
}

func (r *badgeRepository) UnlockBadge(userID string, badgeID string) error {
	now := time.Now()
	ub := models.UserBadge{
		UserID:     userID,
		BadgeID:    badgeID,
		AchievedAt: now,
	}
	return r.db.Where(models.UserBadge{
		UserID:  userID,
		BadgeID: badgeID,
	}).FirstOrCreate(&ub).Error
}

func (r *badgeRepository) GetEquippedBadge(userID string) (models.UserBadge, error) {
	var ub models.UserBadge
	err := r.db.Preload("Badge").Where("user_id = ? AND is_equipped = ?", userID, true).First(&ub).Error
	return ub, err
}

func (r *badgeRepository) EquipBadge(userID string, badgeID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&models.UserBadge{}).Where("user_id = ?", userID).Update("is_equipped", false)

		return tx.Model(&models.UserBadge{}).Where("user_id = ? AND badge_id = ?", userID, badgeID).Update("is_equipped", true).Error
	})
}
