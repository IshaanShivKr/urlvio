package repository

import (
	"context"

	"github.com/IshaanShivKr/urlvio/internal/model"
)

type LinkRepository interface {
	Create(ctx context.Context, link *model.Link) error
}
