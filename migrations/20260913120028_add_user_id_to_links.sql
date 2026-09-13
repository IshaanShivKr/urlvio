-- +goose Up
ALTER TABLE links
ADD COLUMN user_id TEXT NOT NULL;

-- +goose Down
ALTER TABLE links
DROP COLUMN user_id;
