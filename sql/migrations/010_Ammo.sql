-- +goose Up
CREATE TABLE ammunition (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    cost_gp DECIMAL(10, 2) NOT NULL,
    weight INTEGER,
    quantity INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert ammunition data
INSERT INTO
    ammunition (name, cost_gp, weight, quantity)
VALUES
    ('Arrow, Silver-tipped', 10.0, 0, 1),
    ('Arrows', 5.0, 1, 12),
    ('Bolt, Heavy, Silver-tipped', 10.0, 0, 1),
    ('Bolt, Light, Silver-tipped', 10.0, 0, 1),
    ('Bolts, Heavy', 5.0, 2, 10),
    ('Bolts, Light', 5.0, 2, 20),
    ('Bullet, Sling, Silver', 2.0, 0, 1),
    ('Bullets, Sling, Lead', 0.5, 2, 20);

-- +goose Down
DROP TABLE IF EXISTS ammunition;
