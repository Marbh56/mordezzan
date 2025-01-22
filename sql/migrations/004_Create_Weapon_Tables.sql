-- +goose Up
-- Create tables
CREATE TABLE weapon_properties (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE weapons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    reach INTEGER NOT NULL,
    cost_gp INTEGER NOT NULL,
    weight INTEGER NOT NULL,
    range_short INTEGER,
    range_medium INTEGER,
    range_long INTEGER,
    attacks_per_round TEXT,
    damage TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE weapon_property_links (
    weapon_id INTEGER NOT NULL,
    property_id INTEGER NOT NULL,
    FOREIGN KEY (weapon_id) REFERENCES weapons (id),
    FOREIGN KEY (property_id) REFERENCES weapon_properties (id),
    PRIMARY KEY (weapon_id, property_id)
);

-- Insert weapon properties
INSERT INTO
    weapon_properties (symbol, description)
VALUES
    (
        '↵',
        'Ignores AC bonus provided by opponent''s shield'
    ),
    (
        'Ω',
        '+1 "to hit" opponents wearing plate mail, field plate, or full plate'
    ),
    (
        '+',
        'A "true" two-handed melee weapon; must be wielded with two hands'
    ),
    ('↔', '+1 AC bonus versus melee attacks'),
    (
        '^',
        'Double damage dice when set to receive a charge'
    ),
    (
        '#',
        '4-in-6 chance to dismount a rider on a natural 19 or 20 attack roll'
    ),
    (
        '∇',
        'Double damage dice when used from a charging mount'
    ),
    (
        'o',
        'Base damage improves to 1d10 when mounted on a heavy warhorse'
    );

-- +goose Down
DROP TABLE IF EXISTS weapon_property_links;

DROP TABLE IF EXISTS weapons;

DROP TABLE IF EXISTS weapon_properties;
