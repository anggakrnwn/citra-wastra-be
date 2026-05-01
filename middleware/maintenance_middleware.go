package middleware

import (
	"citra-wastra-be/repository"
	"citra-wastra-be/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func MaintenanceMiddleware(repo repository.SystemRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if strings.Contains(path, "/auth/") {
			ctx.Next()
			return
		}

		config, err := repo.GetConfig("MAINTENANCE_MODE")
		if err == nil && config.Value == "true" {
			authHeader := ctx.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					claims, err := utils.ValidateToken(parts[1])
					if err == nil {
						role, _ := claims["role"].(string)
						if role == "admin" || role == "super_admin" {
							ctx.Next()
							return
						}
					}
				}
			}

			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "the system is under maintenance. Please try again later.",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
