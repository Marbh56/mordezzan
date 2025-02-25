-- +goose Up
-- Create base spells table
CREATE TABLE spells (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    range TEXT NOT NULL,
    duration TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create table for spell class levels
CREATE TABLE spell_class_levels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    spell_id INTEGER NOT NULL,
    class TEXT NOT NULL,
    level INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (spell_id) REFERENCES spells(id) ON DELETE CASCADE,
    UNIQUE(spell_id, class)
);

-- Create indexes for better query performance
CREATE INDEX idx_spells_name ON spells(name);
CREATE INDEX idx_spell_class_levels_spell_id ON spell_class_levels(spell_id);
CREATE INDEX idx_spell_class_levels_class ON spell_class_levels(class);
CREATE INDEX idx_spell_class_levels_level ON spell_class_levels(level);

-- +goose Down
DROP INDEX IF EXISTS idx_spells_name;
DROP INDEX IF EXISTS idx_spell_class_levels_level;
DROP INDEX IF EXISTS idx_spell_class_levels_class;
DROP INDEX IF EXISTS idx_spell_class_levels_spell_id;
DROP TABLE IF EXISTS spell_class_levels;
DROP TABLE IF EXISTS spells;