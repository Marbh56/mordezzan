-- name: GetMagicalWeapons :many
SELECT
    mw.id as magical_weapon_id,
    mw.base_weapon_id,
    mw.enhancement_bonus,
    mw.cost_gp,
    mw.xp_value,
    mw.created_at,
    mw.updated_at,
    w.id as weapon_id,
    w.name as base_weapon_name,
    w.damage as base_damage,
    w.attacks_per_round as base_attacks,
    w.weight
FROM
    magical_weapons mw
    JOIN weapons w ON mw.base_weapon_id = w.id
ORDER BY
    w.name,
    mw.enhancement_bonus;

-- name: GetMagicalWeapon :one
SELECT
    mw.id as magical_weapon_id,
    mw.base_weapon_id,
    mw.enhancement_bonus,
    mw.cost_gp,
    mw.xp_value,
    mw.created_at,
    mw.updated_at,
    w.id as weapon_id,
    w.name as base_weapon_name,
    w.damage as base_damage,
    w.attacks_per_round as base_attacks,
    w.weight
FROM
    magical_weapons mw
    JOIN weapons w ON mw.base_weapon_id = w.id
WHERE
    mw.id = ?;

-- name: AddMagicalWeaponToInventory :one
INSERT INTO
    character_inventory (
        character_id,
        item_type,
        item_id,
        magical_weapon_id,
        quantity,
        container_inventory_id,
        equipment_slot_id,
        notes
    )
VALUES
    (
        ?, -- character_id
        'weapon',
        (
            SELECT
                base_weapon_id
            FROM
                magical_weapons
            WHERE
                magical_weapons.id = ?
        ), -- item_id from magical weapon
        ?, -- magical_weapon_id
        ?, -- quantity
        ?, -- container_inventory_id
        ?, -- equipment_slot_id
        ? -- notes
    ) RETURNING *;
