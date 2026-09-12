package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeaders(t *testing.T) {
	router := gin.New()
	router.Use(SecurityHeaders())

	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	expectedHeaders := map[string]string{
		"X-Frame-Options":       "DENY",
		"X-Content-Type-Options": "nosniff",
		"Referrer-Policy":       "strict-origin",
		"Permissions-Policy":    "geolocation=(), camera=(), microphone=()",
	}

	for header, expected := range expectedHeaders {
		if got := rec.Header().Get(header); got != expected {
			t.Errorf("expected %s=%q, got %q", header, expected, got)
		}
	}
}

func TestSecurityHeaders_AllowsRequest(t *testing.T) {
	router := gin.New()
	router.Use(SecurityHeaders())

	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}
