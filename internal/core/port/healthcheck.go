package port

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthCheckResponse struct {
	TimeStamp string `json:"timestamp"`
	Version   string `json:"version"`
}

func HealthCheck(c *gin.Context) {
	health := HealthCheckResponse{
		TimeStamp: time.Now().String(),
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, health)
}
