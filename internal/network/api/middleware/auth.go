package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/auth"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/google/uuid"
)

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Authentication required",
				Code:    http.StatusUnauthorized,
				Details: "Authorization header is missing",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Invalid authorization format",
				Code:    http.StatusUnauthorized,
				Details: "Format must be: Bearer <token>",
			})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Invalid or expiered token",
				Code:    http.StatusUnauthorized,
				Details: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), "user_id", claims.UserID),
		)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.UUID{}, false
	}

	userUUID, ok := userID.(uuid.UUID)
	return userUUID, ok
}
