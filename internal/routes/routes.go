package routes

import (
	"github.com/IshaanShivKr/urlvio/internal/handler"
	"github.com/gin-gonic/gin"
)

func Register(
	router *gin.Engine,
	healthHandler *handler.HealthHandler,
) {
	healthz := router.Group("/healthz")
	{
		healthz.GET("/live", healthHandler.Live)
		healthz.GET("/ready", healthHandler.Ready)
	}
}
