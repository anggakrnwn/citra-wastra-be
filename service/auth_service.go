package service

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/models"
	"citra-wastra-be/repository"
	"citra-wastra-be/utils"
	"errors"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (dto.UserResponse, error)
	Login(req dto.LoginRequest) (dto.LoginResponse, error)
	GetProfile(userID string) (dto.UserResponse, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(req dto.RegisterRequest) (dto.UserResponse, error) {

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	existingUser, err := s.repo.FindByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserResponse{}, err
	}
	if existingUser != nil && existingUser.ID != "" {
		return dto.UserResponse{}, ErrEmailTaken
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return dto.UserResponse{}, err
	}
	user.Password = hashedPassword

	err = s.repo.Create(&user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (dto.LoginResponse, error) {

	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Token: token,
		Data: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

func (s *authService) GetProfile(userID string) (dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}
