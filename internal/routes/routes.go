package routes

import (
	"github.com/IshaanShivKr/urlvio/internal/handler"
	"github.com/IshaanShivKr/urlvio/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Register(
	router *gin.Engine,
	healthHandler *handler.HealthHandler,
	linkHandler *handler.LinkHandler,
	authMiddleware gin.HandlerFunc,
) {
	healthz := router.Group("/healthz")
	{
		healthz.GET("/live", healthHandler.Live)
		healthz.GET("/ready", healthHandler.Ready)
	}

	api := router.Group("/api/v1")
	api.Use(middleware.RateLimiter())
	{
		urls := api.Group("/urls")
		urls.Use(authMiddleware)
		{
			urls.POST("", linkHandler.Create)
			urls.GET("/:shortCode/stats", linkHandler.Stats)
			urls.GET("/:shortCode", linkHandler.Get)
			urls.DELETE("/:shortCode", linkHandler.Delete)
			urls.PUT("/:shortCode", linkHandler.Update)
		}
	}

	router.StaticFile("/docs/api", "docs/index.html")
	router.StaticFile("/docs/openapi.yaml", "docs/openapi.yaml")

	router.GET("/:shortCode", linkHandler.Redirect)
}
