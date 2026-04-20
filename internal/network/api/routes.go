package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/auth"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/handlers"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/middleware"
)

func SetupRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	taskHandler *handlers.TaskHandler,
	eventHandler *handlers.EventHandler,
	jwtService *auth.JWTService,
) {
	v1 := router.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", userHandler.Register)
			authGroup.POST("/login", userHandler.Login)
		}
	}

	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		protected.GET("/users/me", userHandler.GetProfile)
		protected.PATCH("/users/:id", userHandler.UpdateUserInfo)
		protected.PATCH("/users/:id", userHandler.UpdateUserPassword)
		protected.DELETE("/users/:id", userHandler.DeleteUser)

		protected.GET("/tasks", taskHandler.GetTask)
		protected.GET("/tasks", taskHandler.GetUserTasks)
		protected.POST("/tasks", taskHandler.CreateTask)
		protected.PUT("/tasks/:id", taskHandler.UpdateTask)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)

		protected.GET("/events", eventHandler.GetEvent)
		protected.GET("/events", eventHandler.GetUserEvents)
		protected.POST("/events", eventHandler.CreateEvent)
		protected.PUT("/events/:id", eventHandler.UpdateEvent)
		protected.DELETE("/events/:id", eventHandler.DeleteEvent)

		// Calendar routes (to be implemented)
		// protected.GET("/calendar", calendarHandler.GetCalendar)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})
}
