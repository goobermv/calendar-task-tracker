package domain

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
)

func HandleValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, dto.ErrorResponse{
		Error:   "Validation failed",
		Code:    http.StatusBadRequest,
		Details: err.Error(),
	})
}

func HandleAlreadyExistsError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "user with this email already exists" {
		statusCode = http.StatusConflict
	}
	if err.Error() == "user with this username already exists" {
		statusCode = http.StatusConflict
	}

	c.JSON(statusCode, dto.ErrorResponse{
		Error:   err.Error(),
		Code:    statusCode,
		Details: "Please check your input and try again",
	})
}

func HandleIncorrectInfoError(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
		Error:   err.Error(),
		Code:    http.StatusUnauthorized,
		Details: "Email or password is incorrect",
	})
}

func HandleUnathorizedError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
		Error:   "Unauthorized",
		Code:    http.StatusUnauthorized,
		Details: "User not authenticated",
	})
}

func HandleUserNotFoundError(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, dto.ErrorResponse{
		Error:   "User not found",
		Code:    http.StatusNotFound,
		Details: err.Error(),
	})
}

func HandleUpdateUserError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "user not found" {
		statusCode = http.StatusNotFound
	} else if err.Error() == "you do not have the permissions to update this user" {
		statusCode = http.StatusForbidden
	}
	c.JSON(statusCode, dto.ErrorResponse{
		Error:   "Failed to update user",
		Code:    statusCode,
		Details: err.Error(),
	})
}

func HandleDeleteUserError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "user not found" {
		statusCode = http.StatusNotFound
	} else if err.Error() == "you do not have the permissions to delete this user" {
		statusCode = http.StatusForbidden
	}
	c.JSON(statusCode, dto.ErrorResponse{
		Error:   "Failed to delete user",
		Code:    statusCode,
		Details: err.Error(),
	})
}

func HandleCreateTaskError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "user not found" {
		statusCode = http.StatusNotFound
	}
	c.JSON(statusCode, dto.ErrorResponse{
		Error:   "Failed to create task",
		Code:    statusCode,
		Details: err.Error(),
	})
}

func HandleInvalidTaskIDError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, dto.ErrorResponse{
		Error:   "Invalid task ID",
		Code:    http.StatusBadRequest,
		Details: "Task ID must be a valid UUID",
	})
}

func HandleTaskError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "task not found" {
		statusCode = http.StatusNotFound
	} else if err.Error() == "you don't have permission to view this task" {
		statusCode = http.StatusForbidden
	}
	c.JSON(statusCode, dto.ErrorResponse{
		Error:   "Failed to get task",
		Code:    statusCode,
		Details: err.Error(),
	})
}

func HandleUpdateTaskError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "task not found" {
		statusCode = http.StatusNotFound
	} else if err.Error() == "you do not have the permissions to update this task" {
		statusCode = http.StatusForbidden
	}
	c.JSON(statusCode, dto.ErrorResponse{
		Error:   "Failed to update task",
		Code:    statusCode,
		Details: err.Error(),
	})
}

func HandleDeleteTaskError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "task not found" {
		statusCode = http.StatusNotFound
	} else if err.Error() == "you do not have the permissions to delete this task" {
		statusCode = http.StatusForbidden
	}
	c.JSON(statusCode, dto.ErrorResponse{
		Error:   "Failed to delete task",
		Code:    statusCode,
		Details: err.Error(),
	})
}

func HandleCreateEventError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "user not found" {
		statusCode = http.StatusNotFound
	}
	c.JSON(statusCode, dto.ErrorResponse{
		Error:   "Failed to create event",
		Code:    statusCode,
		Details: err.Error(),
	})
}

func HandleInvalidEventIDError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, dto.ErrorResponse{
		Error:   "Invalid event ID",
		Code:    http.StatusBadRequest,
		Details: "Event ID must be a valid UUID",
	})
}

func HandleEventError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "event not found" {
		statusCode = http.StatusNotFound
	} else if err.Error() == "you don't have permission to view this event" {
		statusCode = http.StatusForbidden
	}
	c.JSON(statusCode, dto.ErrorResponse{
		Error:   "Failed to get event",
		Code:    statusCode,
		Details: err.Error(),
	})
}

func HandleUpdateEventError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "event not found" {
		statusCode = http.StatusNotFound
	} else if err.Error() == "you do not have the permissions to update this event" {
		statusCode = http.StatusForbidden
	}
	c.JSON(statusCode, dto.ErrorResponse{
		Error:   "Failed to update event",
		Code:    statusCode,
		Details: err.Error(),
	})
}

func HandleDeleteEventError(c *gin.Context, err error) {
	statusCode := http.StatusBadRequest
	if err.Error() == "event not found" {
		statusCode = http.StatusNotFound
	} else if err.Error() == "you do not have the permissions to delete this event" {
		statusCode = http.StatusForbidden
	}
	c.JSON(statusCode, dto.ErrorResponse{
		Error:   "Failed to delete event",
		Code:    statusCode,
		Details: err.Error(),
	})
}
