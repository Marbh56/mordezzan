-- +goose Up
CREATE TABLE ammunition (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    cost_gp DECIMAL(10, 2) NOT NULL,
    weight INTEGER,
    quantity INTEGER,
    enhancement_bonus INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert standard ammunition data
INSERT INTO
    ammunition (name, cost_gp, weight, quantity, enhancement_bonus)
VALUES
    ('Arrow, Silver-tipped', 10.0, 0, 1, 0),
    ('Arrows', 5.0, 1, 12, 0),
    ('Bolt, Heavy, Silver-tipped', 10.0, 0, 1, 0),
    ('Bolt, Light, Silver-tipped', 10.0, 0, 1, 0),
    ('Bolts, Heavy', 5.0, 2, 10, 0),
    ('Bolts, Light', 5.0, 2, 20, 0),
    ('Bullet, Sling, Silver', 2.0, 0, 1, 0),
    ('Bullets, Sling, Lead', 0.5, 2, 20, 0);

-- +goose Down
DROP TABLE IF EXISTS ammunition;