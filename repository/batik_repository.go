package repository

import (
	"citra-wastra-be/models"

	"gorm.io/gorm"
)

type BatikRepository interface {
	Create(batik *models.Batik) error
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
