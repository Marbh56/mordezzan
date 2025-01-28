-- +goose Up
INSERT INTO
    equipment (name, cost_gp, weight, description)
VALUES
    (
        'Mirror, Steel',
        5.0,
        1,
        'Polished steel; 3 × 5 inches'
    ),
    ('Mirror, Silver', 20.0, 1, '3 × 5 inches'),
    ('Nails', 0.1, 0, 'Set of 20, iron, 4-inch'),
    ('Needle, Blowgun', 0.1, 0, NULL),
    ('Needle, Sewing', 0.01, 0, NULL),
    (
        'Oil, Incendiary',
        35.0,
        1,
        '8-oz. flask, Greek fire'
    ),
    (
        'Oil, Lamp',
        0.1,
        1,
        '8-oz. flask, 6-hour burn time'
    ),
    ('Padlock and Key', 10.0, 1, NULL),
    ('Parchment', 1.0, 0, 'Single sheet'),
    ('Pole', 0.3, 5, 'Wooden; 10-ft.'),
    (
        'Pouch, Hard Leather',
        5.0,
        0,
        'Holds 6 potions or 3 scrolls; includes belt loops'
    ),
    (
        'Pouch, Soft Leather',
        0.07,
        0,
        'Holds 6–9 cubic inches of material; includes drawstring'
    ),
    ('Pry Bar', 1.0, 3, 'Iron; 16-inch'),
    ('Ring, Signet', 5.0, 0, 'Pewter'),
    ('Rope, Hemp', 1.0, 5, '50-ft., ½-inch'),
    ('Rope, Silk', 15.0, 2, '50-ft., ¼-inch'),
    (
        'Rope Ladder, Hemp',
        5.0,
        12,
        '50-ft., 2 parallel hemp ropes connected by short wooden crosspieces'
    ),
    (
        'Rope Ladder, Silk',
        50.0,
        6,
        '50-ft., 2 parallel silk ropes connected by short wooden crosspieces'
    ),
    (
        'Sack, Large',
        0.3,
        0,
        'Cloth or leather; 40-lb. capacity'
    ),
    (
        'Sack, Small',
        0.1,
        0,
        'Cloth or leather; 20-lb. capacity'
    ),
    (
        'Scabbard, Leather',
        0.5,
        0,
        'With baldric; normally included with sword purchase'
    ),
    ('Scabbard, Metal', 0.8, 1, 'With baldric'),
    (
        'Sheath, Dagger',
        0.3,
        0,
        'Leather; normally included with dagger purchase'
    ),
    ('Skis', 10.0, 8, 'Pair; includes poles'),
    ('Soap', 0.5, 1, 'Bar'),
    ('Spikes, Iron', 0.1, 1, 'Set of 4, 9-inch'),
    ('Spyglass', 750.0, 1, '×3 magnification'),
    (
        'Stakes and Wooden Mallet',
        1.0,
        2,
        'Set of 4 stakes with mallet'
    ),
    ('Tent, 1-Person', 5.0, 5, 'Canvas'),
    ('Tent, 2-Person', 7.0, 7, 'Canvas'),
    ('Tent, 4-Person', 12.0, 10, 'Canvas'),
    (
        'Thieves Tools',
        25.0,
        3,
        'File, oil dropper, picks, pincers, skeleton keys, small hammer, small saw, small wedge, wire'
    ),
    (
        'Tinderbox',
        2.0,
        1,
        'Contains flint and steel, paraffin, and wood powder'
    ),
    (
        'Torch',
        0.02,
        1,
        '1- to 2-hour burn time, 30-ft. radius of light; 1d4 hp damage as single-use weapon'
    ),
    ('Water-/Wineskin', 1.0, 0, '½-gallon capacity'),
    ('Wax, Bees-', 0.03, 1, 'Block'),
    ('Wire', 3.0, 0, '100-ft. spool, 50-lb. test'),
    (
        'Wolfsbane',
        25.0,
        0,
        'Dried bunch, 2:6 chance to drive off lycanthropes if affixed to spear tip'
    ),
    ('Writing Stick', 0.1, 0, 'Charcoal');

-- +goose Down
DELETE FROM equipment
WHERE
    name IN (
        'Mirror, Steel',
        'Mirror, Silver',
        'Nails',
        'Needle, Blowgun',
        'Needle, Sewing',
        'Oil, Incendiary',
        'Oil, Lamp',
        'Padlock and Key',
        'Parchment',
        'Pole',
        'Pouch, Hard Leather',
        'Pouch, Soft Leather',
        'Pry Bar',
        'Ring, Signet',
        'Rope, Hemp',
        'Rope, Silk',
        'Rope Ladder, Hemp',
        'Rope Ladder, Silk',
        'Sack, Large',
        'Sack, Small',
        'Scabbard, Leather',
        'Scabbard, Metal',
        'Sheath, Dagger',
        'Skis',
        'Soap',
        'Spikes, Iron',
        'Spyglass',
        'Stakes and Wooden Mallet',
        'Tent, 1-Person',
        'Tent, 2-Person',
        'Tent, 4-Person',
        'Thieves Tools',
        'Tinderbox',
        'Torch',
        'Water-/Wineskin',
        'Wax, Bees-',
        'Wire',
        'Wolfsbane',
        'Writing Stick'
    );
