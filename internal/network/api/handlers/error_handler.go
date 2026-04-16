package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
)

func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrEmailExists), errors.Is(err, domain.ErrUsernameExists):
		statusCode := http.StatusBadRequest
		if errors.Is(err, domain.ErrEmailExists) {
			statusCode = http.StatusConflict
		}
		if errors.Is(err, domain.ErrUsernameExists) {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   err.Error(),
			Code:    statusCode,
			Details: "Please check your input and try again",
		})

	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   err.Error(),
			Code:    http.StatusUnauthorized,
			Details: "Email or password is incorrect",
		})

	case errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Details: "User not authenticated",
		})

	case errors.Is(err, domain.ErrUserNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "User not found",
			Code:    http.StatusNotFound,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrNoPermissionUpdateUser):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Failed to update user",
			Code:    http.StatusForbidden,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrNoPermissionDeleteUser):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Failed to delete user",
			Code:    http.StatusForbidden,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrFailedToCreateTask):
		statusCode := http.StatusBadRequest
		if errors.Is(err, domain.ErrUserNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to create task",
			Code:    statusCode,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrNoPermissionViewTask):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Failed to get task",
			Code:    http.StatusForbidden,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrNoPermissionUpdateTask):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Failed to update task",
			Code:    http.StatusForbidden,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrNoPermissionDeleteTask):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Failed to delete task",
			Code:    http.StatusForbidden,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrFailedToCreateEvent):
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to create event",
			Code:    statusCode,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrNoPermissionViewEvent):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Failed to get event",
			Code:    http.StatusForbidden,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrNoPermissionUpdateEvent):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Failed to update event",
			Code:    http.StatusForbidden,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrNoPermissionDeleteEvent):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Failed to delete event",
			Code:    http.StatusForbidden,
			Details: err.Error(),
		})

	case errors.Is(err, domain.ErrInvalidUUID):
		resource := detectResourceType(c)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid " + resource + " ID",
			Code:    http.StatusBadRequest,
			Details: resource + " ID must be a valid UUID",
		})

	case errors.Is(err, domain.ErrForbidden):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Forbidden",
			Code:    http.StatusForbidden,
			Details: err.Error(),
		})

	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal server error",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
	}
}

func detectResourceType(c *gin.Context) string {
	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}

	switch {
	case contains(path, "task"):
		return "task"
	case contains(path, "event"):
		return "event"
	case contains(path, "user"):
		return "user"
	default:
		return "resource"
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
