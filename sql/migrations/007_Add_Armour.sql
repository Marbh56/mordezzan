-- +goose Up
CREATE TABLE armor (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    armor_class INTEGER NOT NULL,
    cost_gp INTEGER NOT NULL,
    damage_reduction INTEGER NOT NULL,
    weight INTEGER NOT NULL,
    armor_type TEXT NOT NULL CHECK (armor_type IN ('Light', 'Medium', 'Heavy')),
    movement_rate INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert armor data
INSERT INTO
    armor (
        name,
        armor_class,
        cost_gp,
        damage_reduction,
        weight,
        armor_type,
        movement_rate
    )
VALUES
    ('None', 9, 0, 0, 0, 'Light', 40),
    ('Padded', 8, 10, 0, 10, 'Light', 40),
    ('Leather', 7, 15, 0, 15, 'Light', 40),
    ('Studded', 6, 25, 0, 20, 'Light', 40),
    ('Scale Mail', 6, 50, 1, 25, 'Medium', 30),
    ('Chain Mail', 5, 75, 1, 30, 'Medium', 30),
    ('Laminated', 5, 75, 1, 30, 'Medium', 30),
    ('Banded Mail', 4, 150, 1, 35, 'Medium', 30),
    ('Splint Mail', 4, 150, 1, 35, 'Medium', 30),
    ('Plate Mail', 3, 350, 2, 40, 'Heavy', 20),
    ('Field Plate', 2, 1000, 2, 50, 'Heavy', 20),
    ('Full Plate', 1, 2000, 2, 60, 'Heavy', 20);

-- +goose Down
DROP TABLE IF EXISTS armor;
