package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IshaanShivKr/urlvio/internal/model"
	"github.com/jackc/pgx/v5"
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
			user_id,
			url,
			code,
			created_at,
			updated_at,
			access_count
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := r.db.Exec(
		ctx,
		query,
		link.ID,
		link.UserID,
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

func (r *PostgresRepository) Get(ctx context.Context, userID, code string) (*model.Link, error) {
	const query = `
		SELECT
			id,
			user_id,
			url,
			code,
			created_at,
			updated_at,
			access_count
		FROM links
		WHERE user_id = $1 AND code = $2
	`

	link := &model.Link{}

	if err := r.db.QueryRow(ctx, query, userID, code).Scan(
		&link.ID,
		&link.UserID,
		&link.URL,
		&link.Code,
		&link.CreatedAt,
		&link.UpdatedAt,
		&link.AccessCount,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("query link: %w", err)
	}

	return link, nil
}

func (r *PostgresRepository) List(ctx context.Context, userID string) ([]*model.Link, error) {
	const query = `
		SELECT
			id,
			user_id,
			url,
			code,
			created_at,
			updated_at,
			access_count
		FROM links
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list links: %w", err)
	}
	defer rows.Close()

	var links []*model.Link

	for rows.Next() {
		link := &model.Link{}

		if err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.URL,
			&link.Code,
			&link.CreatedAt,
			&link.UpdatedAt,
			&link.AccessCount,
		); err != nil {
			return nil, fmt.Errorf("scan link: %w", err)
		}

		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate links: %w", err)
	}

	return links, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, userID, code string) error {
	const query = `
		DELETE FROM links
		WHERE user_id = $1 AND code = $2
	`

	result, err := r.db.Exec(ctx, query, userID, code)
	if err != nil {
		return fmt.Errorf("delete link: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, userID, code, rawURL string) (*model.Link, error) {
	const query = `
		UPDATE links
		SET
			url = $1,
			updated_at = $2
		WHERE user_id = $3 AND code = $4
		RETURNING
			id,
			user_id,
			url,
			code,
			created_at,
			updated_at,
			access_count
	`

	link := &model.Link{}

	if err := r.db.QueryRow(
		ctx,
		query,
		rawURL,
		time.Now(),
		userID,
		code,
	).Scan(
		&link.ID,
		&link.UserID,
		&link.URL,
		&link.Code,
		&link.CreatedAt,
		&link.UpdatedAt,
		&link.AccessCount,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("update link: %w", err)
	}

	return link, nil
}

func (r *PostgresRepository) GetAndIncrement(ctx context.Context, code string) (*model.Link, error) {
	const query = `
		UPDATE links
		SET access_count = access_count + 1
		WHERE code = $1
		RETURNING
			id,
			user_id,
			url,
			code,
			created_at,
			updated_at,
			access_count
	`

	link := &model.Link{}

	if err := r.db.QueryRow(
		ctx,
		query,
		code,
	).Scan(
		&link.ID,
		&link.UserID,
		&link.URL,
		&link.Code,
		&link.CreatedAt,
		&link.UpdatedAt,
		&link.AccessCount,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get and increment link: %w", err)
	}

	return link, nil
}
