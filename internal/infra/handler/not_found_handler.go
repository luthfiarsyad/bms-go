package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
}