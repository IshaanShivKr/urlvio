package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/IshaanShivKr/urlvio/internal/service"
	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	service *service.LinkService
	baseURL string
}

func NewLinkHandler(service *service.LinkService, baseURL string) *LinkHandler {
	return &LinkHandler{
		service: service,
		baseURL: baseURL,
	}
}

func (h *LinkHandler) Create(c *gin.Context) {
	var req CreateLinkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	link, err := h.service.Create(c.Request.Context(), req.URL)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrURLRequired),
			errors.Is(err, service.ErrInvalidURL):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		default:
			slog.Error("failed to create link", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, newLinkResponse(
		h.baseURL,
		link.Code,
		link.URL,
		link.CreatedAt,
	))
}

func newLinkResponse(baseURL, code, rawURL string, createdAt time.Time) CreateLinkResponse {
	return CreateLinkResponse{
		Code:      code,
		ShortURL:  strings.TrimRight(baseURL, "/") + "/" + code,
		URL:       rawURL,
		CreatedAt: createdAt,
	}
}
