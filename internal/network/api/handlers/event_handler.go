package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/domain"
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
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	var req dto.CreateEventRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		HandleError(c, err)
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
		HandleError(c, err)
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
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	eventIDstr := c.Param("id")
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {

		return
	}

	event, err := h.eventService.GetEventByID(eventID, userID)
	if err != nil {
		HandleError(c, err)
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
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	eventIDstr := c.Param("id")
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		HandleError(c, err)
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	event, err := h.eventService.UpdateEvent(eventID, userID, req)
	if err != nil {
		HandleError(c, err)
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
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = h.eventService.DeleteEvent(eventID, userID)
	if err != nil {

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// implement other handlers for functions
