package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/middleware"
	taskUsecase "github.com/goobermv/calendar-task-tracker/internal/usescases/task"
	"github.com/google/uuid"
)

type TaskHandler struct {
	taskService *taskUsecase.Service
}

func NewTaskHandler(taskService *taskUsecase.Service) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userIDStr, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Details: "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	useCaseReq := taskUsecase.CreateTaskRequest{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
	}

	task, err := h.taskService.CreateTask(useCaseReq)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to create task",
			Code:    statusCode,
			Details: err.Error(),
		})
		return
	}

	response := dto.TaskResponse{
		ID:          task.ID.String(),
		UserID:      task.UserID.String(),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userIDStr, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Details: "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid task ID",
			Code:    http.StatusBadRequest,
			Details: "Task ID must be a valid UUID",
		})
		return
	}

	task, err := h.taskService.GetTaskByID(taskID, userID)
	if err != nil {
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
		return
	}

	response := dto.TaskResponse{
		ID:          task.ID.String(),
		UserID:      task.UserID.String(),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

//implement GetUserTasks function handler

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userIDStr, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Details: "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid task ID",
			Code:    http.StatusBadRequest,
			Details: "Task ID must be a valid UUID",
		})
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	useCaseReq := taskUsecase.UpdateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	task, err := h.taskService.UpdateTask(taskID, userID, useCaseReq)
	if err != nil {
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
		return
	}

	resonse := dto.TaskResponse{
		ID:          task.ID.String(),
		UserID:      task.UserID.String(),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	c.JSON(http.StatusOK, resonse)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userIDStr, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Details: "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid task ID",
			Code:    http.StatusBadRequest,
			Details: "Task ID must be a valid UUID",
		})
		return
	}

	err = h.taskService.DeleteTask(taskID, userID)
	if err != nil {
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
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

/*
func (h *TaskHandler) CompleteTask(c *gin.Context) {
    // 1. Get authenticated user ID
    userIDStr, exists := middleware.GetUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
            Error:   "Unauthorized",
            Code:    http.StatusUnauthorized,
            Details: "User not authenticated",
        })
        return
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "Invalid user ID",
            Code:    http.StatusBadRequest,
            Details: err.Error(),
        })
        return
    }

    // 2. Parse task ID from URL
    taskIDStr := c.Param("id")
    taskID, err := uuid.Parse(taskIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "Invalid task ID",
            Code:    http.StatusBadRequest,
            Details: "Task ID must be a valid UUID",
        })
        return
    }

    // 3. Execute use case
    task, err := h.taskService.CompleteTask(taskID, userID)
    if err != nil {
        statusCode := http.StatusBadRequest
        if err.Error() == "task not found" {
            statusCode = http.StatusNotFound
        } else if err.Error() == "you don't have permission to update this task" {
            statusCode = http.StatusForbidden
        }
        c.JSON(statusCode, dto.ErrorResponse{
            Error:   "Failed to complete task",
            Code:    statusCode,
            Details: err.Error(),
        })
        return
    }

    // 4. Convert to response DTO
    response := dto.TaskResponse{
        ID:          task.ID.String(),
        UserID:      task.UserID.String(),
        Title:       task.Title,
        Description: task.Description,
        Status:      task.Status,
        DueDate:     task.DueDate,
        Priority:    task.Priority,
        CreatedAt:   task.CreatedAt,
        UpdatedAt:   task.UpdatedAt,
    }

    // 5. Send response
    c.JSON(http.StatusOK, response)
}
*/
