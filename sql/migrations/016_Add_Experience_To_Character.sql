-- +goose Up
ALTER TABLE characters
ADD COLUMN experience_points INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE characters
DROP COLUMN experience_points;
