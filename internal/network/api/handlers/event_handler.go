package handlers

import (
	"net/http"
	"strconv"

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

	c.JSON(http.StatusOK, event)
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

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) GetUserEvents(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	page, limit := getPaginationParams(c)

	events, totalCount, err := h.eventService.GetUserEvents(userID, page, limit)
	if err != nil {
		HandleError(c, domain.ErrFailedToGetUserEvents)
	}

	c.JSON(http.StatusOK, gin.H{
		"events":      events,
		"total_count": totalCount,
		"page":        page,
		"limit":       limit,
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

	c.JSON(http.StatusOK, event)
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

	events, totalCount, err := h.eventService.AdminGetAllEvents(page, limit)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events":      events,
		"total_count": totalCount,
		"page":        page,
		"limit":       limit,
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

	c.JSON(http.StatusOK, event)
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
