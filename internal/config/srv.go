package config

import (
	"event-processing-pipeline/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Engine() *gin.Engine {
	return gin.Default()
}

func Routers(router *gin.Engine) *gin.Engine {
	db := NewMySQLDB()
	eventController := api.NewEventController(db)

	router.POST("/events", eventController.HandleSingleEvent)
	router.POST("/events/batch", eventController.HandleEventsBatch)
	router.GET("/metrics", eventController.GetMetrics)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}
