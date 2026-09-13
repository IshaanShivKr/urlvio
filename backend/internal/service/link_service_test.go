package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/IshaanShivKr/urlvio/internal/model"
	"github.com/IshaanShivKr/urlvio/internal/repository"
	"github.com/google/uuid"
)

const testUserID = "user_test_123"

type mockLinkRepository struct {
	createFunc          func(context.Context, *model.Link) error
	getFunc             func(context.Context, string, string) (*model.Link, error)
	deleteFunc          func(context.Context, string, string) error
	updateFunc          func(context.Context, string, string, string) (*model.Link, error)
	getAndIncrementFunc func(context.Context, string) (*model.Link, error)
}

func (m *mockLinkRepository) Create(ctx context.Context, link *model.Link) error {
	return m.createFunc(ctx, link)
}

func (m *mockLinkRepository) Get(ctx context.Context, userID, code string) (*model.Link, error) {
	return m.getFunc(ctx, userID, code)
}

func (m *mockLinkRepository) Delete(ctx context.Context, userID, code string) error {
	return m.deleteFunc(ctx, userID, code)
}

func (m *mockLinkRepository) Update(ctx context.Context, userID, code, rawURL string) (*model.Link, error) {
	return m.updateFunc(ctx, userID, code, rawURL)
}

func (m *mockLinkRepository) GetAndIncrement(ctx context.Context, code string) (*model.Link, error) {
	return m.getAndIncrementFunc(ctx, code)
}

func newTestLink(code, rawURL string) *model.Link {
	now := time.Now()

	return &model.Link{
		ID:          uuid.New(),
		UserID:      testUserID,
		URL:         rawURL,
		Code:        code,
		CreatedAt:   now,
		UpdatedAt:   now,
		AccessCount: 0,
	}
}

func TestLinkService_Get(t *testing.T) {
	expected := newTestLink("abc123", "https://example.com")

	repo := &mockLinkRepository{
		getFunc: func(ctx context.Context, testUserID, code string) (*model.Link, error) {
			if code != "abc123" {
				t.Fatalf("expected code abc123, got %q", code)
			}

			return expected, nil
		},
	}

	service := NewLinkService(repo)

	link, err := service.Get(context.Background(), testUserID, " abc123 ")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if link != expected {
		t.Fatalf("expected link %v, got %v", expected, link)
	}
}

func TestLinkService_Get_InvalidCode(t *testing.T) {
	repo := &mockLinkRepository{
		getFunc: func(ctx context.Context, testUserID, code string) (*model.Link, error) {
			t.Fatal("repository should not be called")
			return nil, nil
		},
	}

	service := NewLinkService(repo)

	_, err := service.Get(context.Background(), testUserID, "abc")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLinkService_Get_NotFound(t *testing.T) {
	repo := &mockLinkRepository{
		getFunc: func(ctx context.Context, testUserID, code string) (*model.Link, error) {
			return nil, repository.ErrNotFound
		},
	}

	service := NewLinkService(repo)

	_, err := service.Get(context.Background(), testUserID, "abc123")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLinkService_Get_RepositoryError(t *testing.T) {
	repoErr := errors.New("database unavailable")

	repo := &mockLinkRepository{
		getFunc: func(ctx context.Context, testUserID, code string) (*model.Link, error) {
			return nil, repoErr
		},
	}

	service := NewLinkService(repo)

	_, err := service.Get(context.Background(), testUserID, "abc123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repository error, got %v", err)
	}
}

func TestLinkService_Create_EmptyURL(t *testing.T) {
	repo := &mockLinkRepository{
		createFunc: func(ctx context.Context, link *model.Link) error {
			t.Fatal("repository should not be called")
			return nil
		},
	}

	service := NewLinkService(repo)

	_, err := service.Create(context.Background(), testUserID, "   ")

	if !errors.Is(err, ErrURLRequired) {
		t.Fatalf("expected ErrURLRequired, got %v", err)
	}
}

func TestLinkService_Create_InvalidURL(t *testing.T) {
	repo := &mockLinkRepository{
		createFunc: func(ctx context.Context, link *model.Link) error {
			t.Fatal("repository should not be called")
			return nil
		},
	}

	service := NewLinkService(repo)

	_, err := service.Create(context.Background(), testUserID, "not-a-url")

	if !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestLinkService_Create_Success(t *testing.T) {
	repo := &mockLinkRepository{
		createFunc: func(ctx context.Context, link *model.Link) error {
			return nil
		},
	}

	service := NewLinkService(repo)

	link, err := service.Create(context.Background(), testUserID, " https://example.com ")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if link.URL != "https://example.com" {
		t.Fatalf("expected trimmed URL, got %q", link.URL)
	}

	if len(link.Code) != codeLength {
		t.Fatalf("expected code length %d, got %d", codeLength, len(link.Code))
	}

	if link.ID == uuid.Nil {
		t.Fatal("expected non-zero UUID")
	}

	if link.AccessCount != 0 {
		t.Fatalf("expected access count 0, got %d", link.AccessCount)
	}
}

func TestLinkService_Create_RetriesOnCodeCollision(t *testing.T) {
	createCalls := 0

	repo := &mockLinkRepository{
		createFunc: func(ctx context.Context, link *model.Link) error {
			createCalls++

			if createCalls == 1 {
				return repository.ErrAlreadyExists
			}

			return nil
		},
	}

	service := NewLinkService(repo)

	link, err := service.Create(context.Background(), testUserID, "https://example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if link == nil {
		t.Fatal("expected link, got nil")
	}

	if createCalls != 2 {
		t.Fatalf("expected 2 create attempts, got %d", createCalls)
	}
}

func TestLinkService_Create_RepositoryError(t *testing.T) {
	repoErr := errors.New("database unavailable")

	repo := &mockLinkRepository{
		createFunc: func(ctx context.Context, link *model.Link) error {
			return repoErr
		},
	}

	service := NewLinkService(repo)

	_, err := service.Create(context.Background(), testUserID, "https://example.com")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repository error, got %v", err)
	}
}

func TestLinkService_Create_RetryExhaustion(t *testing.T) {
	createCalls := 0

	repo := &mockLinkRepository{
		createFunc: func(ctx context.Context, link *model.Link) error {
			createCalls++
			return repository.ErrAlreadyExists
		},
	}

	service := NewLinkService(repo)

	_, err := service.Create(context.Background(), testUserID, "https://example.com")

	if !errors.Is(err, ErrCodeGeneration) {
		t.Fatalf("expected ErrCodeGeneration, got %v", err)
	}

	if createCalls != maxRetries {
		t.Fatalf("expected %d create attempts, got %d", maxRetries, createCalls)
	}
}

func TestLinkService_Update_Success(t *testing.T) {
	expected := newTestLink("abc123", "https://new-example.com")

	repo := &mockLinkRepository{
		updateFunc: func(ctx context.Context, testUserID, code, rawURL string) (*model.Link, error) {
			if code != "abc123" {
				t.Fatalf("expected code abc123, got %q", code)
			}

			if rawURL != "https://new-example.com" {
				t.Fatalf("expected URL https://new-example.com, got %q", rawURL)
			}

			return expected, nil
		},
	}

	service := NewLinkService(repo)

	link, err := service.Update(
		context.Background(),
		testUserID,
		" abc123 ",
		" https://new-example.com ",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if link != expected {
		t.Fatalf("expected link %v, got %v", expected, link)
	}
}

func TestLinkService_Update_InvalidCode(t *testing.T) {
	repo := &mockLinkRepository{
		updateFunc: func(ctx context.Context, testUserID, code, rawURL string) (*model.Link, error) {
			t.Fatal("repository should not be called")
			return nil, nil
		},
	}

	service := NewLinkService(repo)

	_, err := service.Update(
		context.Background(),
		testUserID,
		"abc",
		"https://example.com",
	)

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLinkService_Update_EmptyURL(t *testing.T) {
	repo := &mockLinkRepository{
		updateFunc: func(ctx context.Context, testUserID, code, rawURL string) (*model.Link, error) {
			t.Fatal("repository should not be called")
			return nil, nil
		},
	}

	service := NewLinkService(repo)

	_, err := service.Update(
		context.Background(),
		testUserID,
		"abc123",
		"   ",
	)

	if !errors.Is(err, ErrURLRequired) {
		t.Fatalf("expected ErrURLRequired, got %v", err)
	}
}

func TestLinkService_Update_InvalidURL(t *testing.T) {
	repo := &mockLinkRepository{
		updateFunc: func(ctx context.Context, testUserID, code, rawURL string) (*model.Link, error) {
			t.Fatal("repository should not be called")
			return nil, nil
		},
	}

	service := NewLinkService(repo)

	_, err := service.Update(
		context.Background(),
		testUserID,
		"abc123",
		"ftp://example.com",
	)

	if !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestLinkService_Update_NotFound(t *testing.T) {
	repo := &mockLinkRepository{
		updateFunc: func(ctx context.Context, testUserID, code, rawURL string) (*model.Link, error) {
			return nil, repository.ErrNotFound
		},
	}

	service := NewLinkService(repo)

	_, err := service.Update(
		context.Background(),
		testUserID,
		"abc123",
		"https://example.com",
	)

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLinkService_Delete_Success(t *testing.T) {
	repo := &mockLinkRepository{
		deleteFunc: func(ctx context.Context, testUserID, code string) error {
			if code != "abc123" {
				t.Fatalf("expected code abc123, got %q", code)
			}

			return nil
		},
	}

	service := NewLinkService(repo)

	err := service.Delete(context.Background(), testUserID, " abc123 ")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestLinkService_Delete_InvalidCode(t *testing.T) {
	repo := &mockLinkRepository{
		deleteFunc: func(ctx context.Context, testUserID, code string) error {
			t.Fatal("repository should not be called")
			return nil
		},
	}

	service := NewLinkService(repo)

	err := service.Delete(context.Background(), testUserID, "abc")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLinkService_Delete_NotFound(t *testing.T) {
	repo := &mockLinkRepository{
		deleteFunc: func(ctx context.Context, testUserID, code string) error {
			return repository.ErrNotFound
		},
	}

	service := NewLinkService(repo)

	err := service.Delete(context.Background(), testUserID, "abc123")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLinkService_GetAndIncrement_Success(t *testing.T) {
	expected := newTestLink("abc123", "https://example.com")
	expected.AccessCount = 1

	repo := &mockLinkRepository{
		getAndIncrementFunc: func(ctx context.Context, code string) (*model.Link, error) {
			if code != "abc123" {
				t.Fatalf("expected code abc123, got %q", code)
			}

			return expected, nil
		},
	}

	service := NewLinkService(repo)

	link, err := service.GetAndIncrement(context.Background(), " abc123 ")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if link != expected {
		t.Fatalf("expected link %v, got %v", expected, link)
	}
}

func TestLinkService_GetAndIncrement_InvalidCode(t *testing.T) {
	repo := &mockLinkRepository{
		getAndIncrementFunc: func(ctx context.Context, code string) (*model.Link, error) {
			t.Fatal("repository should not be called")
			return nil, nil
		},
	}

	service := NewLinkService(repo)

	_, err := service.GetAndIncrement(context.Background(), "abc")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLinkService_GetAndIncrement_NotFound(t *testing.T) {
	repo := &mockLinkRepository{
		getAndIncrementFunc: func(ctx context.Context, code string) (*model.Link, error) {
			return nil, repository.ErrNotFound
		},
	}

	service := NewLinkService(repo)

	_, err := service.GetAndIncrement(context.Background(), "abc123")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLinkService_GetAndIncrement_RepositoryError(t *testing.T) {
	repoErr := errors.New("database unavailable")

	repo := &mockLinkRepository{
		getAndIncrementFunc: func(ctx context.Context, code string) (*model.Link, error) {
			return nil, repoErr
		},
	}

	service := NewLinkService(repo)

	_, err := service.GetAndIncrement(context.Background(), "abc123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repository error, got %v", err)
	}
}
