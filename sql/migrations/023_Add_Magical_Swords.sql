-- +goose Up
-- Create table for magical weapons
CREATE TABLE magical_weapons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    base_weapon_id INTEGER NOT NULL,
    enhancement_bonus INTEGER NOT NULL,
    cost_gp INTEGER NOT NULL,
    xp_value INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (base_weapon_id) REFERENCES weapons (id)
);

-- Insert magical daggers
INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    1,
    1250,
    250
FROM
    weapons w
WHERE
    w.name = 'Dagger';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    2,
    3000,
    500
FROM
    weapons w
WHERE
    w.name = 'Dagger';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    3,
    10000,
    1000
FROM
    weapons w
WHERE
    w.name = 'Dagger';

-- Insert magical falcatas
INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    1,
    1750,
    350
FROM
    weapons w
WHERE
    w.name = 'Falcata';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    2,
    4000,
    650
FROM
    weapons w
WHERE
    w.name = 'Falcata';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    3,
    12500,
    1250
FROM
    weapons w
WHERE
    w.name = 'Falcata';

-- Insert magical short scimitars
INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    1,
    1750,
    350
FROM
    weapons w
WHERE
    w.name = 'Scimitar, Short';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    2,
    4000,
    650
FROM
    weapons w
WHERE
    w.name = 'Scimitar, Short';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    3,
    12500,
    1250
FROM
    weapons w
WHERE
    w.name = 'Scimitar, Short';

-- Insert magical long scimitars
INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    1,
    2250,
    450
FROM
    weapons w
WHERE
    w.name = 'Scimitar, Long';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    2,
    4800,
    800
FROM
    weapons w
WHERE
    w.name = 'Scimitar, Long';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    3,
    15000,
    1500
FROM
    weapons w
WHERE
    w.name = 'Scimitar, Long';

-- Insert magical two-handed scimitars
INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    1,
    3000,
    600
FROM
    weapons w
WHERE
    w.name = 'Scimitar, Two-handed';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    2,
    6000,
    1000
FROM
    weapons w
WHERE
    w.name = 'Scimitar, Two-handed';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    3,
    18000,
    1800
FROM
    weapons w
WHERE
    w.name = 'Scimitar, Two-handed';

-- Insert magical short swords
INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    1,
    1750,
    350
FROM
    weapons w
WHERE
    w.name = 'Sword, Short';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    2,
    4000,
    650
FROM
    weapons w
WHERE
    w.name = 'Sword, Short';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    3,
    12500,
    1250
FROM
    weapons w
WHERE
    w.name = 'Sword, Short';

-- Insert magical broad swords
INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    1,
    2250,
    450
FROM
    weapons w
WHERE
    w.name = 'Sword, Broad';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    2,
    4800,
    800
FROM
    weapons w
WHERE
    w.name = 'Sword, Broad';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    3,
    15000,
    1500
FROM
    weapons w
WHERE
    w.name = 'Sword, Broad';

-- Insert magical long swords
INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    1,
    2250,
    450
FROM
    weapons w
WHERE
    w.name = 'Sword, Long';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    2,
    4800,
    800
FROM
    weapons w
WHERE
    w.name = 'Sword, Long';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    3,
    15000,
    1500
FROM
    weapons w
WHERE
    w.name = 'Sword, Long';

-- Insert magical bastard swords
INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    1,
    2500,
    500
FROM
    weapons w
WHERE
    w.name = 'Sword, Bastard';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    2,
    5000,
    850
FROM
    weapons w
WHERE
    w.name = 'Sword, Bastard';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    3,
    16000,
    1600
FROM
    weapons w
WHERE
    w.name = 'Sword, Bastard';

-- Insert magical two-handed swords
INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    1,
    3000,
    600
FROM
    weapons w
WHERE
    w.name = 'Sword, Two-handed';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    2,
    6000,
    1000
FROM
    weapons w
WHERE
    w.name = 'Sword, Two-handed';

INSERT INTO
    magical_weapons (
        base_weapon_id,
        enhancement_bonus,
        cost_gp,
        xp_value
    )
SELECT
    w.id,
    3,
    18000,
    1800
FROM
    weapons w
WHERE
    w.name = 'Sword, Two-handed';

-- +goose Down
DROP TABLE IF EXISTS magical_weapons;
