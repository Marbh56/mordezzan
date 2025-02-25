-- +goose Up
CREATE TABLE shields (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    cost_gp INTEGER NOT NULL,
    weight INTEGER NOT NULL,
    defense_bonus INTEGER NOT NULL,
    enhancement_bonus INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert standard shield data
INSERT INTO
    shields (name, cost_gp, weight, defense_bonus, enhancement_bonus)
VALUES
    ('Small Shield', 5, 5, 1, 0),
    ('Large Shield', 10, 10, 2, 0);

-- Insert magical shield data
INSERT INTO
    shields (name, cost_gp, weight, defense_bonus, enhancement_bonus)
VALUES
    -- Small Shields
    ('Small Shield +1', 2750, 5, 2, 1), 
    ('Small Shield +2', 4750, 5, 2, 2), 
    ('Small Shield +3', 7500, 5, 2, 3), 
    
    -- Large Shields
    ('Large Shield +1', 3500, 10, 3, 1), 
    ('Large Shield +2', 7000, 10, 3, 2), 
    ('Large Shield +3', 10000, 10, 3, 3);

-- +goose Down
DROP TABLE IF EXISTS shields;