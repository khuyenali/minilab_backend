package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCorsMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.Use(corsMiddleware())
	
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// Test CORS headers
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestOptionsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.Use(corsMiddleware())
	
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// Test OPTIONS request
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
} 