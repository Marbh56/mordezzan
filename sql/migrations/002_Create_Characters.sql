-- +goose Up
CREATE TABLE characters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    class TEXT NOT NULL DEFAULT 'Fighter',
    level INTEGER NOT NULL DEFAULT 1,
    max_hp INTEGER NOT NULL,
    current_hp INTEGER NOT NULL,
    strength INTEGER NOT NULL,
    dexterity INTEGER NOT NULL,
    constitution INTEGER NOT NULL,
    intelligence INTEGER NOT NULL,
    wisdom INTEGER NOT NULL,
    charisma INTEGER NOT NULL,
    experience_points INTEGER NOT NULL DEFAULT 0,
    platinum_pieces INTEGER NOT NULL DEFAULT 0,
    gold_pieces INTEGER NOT NULL DEFAULT 0,
    electrum_pieces INTEGER NOT NULL DEFAULT 0,
    silver_pieces INTEGER NOT NULL DEFAULT 0,
    copper_pieces INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE IF EXISTS characters;