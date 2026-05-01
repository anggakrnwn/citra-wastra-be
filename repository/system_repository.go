package repository

import (
	"citra-wastra-be/models"

	"gorm.io/gorm"
)

type SystemRepository interface {
	CreateAuditLog(log *models.AuditLog) error
	GetAuditLogs(limit, offset int) ([]models.AuditLog, int64, error)

	GetConfig(key string) (*models.SystemConfig, error)
	SetConfig(config *models.SystemConfig) error
	GetAllConfigs() ([]models.SystemConfig, error)
	DeleteConfig(key string) error
	ResetTestData() error
}

type systemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db *gorm.DB) SystemRepository {
	return &systemRepository{db}
}

func (r *systemRepository) CreateAuditLog(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *systemRepository) GetAuditLogs(limit, offset int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64
	r.db.Model(&models.AuditLog{}).Count(&total)
	err := r.db.Preload("Admin").Limit(limit).Offset(offset).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

func (r *systemRepository) GetConfig(key string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := r.db.Where("key = ?", key).First(&config).Error
	return &config, err
}

func (r *systemRepository) SetConfig(config *models.SystemConfig) error {
	return r.db.Where("key = ?", config.Key).
		Assign(models.SystemConfig{Value: config.Value, Category: config.Category}).
		FirstOrCreate(config).Error
}

func (r *systemRepository) GetAllConfigs() ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	err := r.db.Find(&configs).Error
	return configs, err
}

func (r *systemRepository) DeleteConfig(key string) error {
	return r.db.Delete(&models.SystemConfig{}, "key = ?", key).Error
}

func (r *systemRepository) ResetTestData() error {
	tx := r.db.Begin()

	if err := tx.Exec("TRUNCATE TABLE audit_logs CASCADE").Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Exec("TRUNCATE TABLE batiks CASCADE").Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Exec("TRUNCATE TABLE user_level_progresses CASCADE").Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Exec("TRUNCATE TABLE user_badges CASCADE").Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
