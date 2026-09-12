package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IshaanShivKr/urlvio/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, link *model.Link) error {
	const query = `
		INSERT INTO links (
			id,
			url,
			code,
			created_at,
			updated_at,
			access_count
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	if _, err := r.db.Exec(
		ctx,
		query,
		link.ID,
		link.URL,
		link.Code,
		link.CreatedAt,
		link.UpdatedAt,
		link.AccessCount,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "links_code_key" {
			return ErrAlreadyExists
		}

		return fmt.Errorf("create link: %w", err)
	}

	return nil
}
