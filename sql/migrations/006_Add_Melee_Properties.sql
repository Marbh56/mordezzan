-- +goose Up
-- Chain Whip (↵)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Chain Whip'
    AND p.symbol = '↵';

-- Flail, Horseman's (↵)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Flail, Horseman''s'
    AND p.symbol = '↵';

-- Flail, Footman's (↵ +)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Flail, Footman''s'
    AND p.symbol IN ('↵', '+');

-- Halberd (+ ^ #)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Halberd'
    AND p.symbol IN ('+', '^', '#');

-- Hammer, Great (+ #)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Hammer, Great'
    AND p.symbol IN ('+', '#');

-- Lance (^ ∇ o)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Lance'
    AND p.symbol IN ('^', '∇', 'o');

-- Mace, Great (+ #)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Mace, Great'
    AND p.symbol IN ('+', '#');

-- Morning Star (Ω)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Morning Star'
    AND p.symbol = 'Ω';

-- Pick, Horseman's (Ω)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Pick, Horseman''s'
    AND p.symbol = 'Ω';

-- Pick, War (Ω)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Pick, War'
    AND p.symbol = 'Ω';

-- Pike (+ ^)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Pike'
    AND p.symbol IN ('+', '^');

-- Quarterstaff (↔)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Quarterstaff'
    AND p.symbol = '↔';

-- Scimitar, Two-handed (+)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Scimitar, Two-handed'
    AND p.symbol = '+';

-- Spear, Short (^)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Spear, Short'
    AND p.symbol = '^';

-- Spear, Long (^)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Spear, Long'
    AND p.symbol = '^';

-- Spear, Great (+ ^ ∇)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Spear, Great'
    AND p.symbol IN ('+', '^', '∇');

-- Spiked Staff (+ ^ #)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Spiked Staff'
    AND p.symbol IN ('+', '^', '#');

-- Sword, Two-handed (+)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Sword, Two-handed'
    AND p.symbol = '+';

-- Tonfa (↔)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Tonfa'
    AND p.symbol = '↔';

-- Trident, Hand (↔)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Trident, Hand'
    AND p.symbol = '↔';

-- Trident, Long (^)
INSERT INTO
    weapon_property_links
SELECT
    w.id,
    p.id
FROM
    weapons w,
    weapon_properties p
WHERE
    w.name = 'Trident, Long'
    AND p.symbol = '^';

-- +goose Down
DELETE FROM weapon_property_links;
