-- +goose Up
INSERT INTO
    containers (
        name,
        cost_gp,
        weight,
        capacity_weight,
        capacity_items,
        description
    )
VALUES
    (
        'Arrow Quiver',
        1.0,
        1,
        2,
        12,
        'Leather quiver for storing arrows'
    ),
    ('Backpack', 5.0, 2, 50, NULL, 'Leather backpack'),
    (
        'Bolt Case',
        1.0,
        1,
        2,
        20,
        'Leather case for storing crossbow bolts'
    ),
    (
        'Map Case, Brass',
        5.0,
        0,
        1,
        6,
        'Brass case for storing maps or scrolls'
    ),
    (
        'Map Case, Leather',
        2.0,
        0,
        1,
        6,
        'Leather case for storing maps or scrolls'
    ),
    (
        'Pouch, Hard Leather',
        5.0,
        0,
        3,
        6,
        'Hard leather pouch with belt loops for potions or scrolls'
    ),
    (
        'Pouch, Soft Leather',
        0.07,
        0,
        5,
        NULL,
        'Soft leather pouch with drawstring'
    ),
    (
        'Sack, Large',
        0.3,
        0,
        40,
        NULL,
        'Large cloth or leather sack'
    ),
    (
        'Sack, Small',
        0.1,
        0,
        20,
        NULL,
        'Small cloth or leather sack'
    );

-- Add type restrictions
INSERT INTO
    container_allowed_types (container_id, item_type, ammo_type)
VALUES
    (
        (
            SELECT
                id
            FROM
                containers
            WHERE
                name = 'Arrow Quiver'
        ),
        'ammunition',
        'arrow'
    ),
    (
        (
            SELECT
                id
            FROM
                containers
            WHERE
                name = 'Bolt Case'
        ),
        'ammunition',
        'bolt'
    );

-- Add general container type permissions
INSERT INTO
    container_allowed_types (container_id, item_type)
SELECT
    c.id,
    'equipment'
FROM
    containers c
WHERE
    c.name IN (
        'Backpack',
        'Pouch, Hard Leather',
        'Pouch, Soft Leather',
        'Sack, Large',
        'Sack, Small'
    );

-- Add map/scroll specific permissions
INSERT INTO
    container_allowed_types (container_id, item_type)
SELECT
    c.id,
    'scroll'
FROM
    containers c
WHERE
    c.name IN ('Map Case, Brass', 'Map Case, Leather');

-- +goose Down
DELETE FROM container_allowed_types;

DELETE FROM containers;
