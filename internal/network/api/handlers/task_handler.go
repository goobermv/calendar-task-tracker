package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/middleware"
	taskUsecase "github.com/goobermv/calendar-task-tracker/internal/usescases/task"
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
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	useCaseReq := dto.CreateTaskRequest{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
	}

	task, err := h.taskService.CreateTask(useCaseReq)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	taskID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrTaskNotFound)
		return
	}

	task, err := h.taskService.GetTaskByID(taskID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func getPaginationParams(c *gin.Context) (int, int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return page, limit
}

func (h *TaskHandler) GetUserTasks(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	page, limit := getPaginationParams(c)

	tasks, totalCount, err := h.taskService.GetUserTasks(userID, page, limit)
	if err != nil {
		HandleError(c, domain.ErrFailedToGetUserTasks)
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":       tasks,
		"total_count": totalCount,
		"page":        page,
		"limit":       limit,
	})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	taskID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrTaskNotFound)
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	task, err := h.taskService.UpdateTask(taskID, userID, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	taskID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrTaskNotFound)
		return
	}

	err := h.taskService.DeleteTask(taskID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *TaskHandler) AdminGetAllTasks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	tasks, totalCount, err := h.taskService.AdminGetAllTasks(page, limit)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":       tasks,
		"total_count": totalCount,
		"page":        page,
		"limit":       limit,
	})
}

func (h *TaskHandler) AdminUpdateTask(c *gin.Context) {
	taskID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrTaskNotFound)
		return
	}

	var req dto.AdminUpdateTaskRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	task, err := h.taskService.AdminUpdateTask(taskID, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) AdminDeleteTask(c *gin.Context) {
	taskID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrTaskNotFound)
		return
	}

	err := h.taskService.AdminDeleteTask(taskID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
