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

	errInvalidTokenPayload = "invalid token payload"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			abortWithError(ctx, http.StatusUnauthorized, "authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			abortWithError(ctx, http.StatusBadRequest, "invalid authorization format")
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			abortWithError(ctx, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			abortWithError(ctx, http.StatusUnauthorized, errInvalidTokenPayload)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			abortWithError(ctx, http.StatusUnauthorized, errInvalidTokenPayload)
			return
		}

		ctx.Set(UserIDKey, userID)
		ctx.Set(RoleKey, role)

		ctx.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get(RoleKey)
		if !exists || role.(string) != "admin" {
			abortWithError(ctx, http.StatusForbidden, "access denied: admin role required")
			return
		}
		ctx.Next()
	}
}

func abortWithError(ctx *gin.Context, status int, msg string) {
	ctx.JSON(status, gin.H{
		"success": false,
		"error":   msg,
	})
	ctx.Abort()
}
