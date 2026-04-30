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
	return r.db.Save(config).Error
}

func (r *systemRepository) GetAllConfigs() ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	err := r.db.Find(&configs).Error
	return configs, err
}

func (r *systemRepository) DeleteConfig(key string) error {
	return r.db.Delete(&models.SystemConfig{}, "key = ?", key).Error
}
