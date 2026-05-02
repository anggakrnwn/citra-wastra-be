package service

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/models"
	"citra-wastra-be/repository"
	"citra-wastra-be/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (dto.LoginResponse, error)
	Login(req dto.LoginRequest) (dto.LoginResponse, error)
	GetProfile(userID string) (dto.UserResponse, error)
	GoogleLogin(ctx context.Context, idToken string) (dto.LoginResponse, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(req dto.RegisterRequest) (dto.LoginResponse, error) {
	user := models.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   req.Password,
		Role:       "user",
		IsVerified: true,
	}

	existingUser, err := s.repo.FindByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.LoginResponse{}, err
	}
	if existingUser != nil && existingUser.ID != "" {
		return dto.LoginResponse{}, ErrEmailTaken
	}

	existingUser, err = s.repo.FindByUsername(user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.LoginResponse{}, err
	}
	if existingUser != nil && existingUser.ID != "" {
		return dto.LoginResponse{}, ErrUsernameTaken
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	user.Password = hashedPassword

	err = s.repo.Create(&user)
	if err != nil {
		return dto.LoginResponse{}, err
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
			XP:       user.XP,
			Level:    user.XP/1000 + 1,
		},
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (dto.LoginResponse, error) {

	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	if user.OAuthProvider != "" {
		return dto.LoginResponse{}, fmt.Errorf("this account is linked to %s, please login with %s", user.OAuthProvider, user.OAuthProvider)
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
			XP:       user.XP,
			Level:    user.XP/1000 + 1,
		},
	}, nil
}

func (s *authService) GoogleLogin(ctx context.Context, tokenStr string) (dto.LoginResponse, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		return dto.LoginResponse{}, fmt.Errorf("GOOGLE_CLIENT_ID is not configured")
	}

	payload, err := idtoken.Validate(ctx, tokenStr, clientID)
	if err != nil {
		log.Printf("[OAUTH] Token validation failed: %v", err)
		return dto.LoginResponse{}, fmt.Errorf("invalid google token")
	}

	email, _ := payload.Claims["email"].(string)
	oauthID, _ := payload.Claims["sub"].(string)
	name, _ := payload.Claims["name"].(string)

	user, err := s.repo.FindByOAuthID("google", oauthID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.LoginResponse{}, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || user.ID == "" {
		existing, _ := s.repo.FindByEmail(email)
		if existing != nil && existing.ID != "" {
			existing.OAuthProvider = "google"
			existing.OAuthID = oauthID
			existing.IsVerified = true
			user = existing
			err = s.repo.Update(user)
		} else {
			baseUsername := name
			if baseUsername == "" {
				baseUsername = "user_" + oauthID[:8]
			}

			user = &models.User{
				Username:      baseUsername,
				Email:         email,
				OAuthProvider: "google",
				OAuthID:       oauthID,
				IsVerified:    true,
			}
			err = s.repo.Create(user)
		}

		if err != nil {
			log.Printf("failed to save/update user: %v", err)
			return dto.LoginResponse{}, fmt.Errorf("failed to process user data")
		}
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
			XP:       user.XP,
			Level:    user.XP/1000 + 1,
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
		XP:       user.XP,
		Level:    user.XP/1000 + 1,
	}, nil
}
