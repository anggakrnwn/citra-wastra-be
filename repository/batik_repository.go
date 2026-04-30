package repository

import (
	"citra-wastra-be/models"

	"gorm.io/gorm"
)

type BatikRepository interface {
	Create(batik *models.Batik) error
	GetAllLogs(limit, offset int) ([]models.Batik, int64, error)

	// Catalog
	CreateCatalog(catalog *models.BatikCatalog) error
	UpdateCatalog(catalog *models.BatikCatalog) error
	DeleteCatalog(id string) error
	GetAllCatalog() []models.BatikCatalog
}

type batikRepository struct {
	db *gorm.DB
}

func NewBatikRepository(db *gorm.DB) BatikRepository {
	return &batikRepository{db}
}

func (r *batikRepository) Create(batik *models.Batik) error {
	return r.db.Create(batik).Error
}

func (r *batikRepository) GetAllLogs(limit, offset int) ([]models.Batik, int64, error) {
	var logs []models.Batik
	var total int64

	db := r.db.Model(&models.Batik{})
	db.Count(&total)

	err := db.Preload("User").Limit(limit).Offset(offset).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

func (r *batikRepository) CreateCatalog(catalog *models.BatikCatalog) error {
	return r.db.Create(catalog).Error
}

func (r *batikRepository) UpdateCatalog(catalog *models.BatikCatalog) error {
	return r.db.Save(catalog).Error
}

func (r *batikRepository) DeleteCatalog(id string) error {
	return r.db.Delete(&models.BatikCatalog{}, "id = ?", id).Error
}

func (r *batikRepository) GetAllCatalog() []models.BatikCatalog {
	var catalogs []models.BatikCatalog
	r.db.Find(&catalogs)
	return catalogs
}
