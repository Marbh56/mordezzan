-- +goose Up
INSERT INTO
    shields (name, cost_gp, weight, defense_bonus)
VALUES
    ('Small Shield +1', 2750, 5, 2), 
    ('Small Shield +2', 4750, 5, 3), 
    ('Small Shield +3', 7500, 5, 4), 
    ('Large Shield +1', 3500, 10, 3), 
    ('Large Shield +2', 7000, 10, 4), 
    ('Large Shield +3', 10000, 10, 5);

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
