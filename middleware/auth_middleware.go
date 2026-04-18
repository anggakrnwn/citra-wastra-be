package middleware

import (
	"citra-wastra-be/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims["user_id"])
		ctx.Set("role", claims["role"])

		ctx.Next()
	}
}
