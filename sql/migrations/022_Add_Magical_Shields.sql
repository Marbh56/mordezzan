-- +goose Up
INSERT INTO
    shields (name, cost_gp, weight, defense_bonus)
VALUES
    ('Small Shield +1', 2750, 5, 2), -- Base small shield bonus (1) + magic bonus (1)
    ('Small Shield +2', 4750, 5, 3), -- Base small shield bonus (1) + magic bonus (2)
    ('Small Shield +3', 7500, 5, 4), -- Base small shield bonus (1) + magic bonus (3)
    ('Large Shield +1', 3500, 10, 3), -- Base large shield bonus (2) + magic bonus (1)
    ('Large Shield +2', 7000, 10, 4), -- Base large shield bonus (2) + magic bonus (2)
    ('Large Shield +3', 10000, 10, 5);

-- Base large shield bonus (2) + magic bonus (3)
-- +goose Down
DELETE FROM shields
WHERE
    name IN (
        'Small Shield +1',
        'Small Shield +2',
        'Small Shield +3',
        'Large Shield +1',
        'Large Shield +2',
        'Large Shield +3'
    );
