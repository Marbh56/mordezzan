-- +goose Up
-- Equipment slots reference table
CREATE TABLE equipment_slots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

-- Insert basic equipment slots
INSERT INTO
    equipment_slots (name, description)
VALUES
    ('head', 'Headwear, helmets'),
    ('neck', 'Necklaces, amulets'),
    ('body', 'Body armor, robes'),
    ('back', 'Cloaks, backpacks'),
    ('right_hand', 'Weapons, shields, tools'),
    ('left_hand', 'Weapons, shields, tools'),
    ('waist', 'Belts, girdles'),
    ('legs', 'Leg armor, greaves'),
    ('feet', 'Boots, shoes'),
    (
        'right_ring_1',
        'Ring slot 1 - right index finger'
    ),
    (
        'right_ring_2',
        'Ring slot 2 - right middle finger'
    ),
    ('right_ring_3', 'Ring slot 3 - right ring finger'),
    (
        'right_ring_4',
        'Ring slot 4 - right pinky finger'
    ),
    ('left_ring_1', 'Ring slot 1 - left index finger'),
    ('left_ring_2', 'Ring slot 2 - left middle finger'),
    ('left_ring_3', 'Ring slot 3 - left ring finger'),
    ('left_ring_4', 'Ring slot 4 - left pinky finger');

-- +goose Down
DROP TABLE IF EXISTS equipment_slots;
