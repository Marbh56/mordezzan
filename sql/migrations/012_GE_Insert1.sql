-- +goose Up
INSERT INTO
    equipment (name, cost_gp, weight, description)
VALUES
    ('Bandages, Gauze', 0.1, 0, '10-ft. ball'),
    (
        'Belladonna',
        30.0,
        0,
        'Dried bunch. 1:4 chance cures lycanthropy; 1:10 chance fatal pending death [poison] save'
    ),
    ('Blanket, Winter', 0.5, 3, NULL),
    ('Block and Tackle', 5.0, 5, NULL),
    ('Bow Case', 10.0, 1, 'Leather; waterproof'),
    (
        'Candle, Beeswax',
        0.1,
        0,
        '8-hour burn time, 5-ft. radius of light'
    ),
    (
        'Candle, Tallow',
        0.01,
        0,
        '2-hour burn time, 5-ft. radius of light'
    ),
    (
        'Chain, Iron, Heavy',
        5.0,
        3,
        'Half-inch thick, cost per foot'
    ),
    (
        'Chain, Iron, Light',
        3.0,
        1,
        'Quarter-inch thick, cost per foot'
    ),
    ('Chalk', 0.01, 0, 'Single piece'),
    (
        'Chisel',
        0.5,
        1,
        'Metal-, stone-, or wood-cutting'
    ),
    ('Cord, Sinew', 0.02, 0, '100-ft. ball'),
    (
        'Crampons',
        2.0,
        1,
        'Pair, includes ice axe: WC 1, 1d3 hp damage'
    ),
    ('Dice, Ivory', 0.2, 0, 'Pair'),
    ('Fishing Hooks', 0.5, 0, 'Set of 12'),
    ('Fishing Net', 3.0, 3, '10 × 10 ft.'),
    ('Fishing String', 0.01, 0, '100-ft. ball'),
    (
        'Glue',
        0.03,
        2,
        '1-qt. clay pot, powdered; must add hot water'
    ),
    ('Grappling Hook', 15.0, 3, 'Iron'),
    ('Grease', 0.02, 2, '1-qt. clay pot'),
    (
        'Hammer, Small',
        0.5,
        2,
        'Iron, WC 1, 1d2 hp damage'
    ),
    (
        'Helmet',
        10.0,
        2,
        'Normally included with armour purchase'
    ),
    ('Horn, Drinking', 0.1, 1, NULL),
    ('Horn, Hunting', 2.0, 1, NULL),
    ('Hourglass', 25.0, 1, 'Brass'),
    ('Ink and Quill', 10.0, 0, NULL),
    (
        'Lantern, Bulls-Eye',
        10.0,
        2,
        '15-ft. radius of light, 60-ft. narrow beam'
    ),
    (
        'Lantern, Hooded',
        7.0,
        2,
        '30-ft. radius of light'
    ),
    (
        'Marbles',
        0.2,
        0,
        'Set of 20, glass or ceramic; in soft leather pouch'
    );

-- +goose Down
DELETE FROM equipment
WHERE
    name IN (
        'Bandages, Gauze',
        'Belladonna',
        'Blanket, Winter',
        'Block and Tackle',
        'Bow Case',
        'Candle, Beeswax',
        'Candle, Tallow',
        'Chain, Iron, Heavy',
        'Chain, Iron, Light',
        'Chalk',
        'Chisel',
        'Cord, Sinew',
        'Crampons',
        'Dice, Ivory',
        'Fishing Hooks',
        'Fishing Net',
        'Fishing String',
        'Glue',
        'Grappling Hook',
        'Grease',
        'Hammer, Small',
        'Helmet',
        'Horn, Drinking',
        'Horn, Hunting',
        'Hourglass',
        'Ink and Quill',
        'Lantern, Bulls-Eye',
        'Lantern, Hooded',
        'Marbles'
    );
