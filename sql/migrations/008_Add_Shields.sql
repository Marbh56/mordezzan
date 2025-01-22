-- +goose Up
CREATE TABLE shields (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    cost_gp INTEGER NOT NULL,
    weight INTEGER NOT NULL,
    defense_bonus INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert shield data
INSERT INTO
    shields (name, cost_gp, weight, defense_bonus)
VALUES
    ('Small Shield', 5, 5, 1),
    ('Large Shield', 10, 10, 2);

-- +goose Down
DROP TABLE IF EXISTS shields;
