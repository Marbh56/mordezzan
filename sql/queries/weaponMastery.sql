-- name: AddWeaponMastery :one
INSERT INTO
    character_weapon_masteries (character_id, weapon_id, mastery_level)
VALUES
    (?, ?, ?) RETURNING *;

-- name: GetCharacterWeaponMasteries :many
SELECT
    cwm.*,
    w.name as weapon_name,
    w.damage as base_damage,
    w.attacks_per_round as base_attacks
FROM
    character_weapon_masteries cwm
    JOIN weapons w ON cwm.weapon_id = w.id
WHERE
    cwm.character_id = ?
ORDER BY
    w.name;

-- name: UpdateWeaponMastery :one
UPDATE character_weapon_masteries
SET
    mastery_level = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE
    character_id = ?
    AND weapon_id = ? RETURNING *;

-- name: RemoveWeaponMastery :exec
DELETE FROM character_weapon_masteries
WHERE
    character_id = ?
    AND weapon_id = ?;

-- name: GetWeaponMastery :one
SELECT
    cwm.*,
    w.name as weapon_name,
    w.damage as base_damage,
    w.attacks_per_round as base_attacks
FROM
    character_weapon_masteries cwm
    JOIN weapons w ON cwm.weapon_id = w.id
WHERE
    cwm.character_id = ?
    AND cwm.weapon_id = ?
LIMIT
    1;

-- name: GetWeaponMasteriesForEquippedWeapons :many
SELECT
    cwm.weapon_id,
    cwm.mastery_level,
    w.name as weapon_name,
    w.damage as base_damage,
    w.attacks_per_round as base_attacks
FROM
    character_weapon_masteries cwm
    JOIN weapons w ON cwm.weapon_id = w.id
    JOIN character_inventory ci ON ci.item_id = cwm.weapon_id
    AND ci.item_type = 'weapon'
    AND ci.character_id = cwm.character_id
WHERE
    cwm.character_id = ?
    AND ci.equipment_slot_id IS NOT NULL;
