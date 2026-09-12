package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/IshaanShivKr/urlvio/internal/model"
	"github.com/IshaanShivKr/urlvio/internal/repository"
	"github.com/google/uuid"
)

const (
	codeLength = 6
	maxRetries = 5
	alphabet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

type LinkService struct {
	repo repository.LinkRepository
}

func NewLinkService(repo repository.LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}

func (s *LinkService) Create(ctx context.Context, rawURL string) (*model.Link, error) {
	rawURL = strings.TrimSpace(rawURL)

	if rawURL == "" {
		return nil, ErrURLRequired
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Host == "" {
		return nil, ErrInvalidURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, ErrInvalidURL
	}

	now := time.Now()

	link := &model.Link{
		ID:          uuid.New(),
		URL:         rawURL,
		CreatedAt:   now,
		UpdatedAt:   now,
		AccessCount: 0,
	}

	for range maxRetries {
		code, err := generateCode()
		if err != nil {
			return nil, fmt.Errorf("generate code: %w", err)
		}

		link.Code = code

		if err := s.repo.Create(ctx, link); err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
				continue
			}

			return nil, fmt.Errorf("create link: %w", err)
		}

		return link, nil
	}

	return nil, fmt.Errorf("%w after %d attempts", ErrCodeGeneration, maxRetries)
}

func (s *LinkService) Get(ctx context.Context, code string) (*model.Link, error) {
	code = strings.TrimSpace(code)

	if len(code) != codeLength {
		return nil, ErrNotFound
	}

	link, err := s.repo.Get(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get link: %w", err)
	}

	return link, nil
}

func (s *LinkService) Delete(ctx context.Context, code string) error {
	code = strings.TrimSpace(code)

	if len(code) != codeLength {
		return ErrNotFound
	}

	if err := s.repo.Delete(ctx, code); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return fmt.Errorf("delete link: %w", err)
	}

	return nil
}

func (s *LinkService) Update(ctx context.Context, code, rawURL string) (*model.Link, error) {
	code = strings.TrimSpace(code)
	rawURL = strings.TrimSpace(rawURL)

	if len(code) != codeLength {
		return nil, ErrNotFound
	}

	if rawURL == "" {
		return nil, ErrURLRequired
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Host == "" {
		return nil, ErrInvalidURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, ErrInvalidURL
	}

	link, err := s.repo.Update(ctx, code, rawURL)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("update link: %w", err)
	}

	return link, nil
}

func generateCode() (string, error) {
	code := make([]byte, codeLength)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", fmt.Errorf("generate random index: %w", err)
		}

		code[i] = alphabet[n.Int64()]
	}

	return string(code), nil
}
