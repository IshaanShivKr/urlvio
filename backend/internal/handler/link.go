package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/IshaanShivKr/urlvio/internal/auth"
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
	userID, ok := auth.UserID(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var req CreateLinkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	link, err := h.service.Create(c.Request.Context(), userID, req.URL)
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
	userID, ok := auth.UserID(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	code := c.Param("shortCode")

	link, err := h.service.Get(c.Request.Context(), userID, code)
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

func (h *LinkHandler) List(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	links, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, links)
}

func (h *LinkHandler) Stats(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	code := c.Param("shortCode")

	link, err := h.service.Get(c.Request.Context(), userID, code)
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

func (h *LinkHandler) Delete(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	code := c.Param("shortCode")

	if err := h.service.Delete(c.Request.Context(), userID, code); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "link not found",
			})
			return
		}

		slog.Error("failed to delete link", "code", code, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *LinkHandler) Update(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	code := c.Param("shortCode")

	var req CreateLinkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	link, err := h.service.Update(c.Request.Context(), userID, code, req.URL)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrURLRequired),
			errors.Is(err, service.ErrInvalidURL):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "link not found",
			})

		default:
			slog.Error("failed to update link", "code", code, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
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

func newLinkResponse(baseURL, code, rawURL string, createdAt, updatedAt time.Time) CreateLinkResponse {
	return CreateLinkResponse{
		Code:      code,
		ShortURL:  strings.TrimRight(baseURL, "/") + "/" + code,
		URL:       rawURL,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (h *LinkHandler) Redirect(c *gin.Context) {
	code := c.Param("shortCode")

	link, err := h.service.GetAndIncrement(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "link not found",
			})
			return
		}

		slog.Error("failed to redirect link", "code", code, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.Redirect(http.StatusFound, link.URL)
}
