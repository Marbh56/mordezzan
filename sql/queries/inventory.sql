-- name: GetCharacterInventory :many
SELECT
    ci.*,
    es.name as slot_name,
    CASE ci.item_type
        WHEN 'equipment' THEN e.name
        WHEN 'weapon' THEN w.name
        WHEN 'armor' THEN a.name
        WHEN 'ammunition' THEN am.name
        WHEN 'container' THEN c.name
        WHEN 'shield' THEN s.name
        WHEN 'ranged_weapon' THEN rw.name
    END as item_name,
    CASE ci.item_type
        WHEN 'equipment' THEN e.weight
        WHEN 'weapon' THEN w.weight
        WHEN 'armor' THEN a.weight
        WHEN 'ammunition' THEN am.weight
        WHEN 'container' THEN c.weight
        WHEN 'shield' THEN s.weight
        WHEN 'ranged_weapon' THEN rw.weight
        ELSE 0
    END as item_weight,
    CASE ci.item_type
        WHEN 'armor' THEN a.movement_rate
        ELSE NULL
    END as movement_rate,
    COALESCE(
        CASE ci.item_type
            WHEN 'weapon' THEN CAST(w.damage AS TEXT)
            WHEN 'ranged_weapon' THEN CAST(rw.damage AS TEXT)
            ELSE NULL
        END,
        NULL
    ) as damage,
    COALESCE(
        CASE ci.item_type
            WHEN 'weapon' THEN CAST(w.attacks_per_round AS TEXT)
            WHEN 'ranged_weapon' THEN CAST(rw.rate_of_fire AS TEXT)
            ELSE NULL
        END,
        NULL
    ) as attacks_per_round
FROM
    character_inventory ci
    LEFT JOIN equipment_slots es ON ci.equipment_slot_id = es.id
    LEFT JOIN equipment e ON ci.item_type = 'equipment'
    AND ci.item_id = e.id
    LEFT JOIN weapons w ON ci.item_type = 'weapon'
    AND ci.item_id = w.id
    LEFT JOIN armor a ON ci.item_type = 'armor'
    AND ci.item_id = a.id
    LEFT JOIN ammunition am ON ci.item_type = 'ammunition'
    AND ci.item_id = am.id
    LEFT JOIN containers c ON ci.item_type = 'container'
    AND ci.item_id = c.id
    LEFT JOIN shields s ON ci.item_type = 'shield'
    AND ci.item_id = s.id
    LEFT JOIN ranged_weapons rw ON ci.item_type = 'ranged_weapon'
    AND ci.item_id = rw.id
WHERE
    ci.character_id = ?
ORDER BY
    ci.container_inventory_id NULLS FIRST,
    ci.equipment_slot_id NULLS LAST,
    item_name;

-- name: AddItemToInventory :one
INSERT INTO
    character_inventory (
        character_id,
        item_type,
        item_id,
        quantity,
        container_inventory_id,
        equipment_slot_id,
        notes
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: GetItemFromInventory :one
SELECT
    *
FROM
    character_inventory
WHERE
    id = ?
    AND character_id = ?
LIMIT
    1;

-- name: UpdateInventoryItem :one
UPDATE character_inventory
SET
    quantity = ?,
    container_inventory_id = ?,
    equipment_slot_id = ?,
    notes = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = ?
    AND character_id = ? RETURNING *;

-- name: RemoveItemFromInventory :exec
DELETE FROM character_inventory
WHERE
    id = ?
    AND character_id = ?;

-- name: GetContainerContents :many
SELECT
    ci.*,
    CASE ci.item_type
        WHEN 'equipment' THEN e.name
        WHEN 'weapon' THEN w.name
        WHEN 'armor' THEN a.name
        WHEN 'ammunition' THEN am.name
        WHEN 'container' THEN c.name
        WHEN 'shield' THEN s.name
        WHEN 'ranged_weapon' THEN rw.name
    END as item_name,
    CASE ci.item_type
        WHEN 'equipment' THEN e.weight
        WHEN 'weapon' THEN w.weight
        WHEN 'armor' THEN a.weight
        WHEN 'ammunition' THEN am.weight
        WHEN 'container' THEN c.weight
        WHEN 'shield' THEN s.weight
        WHEN 'ranged_weapon' THEN rw.weight
        ELSE 0
    END as item_weight
FROM
    character_inventory ci
    LEFT JOIN equipment e ON ci.item_type = 'equipment'
    AND ci.item_id = e.id
    LEFT JOIN weapons w ON ci.item_type = 'weapon'
    AND ci.item_id = w.id
    LEFT JOIN armor a ON ci.item_type = 'armor'
    AND ci.item_id = a.id
    LEFT JOIN ammunition am ON ci.item_type = 'ammunition'
    AND ci.item_id = am.id
    LEFT JOIN containers c ON ci.item_type = 'container'
    AND ci.item_id = c.id
    LEFT JOIN shields s ON ci.item_type = 'shield'
    AND ci.item_id = s.id
    LEFT JOIN ranged_weapons rw ON ci.item_type = 'ranged_weapon'
    AND ci.item_id = rw.id
WHERE
    ci.container_inventory_id = ?
    AND ci.character_id = ?
ORDER BY
    item_name;

-- name: GetEquippedItems :many
SELECT
    ci.*,
    es.name as slot_name,
    CASE ci.item_type
        WHEN 'equipment' THEN e.name
        WHEN 'weapon' THEN w.name
        WHEN 'armor' THEN a.name
        WHEN 'ammunition' THEN am.name
        WHEN 'container' THEN c.name
        WHEN 'shield' THEN s.name
        WHEN 'ranged_weapon' THEN rw.name
    END as item_name,
    CASE ci.item_type
        WHEN 'equipment' THEN e.weight
        WHEN 'weapon' THEN w.weight
        WHEN 'armor' THEN a.weight
        WHEN 'ammunition' THEN am.weight
        WHEN 'container' THEN c.weight
        WHEN 'shield' THEN s.weight
        WHEN 'ranged_weapon' THEN rw.weight
        ELSE 0
    END as item_weight
FROM
    character_inventory ci
    JOIN equipment_slots es ON ci.equipment_slot_id = es.id
    LEFT JOIN equipment e ON ci.item_type = 'equipment'
    AND ci.item_id = e.id
    LEFT JOIN weapons w ON ci.item_type = 'weapon'
    AND ci.item_id = w.id
    LEFT JOIN armor a ON ci.item_type = 'armor'
    AND ci.item_id = a.id
    LEFT JOIN ammunition am ON ci.item_type = 'ammunition'
    AND ci.item_id = am.id
    LEFT JOIN containers c ON ci.item_type = 'container'
    AND ci.item_id = c.id
    LEFT JOIN shields s ON ci.item_type = 'shield'
    AND ci.item_id = s.id
    LEFT JOIN ranged_weapons rw ON ci.item_type = 'ranged_weapon'
    AND ci.item_id = rw.id
WHERE
    ci.character_id = ?
    AND ci.equipment_slot_id IS NOT NULL
ORDER BY
    es.name;
