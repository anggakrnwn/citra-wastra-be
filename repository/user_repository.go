package repository

import (
	"citra-wastra-be/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByID(ID string) (*models.User, error)
	UpdateTotalXP(userID string, totalXP int) error
	GetUsersByIDs(ids []string) ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) UpdateTotalXP(userID string, totalXP int) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("xp", totalXP).Error
}

func (r *userRepository) GetUsersByIDs(ids []string) ([]models.User, error) {
	var users []models.User

	err := r.db.Where("id IN ?", ids).Select("id", "username").Find(&users).Error
	return users, err
}
