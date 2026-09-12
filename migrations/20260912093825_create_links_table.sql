-- +goose Up
CREATE TABLE links (
    id UUID PRIMARY KEY,

    url TEXT NOT NULL
        CHECK (length(trim(url)) > 0),

    code VARCHAR(6) NOT NULL
        UNIQUE
        CHECK (code ~ '^[A-Za-z0-9]{6}$'),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    access_count BIGINT NOT NULL DEFAULT 0
        CHECK (access_count >= 0)
);

-- +goose Down
DROP TABLE links;
