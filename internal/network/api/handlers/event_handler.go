package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/middleware"
	eventUsecase "github.com/goobermv/calendar-task-tracker/internal/usescases/event"
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
		ID:          event.ID,
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

	eventID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrEventNotFound)
		return
	}

	event, err := h.eventService.GetEventByID(eventID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.EventResponse{
		ID:          event.ID,
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

func (h *EventHandler) GetUserEvents(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	events, err := h.eventService.GetUserEvents(userID)
	if err != nil {
		HandleError(c, domain.ErrFailedToGetUserEvents)
	}

	responses := make([]dto.EventResponse, len(events))
	for i, event := range events {
		responses[i] = dto.EventResponse{
			ID:          event.ID,
			UserID:      event.UserID,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			EventType:   event.EventType,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"events": responses,
		"count":  len(responses),
	})
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	eventID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrEventNotFound)
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
		ID:          event.ID,
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

	eventID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrEventNotFound)
		return
	}

	err := h.eventService.DeleteEvent(eventID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *EventHandler) AdminGetAllEvents(c *gin.Context) {
	events, err := h.eventService.AdminGetAllEvents()
	if err != nil {
		HandleError(c, err)
		return
	}

	responses := make([]dto.EventResponse, len(events))
	for i, event := range events {
		responses[i] = dto.EventResponse{
			ID:          event.ID,
			UserID:      event.ID,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			EventType:   event.EventType,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"events": responses,
		"count":  len(responses),
	})
}

func (h *EventHandler) AdminUpdateEvents(c *gin.Context) {
	eventID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrEventNotFound)
		return
	}

	var req dto.AdminUpdateEventRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	event, err := h.eventService.AdminUpdateEvent(eventID, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.EventResponse{
		ID:          event.ID,
		UserID:      event.UserID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		EventType:   event.EventType,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	})
}

func (h *EventHandler) AdminDeleteEvents(c *gin.Context) {
	eventID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrEventNotFound)
		return
	}

	err := h.eventService.AdminDeleteEvent(eventID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
