-- +goose Up
CREATE TABLE magical_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    weight INTEGER NOT NULL DEFAULT 0,
    cost_gp INTEGER NOT NULL DEFAULT 0,
    max_charges INTEGER NOT NULL DEFAULT 1,
    category TEXT NOT NULL CHECK (category IN ('wand', 'potion', 'rod', 'scroll', 'other')),
    effect_description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS magical_items;