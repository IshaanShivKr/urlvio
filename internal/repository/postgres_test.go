package repository

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/IshaanShivKr/urlvio/internal/database"
	"github.com/IshaanShivKr/urlvio/internal/model"
	"github.com/google/uuid"
)

func newTestRepository(t *testing.T) *PostgresRepository {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repo := NewPostgresRepository(db)

	t.Cleanup(func() {
		if _, err := db.Exec(ctx, "TRUNCATE TABLE links"); err != nil {
			t.Fatalf("failed to clean test database: %v", err)
		}
	})

	return repo
}

func newTestLink(code, rawURL string) *model.Link {
	now := time.Now()

	return &model.Link{
		ID:          uuid.New(),
		URL:         rawURL,
		Code:        code,
		CreatedAt:   now,
		UpdatedAt:   now,
		AccessCount: 0,
	}
}

func TestPostgresRepository_Create(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	link := newTestLink("abc123", "https://example.com")

	err := repo.Create(ctx, link)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := repo.Get(ctx, "abc123")
	if err != nil {
		t.Fatalf("expected created link to be found, got %v", err)
	}

	if got.ID != link.ID {
		t.Fatalf("expected ID %v, got %v", link.ID, got.ID)
	}

	if got.URL != link.URL {
		t.Fatalf("expected URL %q, got %q", link.URL, got.URL)
	}

	if got.Code != link.Code {
		t.Fatalf("expected code %q, got %q", link.Code, got.Code)
	}

	if got.AccessCount != link.AccessCount {
		t.Fatalf("expected access count %d, got %d", link.AccessCount, got.AccessCount)
	}
}

func TestPostgresRepository_Create_DuplicateCode(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	first := newTestLink("abc123", "https://example.com")
	second := newTestLink("abc123", "https://google.com")

	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("failed to create first link: %v", err)
	}

	err := repo.Create(ctx, second)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestPostgresRepository_Get(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	link := newTestLink("abc123", "https://example.com")

	if err := repo.Create(ctx, link); err != nil {
		t.Fatalf("failed to create link: %v", err)
	}

	got, err := repo.Get(ctx, "abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != link.ID {
		t.Fatalf("expected ID %v, got %v", link.ID, got.ID)
	}

	if got.URL != link.URL {
		t.Fatalf("expected URL %q, got %q", link.URL, got.URL)
	}

	if got.Code != link.Code {
		t.Fatalf("expected code %q, got %q", link.Code, got.Code)
	}

	if got.AccessCount != link.AccessCount {
		t.Fatalf("expected AccessCount %d, got %d", link.AccessCount, got.AccessCount)
	}
}

func TestPostgresRepository_Get_NotFound(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	_, err := repo.Get(ctx, "abc123")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresRepository_Update(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	link := newTestLink("abc123", "https://example.com")

	if err := repo.Create(ctx, link); err != nil {
		t.Fatalf("failed to create link: %v", err)
	}

	updatedURL := "https://example.org"

	got, err := repo.Update(ctx, "abc123", updatedURL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != link.ID {
		t.Fatalf("expected ID %v, got %v", link.ID, got.ID)
	}

	if got.Code != link.Code {
		t.Fatalf("expected code %q, got %q", link.Code, got.Code)
	}

	if got.URL != updatedURL {
		t.Fatalf("expected URL %q, got %q", updatedURL, got.URL)
	}

	if got.AccessCount != link.AccessCount {
		t.Fatalf("expected AccessCount %d, got %d", link.AccessCount, got.AccessCount)
	}

	if !got.UpdatedAt.After(link.UpdatedAt) {
		t.Fatalf("expected UpdatedAt to be updated, got %v", got.UpdatedAt)
	}
}

func TestPostgresRepository_Update_NotFound(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	_, err := repo.Update(ctx, "abc123", "https://example.org")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresRepository_Delete(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	link := newTestLink("abc123", "https://example.com")

	if err := repo.Create(ctx, link); err != nil {
		t.Fatalf("failed to create link: %v", err)
	}

	if err := repo.Delete(ctx, "abc123"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := repo.Get(ctx, "abc123")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected link to be deleted, got error %v", err)
	}
}

func TestPostgresRepository_Delete_NotFound(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	err := repo.Delete(ctx, "abc123")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresRepository_GetAndIncrement(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	link := newTestLink("abc123", "https://example.com")
	link.AccessCount = 5

	if err := repo.Create(ctx, link); err != nil {
		t.Fatalf("failed to create link: %v", err)
	}

	got, err := repo.GetAndIncrement(ctx, "abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != link.ID {
		t.Fatalf("expected ID %v, got %v", link.ID, got.ID)
	}

	if got.URL != link.URL {
		t.Fatalf("expected URL %q, got %q", link.URL, got.URL)
	}

	if got.Code != link.Code {
		t.Fatalf("expected code %q, got %q", link.Code, got.Code)
	}

	if got.AccessCount != 6 {
		t.Fatalf("expected AccessCount 6, got %d", got.AccessCount)
	}

	stored, err := repo.Get(ctx, "abc123")
	if err != nil {
		t.Fatalf("failed to get updated link: %v", err)
	}

	if stored.AccessCount != 6 {
		t.Fatalf("expected persisted AccessCount 6, got %d", stored.AccessCount)
	}
}

func TestPostgresRepository_GetAndIncrement_NotFound(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	_, err := repo.GetAndIncrement(ctx, "abc123")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresRepository_GetAndIncrement_Concurrent(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	link := newTestLink("abc123", "https://example.com")

	if err := repo.Create(ctx, link); err != nil {
		t.Fatalf("failed to create link: %v", err)
	}

	const goroutines = 100

	errs := make(chan error, goroutines)

	for range goroutines {
		go func() {
			_, err := repo.GetAndIncrement(ctx, "abc123")
			errs <- err
		}()
	}

	for range goroutines {
		if err := <-errs; err != nil {
			t.Fatalf("concurrent increment failed: %v", err)
		}
	}

	got, err := repo.Get(ctx, "abc123")
	if err != nil {
		t.Fatalf("failed to get link after concurrent increments: %v", err)
	}

	if got.AccessCount != goroutines {
		t.Fatalf("expected AccessCount %d, got %d", goroutines, got.AccessCount)
	}
}
