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
		if errors.Is(err, service.ErrURLRequired) ||
			errors.Is(err, service.ErrInvalidURL) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		slog.Error("failed to create link", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, newLinkResponse(
		h.baseURL,
		link.Code,
		link.URL,
		link.CreatedAt,
		link.UpdatedAt,
	))
}

func (h *LinkHandler) Get(c *gin.Context) {
	code := c.Param("shortCode")

	link, err := h.service.Get(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "link not found",
			})
			return
		}

		slog.Error("failed to get link", "code", code, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, newLinkResponse(
		h.baseURL,
		link.Code,
		link.URL,
		link.CreatedAt,
		link.UpdatedAt,
	))
}

func (h *LinkHandler) Stats(c *gin.Context) {
	code := c.Param("shortCode")

	link, err := h.service.Get(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "link not found",
			})
			return
		}

		slog.Error("failed to get link stats", "code", code, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, LinkStatsResponse{
		Code:        link.Code,
		URL:         link.URL,
		AccessCount: link.AccessCount,
		CreatedAt:   link.CreatedAt,
		UpdatedAt:   link.UpdatedAt,
	})
}

func newLinkResponse(baseURL, code, rawURL string, createdAt, updatedAt time.Time) CreateLinkResponse {
	return CreateLinkResponse{
		Code:      code,
		ShortURL:  strings.TrimRight(baseURL, "/") + "/" + code,
		URL:       rawURL,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
