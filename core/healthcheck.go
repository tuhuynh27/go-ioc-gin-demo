package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheck struct {
	Component struct{} `implements:"HealthCheck"`
}

func (h *HealthCheck) Check(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
