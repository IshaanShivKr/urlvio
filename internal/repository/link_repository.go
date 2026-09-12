package repository

import (
	"context"

	"github.com/IshaanShivKr/urlvio/internal/model"
)

type LinkRepository interface {
	Create(ctx context.Context, link *model.Link) error
	Get(ctx context.Context, code string) (*model.Link, error)
	Delete(ctx context.Context, code string) error
}
