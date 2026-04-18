package service

import (
	"citra-wastra-be/models"
	"citra-wastra-be/repository"
	"citra-wastra-be/utils"
	"errors"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(user *models.User) error
	Login(email, password string) (models.User, string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(user *models.User) error {

	existingUser, err := s.repo.FindByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil && existingUser.ID != "" {
		return ErrEmailTaken
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.repo.Create(user)
}

func (s *authService) Login(email, password string) (models.User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return models.User{}, "", ErrInvalidConfig
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return models.User{}, "", ErrInvalidConfig
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return models.User{}, "", err
	}

	return *user, token, nil
}
