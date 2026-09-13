package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/IshaanShivKr/urlvio/internal/model"
	"github.com/IshaanShivKr/urlvio/internal/repository"
	"github.com/IshaanShivKr/urlvio/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockLinkRepository struct {
	createFunc          func(context.Context, *model.Link) error
	getFunc             func(context.Context, string, string) (*model.Link, error)
	listFunc            func(context.Context, string) ([]*model.Link, error)
	deleteFunc          func(context.Context, string, string) error
	updateFunc          func(context.Context, string, string, string) (*model.Link, error)
	getAndIncrementFunc func(context.Context, string) (*model.Link, error)
}

func (m *mockLinkRepository) Create(ctx context.Context, link *model.Link) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, link)
	}
	return nil
}

func (m *mockLinkRepository) Get(ctx context.Context, userID, code string) (*model.Link, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, userID, code)
	}
	return nil, repository.ErrNotFound
}

func (m *mockLinkRepository) List(ctx context.Context, userID string) ([]*model.Link, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockLinkRepository) Delete(ctx context.Context, userID, code string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, userID, code)
	}
	return repository.ErrNotFound
}

func (m *mockLinkRepository) Update(ctx context.Context, userID, code, rawURL string) (*model.Link, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, userID, code, rawURL)
	}
	return nil, repository.ErrNotFound
}

func (m *mockLinkRepository) GetAndIncrement(ctx context.Context, code string) (*model.Link, error) {
	if m.getAndIncrementFunc != nil {
		return m.getAndIncrementFunc(ctx, code)
	}
	return nil, repository.ErrNotFound
}

func newTestHandler(repo repository.LinkRepository) *LinkHandler {
	svc := service.NewLinkService(repo)
	return NewLinkHandler(svc, "http://localhost:8080")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		c.Next()
	})

	return router
}

func TestLinkHandler_Create(t *testing.T) {
	createdAt := time.Now().Truncate(time.Microsecond)
	updatedAt := createdAt

	repo := &mockLinkRepository{
		createFunc: func(ctx context.Context, link *model.Link) error {
			link.Code = "abc123"
			link.CreatedAt = createdAt
			link.UpdatedAt = updatedAt
			return nil
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.POST("/api/v1/urls", handler.Create)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/urls",
		strings.NewReader(`{"url":"https://example.com"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var response CreateLinkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code != "abc123" {
		t.Fatalf("expected code abc123, got %q", response.Code)
	}

	if response.ShortURL != "http://localhost:8080/abc123" {
		t.Fatalf("expected short URL %q, got %q", "http://localhost:8080/abc123", response.ShortURL)
	}

	if response.URL != "https://example.com" {
		t.Fatalf("expected URL %q, got %q", "https://example.com", response.URL)
	}

	if !response.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected CreatedAt %v, got %v", createdAt, response.CreatedAt)
	}

	if !response.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("expected UpdatedAt %v, got %v", updatedAt, response.UpdatedAt)
	}
}

func TestLinkHandler_Create_InvalidJSON(t *testing.T) {
	repo := &mockLinkRepository{}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.POST("/api/v1/urls", handler.Create)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/urls",
		strings.NewReader(`{"url":`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestLinkHandler_Create_InvalidURL(t *testing.T) {
	repo := &mockLinkRepository{}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.POST("/api/v1/urls", handler.Create)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/urls",
		strings.NewReader(`{"url":"not-a-url"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestLinkHandler_Create_RepositoryError(t *testing.T) {
	repo := &mockLinkRepository{
		createFunc: func(ctx context.Context, link *model.Link) error {
			return errors.New("database unavailable")
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.POST("/api/v1/urls", handler.Create)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/urls",
		strings.NewReader(`{"url":"https://example.com"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestLinkHandler_Get(t *testing.T) {
	link := &model.Link{
		ID:          uuid.New(),
		URL:         "https://example.com",
		Code:        "abc123",
		CreatedAt:   time.Now().Truncate(time.Microsecond),
		UpdatedAt:   time.Now().Truncate(time.Microsecond),
		AccessCount: 5,
	}

	repo := &mockLinkRepository{
		getFunc: func(ctx context.Context, userID, code string) (*model.Link, error) {
			if code != "abc123" {
				t.Fatalf("expected code abc123, got %q", code)
			}
			return link, nil
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.GET("/api/v1/urls/:shortCode", handler.Get)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/urls/abc123",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response CreateLinkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code != link.Code {
		t.Fatalf("expected code %q, got %q", link.Code, response.Code)
	}

	if response.URL != link.URL {
		t.Fatalf("expected URL %q, got %q", link.URL, response.URL)
	}
}

func TestLinkHandler_Get_NotFound(t *testing.T) {
	repo := &mockLinkRepository{
		getFunc: func(ctx context.Context, userID, code string) (*model.Link, error) {
			return nil, repository.ErrNotFound
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.GET("/api/v1/urls/:shortCode", handler.Get)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/urls/abc123",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestLinkHandler_Get_RepositoryError(t *testing.T) {
	repo := &mockLinkRepository{
		getFunc: func(ctx context.Context, userID, code string) (*model.Link, error) {
			return nil, errors.New("database unavailable")
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.GET("/api/v1/urls/:shortCode", handler.Get)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/urls/abc123",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestLinkHandler_List(t *testing.T) {
	expected := []*model.Link{
		{
			ID:          uuid.New(),
			UserID:      "user_test_123",
			URL:         "https://example.com",
			Code:        "abc123",
			AccessCount: 5,
		},
		{
			ID:          uuid.New(),
			UserID:      "user_test_123",
			URL:         "https://google.com",
			Code:        "xyz789",
			AccessCount: 2,
		},
	}

	repo := &mockLinkRepository{
		listFunc: func(ctx context.Context, userID string) ([]*model.Link, error) {
			return expected, nil
		},
	}

	handler := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/urls", nil)
	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Set("userID", "user_test_123")

	handler.List(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var got []*model.Link
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(got) != len(expected) {
		t.Fatalf("expected %d links, got %d", len(expected), len(got))
	}

	for i := range expected {
		if got[i].ID != expected[i].ID {
			t.Fatalf("expected ID %v, got %v", expected[i].ID, got[i].ID)
		}

		if got[i].UserID != expected[i].UserID {
			t.Fatalf("expected UserID %q, got %q", expected[i].UserID, got[i].UserID)
		}

		if got[i].URL != expected[i].URL {
			t.Fatalf("expected URL %q, got %q", expected[i].URL, got[i].URL)
		}

		if got[i].Code != expected[i].Code {
			t.Fatalf("expected code %q, got %q", expected[i].Code, got[i].Code)
		}

		if got[i].AccessCount != expected[i].AccessCount {
			t.Fatalf("expected access count %d, got %d",
				expected[i].AccessCount, got[i].AccessCount)
		}
	}
}

func TestLinkHandler_List_Empty(t *testing.T) {
	repo := &mockLinkRepository{
		listFunc: func(ctx context.Context, userID string) ([]*model.Link, error) {
			return []*model.Link{}, nil
		},
	}

	handler := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/urls", nil)
	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Set("userID", "user_test_123")

	handler.List(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var got []*model.Link
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected 0 links, got %d", len(got))
	}
}

func TestLinkHandler_List_Unauthorized(t *testing.T) {
	repo := &mockLinkRepository{
		listFunc: func(ctx context.Context, userID string) ([]*model.Link, error) {
			t.Fatal("repository should not be called")
			return nil, nil
		},
	}

	router := gin.New()
	handler := newTestHandler(repo)
	router.GET("/api/v1/urls", handler.List)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/urls",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestLinkHandler_List_ServiceError(t *testing.T) {
	repo := &mockLinkRepository{
		listFunc: func(ctx context.Context, userID string) ([]*model.Link, error) {
			return nil, errors.New("database unavailable")
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.GET("/api/v1/urls", handler.List)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/urls",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestLinkHandler_Stats(t *testing.T) {
	link := &model.Link{
		ID:          uuid.New(),
		URL:         "https://example.com",
		Code:        "abc123",
		CreatedAt:   time.Now().Truncate(time.Microsecond),
		UpdatedAt:   time.Now().Truncate(time.Microsecond),
		AccessCount: 42,
	}

	repo := &mockLinkRepository{
		getFunc: func(ctx context.Context, userID, code string) (*model.Link, error) {
			if code != "abc123" {
				t.Fatalf("expected code abc123, got %q", code)
			}
			return link, nil
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.GET("/api/v1/urls/:shortCode/stats", handler.Stats)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/urls/abc123/stats",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response LinkStatsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code != link.Code {
		t.Fatalf("expected code %q, got %q", link.Code, response.Code)
	}

	if response.URL != link.URL {
		t.Fatalf("expected URL %q, got %q", link.URL, response.URL)
	}

	if response.AccessCount != link.AccessCount {
		t.Fatalf("expected AccessCount %d, got %d", link.AccessCount, response.AccessCount)
	}

	if !response.CreatedAt.Equal(link.CreatedAt) {
		t.Fatalf("expected CreatedAt %v, got %v", link.CreatedAt, response.CreatedAt)
	}

	if !response.UpdatedAt.Equal(link.UpdatedAt) {
		t.Fatalf("expected UpdatedAt %v, got %v", link.UpdatedAt, response.UpdatedAt)
	}
}

func TestLinkHandler_Stats_NotFound(t *testing.T) {
	repo := &mockLinkRepository{
		getFunc: func(ctx context.Context, userID, code string) (*model.Link, error) {
			return nil, repository.ErrNotFound
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.GET("/api/v1/urls/:shortCode/stats", handler.Stats)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/urls/abc123/stats",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestLinkHandler_Stats_RepositoryError(t *testing.T) {
	repo := &mockLinkRepository{
		getFunc: func(ctx context.Context, userID, code string) (*model.Link, error) {
			return nil, errors.New("database unavailable")
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.GET("/api/v1/urls/:shortCode/stats", handler.Stats)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/urls/abc123/stats",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestLinkHandler_Update(t *testing.T) {
	link := &model.Link{
		ID:          uuid.New(),
		URL:         "https://example.org",
		Code:        "abc123",
		CreatedAt:   time.Now().Truncate(time.Microsecond),
		UpdatedAt:   time.Now().Truncate(time.Microsecond),
		AccessCount: 5,
	}

	repo := &mockLinkRepository{
		updateFunc: func(ctx context.Context, userID, code, rawURL string) (*model.Link, error) {
			if code != "abc123" {
				t.Fatalf("expected code abc123, got %q", code)
			}

			if rawURL != "https://example.org" {
				t.Fatalf("expected URL https://example.org, got %q", rawURL)
			}

			return link, nil
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.PUT("/api/v1/urls/:shortCode", handler.Update)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/urls/abc123",
		strings.NewReader(`{"url":"https://example.org"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response CreateLinkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code != link.Code {
		t.Fatalf("expected code %q, got %q", link.Code, response.Code)
	}

	if response.ShortURL != "http://localhost:8080/abc123" {
		t.Fatalf("expected short URL %q, got %q", "http://localhost:8080/abc123", response.ShortURL)
	}

	if response.URL != link.URL {
		t.Fatalf("expected URL %q, got %q", link.URL, response.URL)
	}
}

func TestLinkHandler_Update_InvalidJSON(t *testing.T) {
	repo := &mockLinkRepository{}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.PUT("/api/v1/urls/:shortCode", handler.Update)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/urls/abc123",
		strings.NewReader(`{"url":`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestLinkHandler_Update_InvalidURL(t *testing.T) {
	repo := &mockLinkRepository{}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.PUT("/api/v1/urls/:shortCode", handler.Update)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/urls/abc123",
		strings.NewReader(`{"url":"not-a-url"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestLinkHandler_Update_NotFound(t *testing.T) {
	repo := &mockLinkRepository{
		updateFunc: func(ctx context.Context, userID, code, rawURL string) (*model.Link, error) {
			return nil, repository.ErrNotFound
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.PUT("/api/v1/urls/:shortCode", handler.Update)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/urls/abc123",
		strings.NewReader(`{"url":"https://example.org"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestLinkHandler_Update_RepositoryError(t *testing.T) {
	repo := &mockLinkRepository{
		updateFunc: func(ctx context.Context, userID, code, rawURL string) (*model.Link, error) {
			return nil, errors.New("database unavailable")
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.PUT("/api/v1/urls/:shortCode", handler.Update)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/urls/abc123",
		strings.NewReader(`{"url":"https://example.org"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestLinkHandler_Delete(t *testing.T) {
	var receivedCode string

	repo := &mockLinkRepository{
		deleteFunc: func(ctx context.Context, userID, code string) error {
			receivedCode = code
			return nil
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.DELETE("/api/v1/urls/:shortCode", handler.Delete)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/urls/abc123",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	if rec.Body.Len() != 0 {
		t.Fatalf("expected empty response body, got %q", rec.Body.String())
	}

	if receivedCode != "abc123" {
		t.Fatalf("expected code abc123, got %q", receivedCode)
	}
}

func TestLinkHandler_Delete_NotFound(t *testing.T) {
	repo := &mockLinkRepository{
		deleteFunc: func(ctx context.Context, userID, code string) error {
			return repository.ErrNotFound
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.DELETE("/api/v1/urls/:shortCode", handler.Delete)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/urls/abc123",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestLinkHandler_Delete_RepositoryError(t *testing.T) {
	repo := &mockLinkRepository{
		deleteFunc: func(ctx context.Context, userID, code string) error {
			return errors.New("database unavailable")
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.DELETE("/api/v1/urls/:shortCode", handler.Delete)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/urls/abc123",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestLinkHandler_Redirect(t *testing.T) {
	link := &model.Link{
		ID:          uuid.New(),
		URL:         "https://example.com",
		Code:        "abc123",
		CreatedAt:   time.Now().Truncate(time.Microsecond),
		UpdatedAt:   time.Now().Truncate(time.Microsecond),
		AccessCount: 1,
	}

	repo := &mockLinkRepository{
		getAndIncrementFunc: func(ctx context.Context, code string) (*model.Link, error) {
			if code != "abc123" {
				t.Fatalf("expected code abc123, got %q", code)
			}
			return link, nil
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.GET("/:shortCode", handler.Redirect)

	req := httptest.NewRequest(
		http.MethodGet,
		"/abc123",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status %d, got %d", http.StatusFound, rec.Code)
	}

	if location := rec.Header().Get("Location"); location != link.URL {
		t.Fatalf("expected Location %q, got %q", link.URL, location)
	}
}

func TestLinkHandler_Redirect_NotFound(t *testing.T) {
	repo := &mockLinkRepository{
		getAndIncrementFunc: func(ctx context.Context, code string) (*model.Link, error) {
			return nil, repository.ErrNotFound
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.GET("/:shortCode", handler.Redirect)

	req := httptest.NewRequest(
		http.MethodGet,
		"/abc123",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestLinkHandler_Redirect_RepositoryError(t *testing.T) {
	repo := &mockLinkRepository{
		getAndIncrementFunc: func(ctx context.Context, code string) (*model.Link, error) {
			return nil, errors.New("database unavailable")
		},
	}

	router := setupRouter()
	handler := newTestHandler(repo)
	router.GET("/:shortCode", handler.Redirect)

	req := httptest.NewRequest(
		http.MethodGet,
		"/abc123",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
