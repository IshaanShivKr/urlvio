package routes

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IshaanShivKr/urlvio/internal/handler"
	"github.com/IshaanShivKr/urlvio/internal/model"
	"github.com/IshaanShivKr/urlvio/internal/service"
	"github.com/gin-gonic/gin"
)

const burstSize = 20

type stubPinger struct{}

func (stubPinger) Ping(ctx context.Context) error { return nil }

type stubLinkRepository struct{}

func (stubLinkRepository) Create(ctx context.Context, link *model.Link) error {
	return nil
}

func (stubLinkRepository) Get(ctx context.Context, code string) (*model.Link, error) {
	return &model.Link{Code: code, URL: "https://example.com"}, nil
}

func (stubLinkRepository) Delete(ctx context.Context, code string) error {
	return nil
}

func (stubLinkRepository) Update(ctx context.Context, code, rawURL string) (*model.Link, error) {
	return &model.Link{Code: code, URL: rawURL}, nil
}

func (stubLinkRepository) GetAndIncrement(ctx context.Context, code string) (*model.Link, error) {
	return &model.Link{Code: code, URL: "https://example.com"}, nil
}

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	healthHandler := handler.NewHealthHandler(stubPinger{})
	linkHandler := handler.NewLinkHandler(
		service.NewLinkService(stubLinkRepository{}),
		"http://localhost:8080",
	)

	router := gin.New()
	Register(router, healthHandler, linkHandler)

	return router
}

func doRequest(router *gin.Engine, method, path, body, remoteAddr string) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if remoteAddr != "" {
		req.RemoteAddr = remoteAddr
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func TestRegister_RouteReachability(t *testing.T) {
	router := newTestRouter()

	cases := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{"live", http.MethodGet, "/healthz/live", "", http.StatusOK},
		{"ready", http.MethodGet, "/healthz/ready", "", http.StatusOK},
		{"create", http.MethodPost, "/api/v1/urls", `{"url":"https://example.com"}`, http.StatusCreated},
		{"stats", http.MethodGet, "/api/v1/urls/abc123/stats", "", http.StatusOK},
		{"get", http.MethodGet, "/api/v1/urls/abc123", "", http.StatusOK},
		{"delete", http.MethodDelete, "/api/v1/urls/abc123", "", http.StatusNoContent},
		{"update", http.MethodPut, "/api/v1/urls/abc123", `{"url":"https://example.org"}`, http.StatusOK},
		{"redirect", http.MethodGet, "/abc123", "", http.StatusFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := doRequest(router, tc.method, tc.path, tc.body, "")

			if rec.Code != tc.wantStatus {
				t.Fatalf(
					"%s %s: expected status %d, got %d",
					tc.method, tc.path, tc.wantStatus, rec.Code,
				)
			}
		})
	}
}

func TestRegister_RateLimiterAppliedToAPIGroup(t *testing.T) {
	router := newTestRouter()

	const ip = "203.0.113.10:12345"

	for i := 0; i < burstSize; i++ {
		rec := doRequest(router, http.MethodGet, "/api/v1/urls/abc123", "", ip)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected status %d, got %d", i+1, http.StatusOK, rec.Code)
		}
	}

	rec := doRequest(router, http.MethodGet, "/api/v1/urls/abc123", "", ip)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf(
			"expected /api/v1/urls to be rate limited after %d requests, got status %d",
			burstSize, rec.Code,
		)
	}
}

func TestRegister_RateLimiterNotAppliedToHealthz(t *testing.T) {
	router := newTestRouter()

	const ip = "203.0.113.11:12345"

	for i := 0; i < burstSize+1; i++ {
		rec := doRequest(router, http.MethodGet, "/healthz/live", "", ip)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"request %d: expected /healthz/live to never be rate limited, got status %d",
				i+1, rec.Code,
			)
		}
	}
}

func TestRegister_RateLimiterNotAppliedToRedirect(t *testing.T) {
	router := newTestRouter()

	const ip = "203.0.113.12:12345"

	for i := 0; i < burstSize+1; i++ {
		rec := doRequest(router, http.MethodGet, "/abc123", "", ip)

		if rec.Code != http.StatusFound {
			t.Fatalf(
				"request %d: expected redirect to never be rate limited, got status %d",
				i+1, rec.Code,
			)
		}
	}
}
