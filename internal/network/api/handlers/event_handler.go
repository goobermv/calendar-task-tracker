package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/middleware"
	eventUsecase "github.com/goobermv/calendar-task-tracker/internal/usescases/event"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventService *eventUsecase.Service
}

func NewEventHandler(eventService *eventUsecase.Service) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}
func (h *EventHandler) CreateEvent(c *gin.Context) {
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

	var req dto.CreateEventRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
	}

	useCaseReq := dto.CreateEventRequest{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		EventType:   req.EventType,
	}

	event, err := h.eventService.CreateEvent(useCaseReq)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to create event",
			Code:    statusCode,
			Details: err.Error(),
		})
		return
	}

	response := dto.EventResponse{
		ID:          event.ID.String(),
		UserID:      event.UserID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		EventType:   event.EventType,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *EventHandler) GetEvent(c *gin.Context) {
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

	eventIDstr := c.Param("id")
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Code:    http.StatusBadRequest,
			Details: "Event ID must be a valid UUID",
		})
		return
	}

	event, err := h.eventService.GetEventByID(eventID, userID)
	if err != nil {
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
		return
	}

	response := dto.EventResponse{
		ID:          event.ID.String(),
		UserID:      event.UserID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		EventType:   event.EventType,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// implement GetUserEvents handler

func (h *EventHandler) UpdateEvent(c *gin.Context) {
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

	eventIDstr := c.Param("id")
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Code:    http.StatusBadRequest,
			Details: "Event ID must be a valid UUID",
		})
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	useCaseReq := dto.UpdateEventRequest{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		EventType:   req.EventType,
	}

	event, err := h.eventService.UpdateEvent(eventID, userID, useCaseReq)
	if err != nil {
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
		return
	}

	response := dto.EventResponse{
		ID:          event.ID.String(),
		UserID:      event.UserID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		EventType:   event.EventType,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
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

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid task ID",
			Code:    http.StatusBadRequest,
			Details: "Task ID must be a valid UUID",
		})
		return
	}

	err = h.eventService.DeleteEvent(eventID, userID)
	if err != nil {
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
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// implement other handlers for functions
