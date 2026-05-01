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
	GetAllUsers(limit, offset int) ([]models.User, int64, error)
	UpdateUserStatus(userID string, isBanned, isActive bool) error
	UpdateUserRole(userID string, role string) error
	DeleteUser(userID string) error
	FindByOAuthID(provider, oauthID string) (*models.User, error)
	Update(user *models.User) error
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

func (r *userRepository) GetAllUsers(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := r.db.Model(&models.User{})
	db.Count(&total)

	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *userRepository) UpdateUserStatus(userID string, isBanned, isActive bool) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_banned": isBanned,
		"is_active": isActive,
	}).Error
}

func (r *userRepository) UpdateUserRole(userID string, role string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

func (r *userRepository) DeleteUser(userID string) error {
	return r.db.Unscoped().Delete(&models.User{}, "id = ?", userID).Error
}

func (r *userRepository) FindByOAuthID(provider, oauthID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("oauth_provider = ? AND oauth_id = ?", provider, oauthID).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}
