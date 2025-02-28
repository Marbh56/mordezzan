-- +goose Up
-- Create the main spells table
CREATE TABLE spells (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    range TEXT NOT NULL,
    duration TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create a table for spell levels by class
CREATE TABLE spell_levels (
    spell_id INTEGER NOT NULL,
    class TEXT NOT NULL,
    level INTEGER NOT NULL,
    PRIMARY KEY (spell_id, class),
    FOREIGN KEY (spell_id) REFERENCES spells (id) ON DELETE CASCADE
);

-- Create index for faster lookup by class and level
CREATE INDEX idx_spell_levels_class_level ON spell_levels (class, level);

-- +goose Down
DROP INDEX IF EXISTS idx_spell_levels_class_level;
DROP TABLE IF EXISTS spell_levels;
DROP TABLE IF EXISTS spells;