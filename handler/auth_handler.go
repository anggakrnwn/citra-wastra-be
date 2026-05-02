package handler

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/middleware"
	"citra-wastra-be/service"
	"citra-wastra-be/utils"
	"errors"
	"log"
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
			"success": false,
			"error":   "validation failed",
			"details": formattedErrors,
		})
		return
	}

	loginResponse, err := h.service.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailTaken):
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "email already registered",
			})
		case errors.Is(err, service.ErrUsernameTaken):
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "username already taken",
			})
		default:
			log.Printf("auth register error (email=%s): %v", req.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "registration failed: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "user created successfully",
		"data":    loginResponse,
	})

}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req struct {
		IDToken string `json:"id_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "id_token is required",
		})
		return
	}

	loginResponse, err := h.service.GoogleLogin(c.Request.Context(), req.IDToken)
	if err != nil {
		log.Printf("google login error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "google login failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login successful",
		"data":    loginResponse,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		formattedErrors := utils.FormatValidationError(err, &input)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "validation failed",
			"details": formattedErrors,
		})
		return
	}

	loginResponse, err := h.service.Login(input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "invalid email or password",
			})
		default:
			log.Printf("auth login error (email=%s): %v", input.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "login failed",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login successful",
		"data":    loginResponse,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	val, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	userID, ok := val.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "invalid user id",
		})
		return
	}

	userResponse, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "profile retrieved successfully",
		"data":    userResponse,
	})
}
