package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/domain"
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
		c.Set("user_type", claims.UserType)

		ctx := context.WithValue(c.Request.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_type", claims.UserType)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")

		if !exists {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Unauthorized",
				Code:    http.StatusUnauthorized,
				Details: "User session information not found",
			})
			c.Abort()
			return
		}

		if userType != domain.UserTypeAdmin {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "Access denied",
				Code:    http.StatusForbidden,
				Details: "Admin privileges required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.UUID{}, false
	}

	if userUUID, ok := userID.(uuid.UUID); ok {
		return userUUID, true
	}

	if userStr, ok := userID.(string); ok {
		parsedUUID, err := uuid.Parse(userStr)
		if err == nil {
			return parsedUUID, true
		}
	}

	return uuid.UUID{}, false
}

func GetID(c *gin.Context) (uuid.UUID, bool) {
	ID, exists := c.Get("id")
	if !exists {
		idParam := c.Param("id")
		if idParam != "" {
			parsed, err := uuid.Parse(idParam)
			return parsed, err == nil
		}
		return uuid.UUID{}, false
	}

	if UUID, ok := ID.(uuid.UUID); ok {
		return UUID, true
	}

	if idStr, ok := ID.(string); ok {
		parsedUUID, err := uuid.Parse(idStr)
		if err == nil {
			return parsedUUID, true
		}
	}

	return uuid.UUID{}, false
}
