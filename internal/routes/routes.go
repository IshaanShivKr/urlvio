package routes

import (
	"github.com/IshaanShivKr/urlvio/internal/handler"
	"github.com/gin-gonic/gin"
)

func Register(
	router *gin.Engine,
	healthHandler *handler.HealthHandler,
	linkHandler *handler.LinkHandler,
) {
	healthz := router.Group("/healthz")
	{
		healthz.GET("/live", healthHandler.Live)
		healthz.GET("/ready", healthHandler.Ready)
	}

	api := router.Group("/api/v1")
	{
		urls := api.Group("/urls")
		{
			urls.POST("", linkHandler.Create)
			urls.GET("/:shortCode/stats", linkHandler.Stats)
			urls.GET("/:shortCode", linkHandler.Get)
			urls.DELETE("/:shortCode", linkHandler.Delete)
		}
	}
}
