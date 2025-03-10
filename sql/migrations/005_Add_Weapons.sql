-- +goose Up
INSERT INTO
    weapons (
        name,
        reach,
        cost_gp,
        weight,
        range_short,
        range_medium,
        range_long,
        attacks_per_round,
        damage,
        enhancement_bonus
    )
VALUES
    ('Axe, Hand', 1, 5, 2, 15, 30, 45, '1/1', '1d6', 0),
    ('Axe, Battle', 2, 10, 5, NULL, NULL, NULL, NULL, '1d8 (1d10)', 0),
    ('Axe, Great', 4, 20, 10, NULL, NULL, NULL, NULL, '2d6', 0),
    ('Cæstuses', 0, 1, 1, NULL, NULL, NULL, NULL, '+1', 0),
    ('Chain Whip', 4, 10, 3, NULL, NULL, NULL, NULL, '1d6', 0),
    ('Club, Light', 1, 1, 2, 10, 20, 30, '1/1', '1d4', 0),
    ('Club, War', 2, 3, 4, NULL, NULL, NULL, NULL, '1d6 (1d8)', 0),
    ('Dagger', 1, 4, 1, 10, 20, 30, '3/2', '1d4', 0),
    ('Dagger, Silver', 1, 25, 1, 10, 20, 30, '3/2', '1d4', 0),
    ('Falcata', 1, 10, 3, NULL, NULL, NULL, NULL, '1d6', 0),
    ('Flail, Horseman''s', 1, 5, 3, NULL, NULL, NULL, NULL, '1d6', 0),
    ('Flail, Footman''s', 3, 10, 10, NULL, NULL, NULL, NULL, '1d10', 0),
    ('Halberd', 4, 15, 8, NULL, NULL, NULL, NULL, '1d10', 0),
    ('Hammer, Horseman''s', 1, 5, 3, 10, 20, 30, '1/1', '1d6', 0),
    ('Hammer, War', 2, 10, 5, NULL, NULL, NULL, NULL, '1d8 (1d10)', 0),
    ('Hammer, Great', 4, 20, 10, NULL, NULL, NULL, NULL, '2d6', 0),
    ('Javelin', 2, 3, 3, 20, 40, 80, '1/1', '1d4 (1d6)', 0),
    ('Lance', 5, 15, 8, NULL, NULL, NULL, NULL, '1d8', 0),
    ('Mace, Horseman''s', 1, 4, 3, NULL, NULL, NULL, NULL, '1d6', 0),
    ('Mace, Footman''s', 2, 10, 5, NULL, NULL, NULL, NULL, '1d8 (1d10)', 0),
    ('Mace, Great', 4, 20, 10, NULL, NULL, NULL, NULL, '2d6', 0),
    ('Morning Star', 2, 15, 5, NULL, NULL, NULL, NULL, '1d8 (1d10)', 0),
    ('Pick, Horseman''s', 1, 5, 3, NULL, NULL, NULL, NULL, '1d6', 0),
    ('Pick, War', 2, 15, 5, NULL, NULL, NULL, NULL, '1d8 (1d10)', 0),
    ('Pike', 6, 7, 12, NULL, NULL, NULL, NULL, '1d8', 0),
    ('Quarterstaff', 3, 5, 5, NULL, NULL, NULL, NULL, '1d6', 0);

-- +goose Down
DELETE FROM weapons 
WHERE name IN (
    'Axe, Hand', 'Axe, Battle', 'Axe, Great', 'Cæstuses', 'Chain Whip',
    'Club, Light', 'Club, War', 'Dagger', 'Dagger, Silver', 'Falcata',
    'Flail, Horseman''s', 'Flail, Footman''s', 'Halberd', 'Hammer, Horseman''s',
    'Hammer, War', 'Hammer, Great', 'Javelin', 'Lance', 'Mace, Horseman''s',
    'Mace, Footman''s', 'Mace, Great', 'Morning Star', 'Pick, Horseman''s', 
    'Pick, War', 'Pike', 'Quarterstaff'
);