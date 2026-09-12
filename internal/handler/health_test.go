package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockPinger struct {
	pingFunc func(context.Context) error
}

func (m *mockPinger) Ping(ctx context.Context) error {
	if m.pingFunc != nil {
		return m.pingFunc(ctx)
	}
	return nil
}

func TestHealthHandler_Live(t *testing.T) {
	handler := NewHealthHandler(&mockPinger{})

	router := gin.New()
	router.GET("/healthz/live", handler.Live)

	req := httptest.NewRequest(http.MethodGet, "/healthz/live", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := `{"status":"ok"}`
	if rec.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestHealthHandler_Ready(t *testing.T) {
	handler := NewHealthHandler(&mockPinger{
		pingFunc: func(ctx context.Context) error {
			return nil
		},
	})

	router := gin.New()
	router.GET("/healthz/ready", handler.Ready)

	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := `{"status":"ok"}`
	if rec.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestHealthHandler_Ready_DatabaseUnavailable(t *testing.T) {
	handler := NewHealthHandler(&mockPinger{
		pingFunc: func(ctx context.Context) error {
			return errors.New("database unavailable")
		},
	})

	router := gin.New()
	router.GET("/healthz/ready", handler.Ready)

	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusServiceUnavailable,
			rec.Code,
		)
	}

	expected := `{"status":"unavailable","error":"database unavailable"}`
	if rec.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rec.Body.String())
	}
}
