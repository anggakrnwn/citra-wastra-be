package handler

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/middleware"
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

	userResponse, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully!",
		"data":    userResponse,
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

	loginResponse, err := h.service.Login(input)
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"data":    loginResponse,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(string)

	userResponse, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User profile retrieved successfully",
		"data":    userResponse,
	})
}
