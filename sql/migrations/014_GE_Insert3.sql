-- +goose Up
INSERT INTO
    equipment (name, cost_gp, weight, description)
VALUES
    -- Clothing
    ('Belt', 0.2, 0, 'Leather'),
    ('Boots, Normal', 1.0, 1, 'Leather buskins'),
    ('Boots, Riding', 3.0, 1, NULL),
    ('Cape', 0.5, 1, 'Wool'),
    ('Cape, Fine', 25.0, 1, 'Fur, leather, or silk'),
    ('Cloak, Hooded', 0.8, 2, 'Wool'),
    (
        'Cloak, Hooded, Fine',
        50.0,
        2,
        'Fur, leather, or silk'
    ),
    (
        'Clothing, Disguise',
        25.0,
        2,
        'Faux exotic outfit, beggar''s clothes, wig, small jar of ochre, small jar of soot'
    ),
    (
        'Clothing, Normal',
        1.0,
        3,
        'Pantaloons, shirt/tunic, underclothes'
    ),
    (
        'Clothing, Religious',
        5.0,
        2,
        'Surplice and cassock; gown; etc.'
    ),
    (
        'Clothing, Special',
        15.0,
        4,
        'Buckskin outfit; fancy clothes; wool/fur winter outfit'
    ),
    ('Coat, Hooded, Fur', 30.0, 2, NULL),
    ('Coat, Hooded, Wool', 1.0, 2, NULL),
    ('Gloves, Fur', 20.0, 0, NULL),
    ('Gloves, Leather', 1.0, 0, NULL),
    ('Hat, Wool', 0.1, 0, NULL),
    ('Hat, Fur', 10.0, 0, NULL),
    ('Leggings, Fur', 10.0, 0, NULL),
    ('Robe', 1.0, 2, 'Wool'),
    ('Robe, Fine', 50.0, 2, 'Fur, silk, or velvet'),
    ('Sandals', 0.1, 0, 'Leather'),
    ('Shoes', 0.2, 1, 'Leather'),
    (
        'Tabard',
        1.0,
        0,
        'Wool; sleeveless jerkin; emblazoned'
    ),
    ('Toga', 0.1, 1, 'Wool'),
    -- Food and Provisions
    ('Biscuits, Hard', 0.1, 1, 'Bag'),
    (
        'Cereal',
        0.1,
        3,
        'Bag; barley, corn, oats, wheat'
    ),
    ('Cheese', 0.3, 3, 'Brick'),
    ('Eggs', 0.05, 1, 'Dozen; boxed'),
    ('Flour', 0.1, 20, 'Sack'),
    ('Honey', 1.0, 5, 'Crock'),
    ('Horse Meal/Grains', 0.5, 25, 'Sack'),
    ('Nuts', 0.5, 1, 'Bag'),
    (
        'Rations, Iron',
        5.0,
        5,
        '1 person, 1 week: salted/smoked meat or fish; dried fruit'
    ),
    (
        'Rations, Standard',
        2.0,
        10,
        '1 person, 1 week: cooked meat or fish; fresh fruit'
    ),
    ('Salt', 0.1, 1, 'Bag'),
    ('Spices, Cooking', 5.0, 0, 'Pouch'),
    ('Sugar, Coarse', 0.3, 1, 'Bag'),
    ('Wine, Common', 0.8, 2, 'Quart'),
    -- Religious and Musical Items
    ('Bell, Small', 0.5, 0, 'Brass'),
    ('Censer', 10.0, 2, 'Brass'),
    ('Flute', 0.1, 0, 'Wooden'),
    (
        'Holy Oil/Water',
        25.0,
        0,
        '8-oz. thin glass phial'
    ),
    ('Holy Symbol, Wooden', 0.5, 0, NULL),
    ('Holy Symbol, Silver', 25.0, 0, NULL),
    ('Holy Symbol, Ivory', 60.0, 0, 'Yellow ivory'),
    ('Holy Symbol, Gold', 75.0, 0, NULL),
    ('Incense Sticks', 5.0, 0, 'Set of 12'),
    ('Mask, Leather', 10.0, 0, 'Dyed leather'),
    ('Mask, Wooden', 10.0, 1, NULL),
    (
        'Mask, Pearl',
        65.0,
        1,
        'Wooden with mother-of-pearl'
    ),
    ('Mask, Silver', 75.0, 1, NULL),
    ('Mask, Ivory', 200.0, 1, 'Yellow ivory'),
    ('Mask, Gold', 250.0, 1, NULL),
    ('Paint, Body', 1.0, 3, 'Crock; e.g., ochre, woad'),
    ('Panpipes', 0.5, 0, NULL),
    ('Prayer Beads, Wooden', 0.01, 0, NULL),
    ('Prayer Beads, Ivory', 5.0, 0, 'Yellow ivory'),
    (
        'Prayer Book',
        100.0,
        1,
        'Leather cover, sewn binding, 50 parchment pages'
    ),
    ('Rattle', 1.0, 0, 'Wooden');

-- +goose Down
DELETE FROM equipment
WHERE
    name IN (
        'Belt',
        'Boots, Normal',
        'Boots, Riding',
        'Cape',
        'Cape, Fine',
        'Cloak, Hooded',
        'Cloak, Hooded, Fine',
        'Clothing, Disguise',
        'Clothing, Normal',
        'Clothing, Religious',
        'Clothing, Special',
        'Coat, Hooded, Fur',
        'Coat, Hooded, Wool',
        'Gloves, Fur',
        'Gloves, Leather',
        'Hat, Wool',
        'Hat, Fur',
        'Leggings, Fur',
        'Robe',
        'Robe, Fine',
        'Sandals',
        'Shoes',
        'Tabard',
        'Toga',
        'Biscuits, Hard',
        'Cereal',
        'Cheese',
        'Eggs',
        'Flour',
        'Honey',
        'Horse Meal/Grains',
        'Nuts',
        'Rations, Iron',
        'Rations, Standard',
        'Salt',
        'Spices, Cooking',
        'Sugar, Coarse',
        'Wine, Common',
        'Bell, Small',
        'Censer',
        'Flute',
        'Holy Oil/Water',
        'Holy Symbol, Wooden',
        'Holy Symbol, Silver',
        'Holy Symbol, Ivory',
        'Holy Symbol, Gold',
        'Incense Sticks',
        'Mask, Leather',
        'Mask, Wooden',
        'Mask, Pearl',
        'Mask, Silver',
        'Mask, Ivory',
        'Mask, Gold',
        'Paint, Body',
        'Panpipes',
        'Prayer Beads, Wooden',
        'Prayer Beads, Ivory',
        'Prayer Book',
        'Rattle'
    );
