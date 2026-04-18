package handler

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/models"
	"citra-wastra-be/service"
	"citra-wastra-be/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		formattedErrors := utils.FormatValidationError(err, &req)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": formattedErrors,
		})
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.service.Register(&user); err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	response := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully!",
		"data":    response,
	})

}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		formattedErrors := utils.FormatValidationError(err, &input)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": formattedErrors,
		})
		return
	}

	user, token, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidConfig) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	userResponse := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	response := dto.LoginResponse{
		Token: token,
		Data:  userResponse,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"data":    response,
	})
}
