-- +goose Up
ALTER TABLE items ADD COLUMN image_path TEXT;

-- +goose Down
ALTER TABLE items DROP COLUMN image_path;
