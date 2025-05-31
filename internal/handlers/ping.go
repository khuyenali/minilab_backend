package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping handles the ping endpoint
// @Summary Ping
// @Description Ping endpoint to test API connectivity
// @Tags ping
// @Produce json
// @Success 200 {object} map[string]interface{} "pong response"
// @Router /api/v1/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"status":  "success",
	})
} 