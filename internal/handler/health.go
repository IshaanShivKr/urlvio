package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Pinger interface {
	Ping(context.Context) error
}

type healthResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type HealthHandler struct {
	db Pinger
}

func NewHealthHandler(db Pinger) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, healthResponse{Status: "ok"})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	if err := h.db.Ping(c.Request.Context()); err != nil {
		slog.Warn("readiness check failed", "error", err)
		c.JSON(http.StatusServiceUnavailable, healthResponse{
			Status: "unavailable",
			Error:  "database unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, healthResponse{Status: "ok"})
}
