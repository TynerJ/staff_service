package handler

import "github.com/gin-gonic/gin"

// HealthHandler handles health of the HTTP server
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck checks if the server is up
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "OK",
	})
}
