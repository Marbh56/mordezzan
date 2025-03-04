-- name: GetCharacterInventoryItems :many
SELECT 
    ci.id, ci.character_id, ci.item_id, ci.item_type, ci.quantity,
    ci.container_id, ci.equipment_slot_id, ci.notes,
    ci.created_at, ci.updated_at,
    CASE 
        WHEN ci.item_type = 'equipment' THEN e.name
        WHEN ci.item_type = 'weapon' THEN w.name
        WHEN ci.item_type = 'armor' THEN a.name
        WHEN ci.item_type = 'ammunition' THEN am.name
        WHEN ci.item_type = 'container' THEN e.name  -- Use equipment name as fallback
        WHEN ci.item_type = 'shield' THEN s.name
        WHEN ci.item_type = 'ranged_weapon' THEN rw.name
    END as item_name,
    CASE 
        WHEN ci.item_type = 'equipment' THEN e.weight
        WHEN ci.item_type = 'weapon' THEN w.weight
        WHEN ci.item_type = 'armor' THEN a.weight
        WHEN ci.item_type = 'ammunition' THEN am.weight
        WHEN ci.item_type = 'container' THEN e.weight  -- Use equipment weight as fallback
        WHEN ci.item_type = 'shield' THEN s.weight
        WHEN ci.item_type = 'ranged_weapon' THEN rw.weight
        ELSE 0
    END as item_weight,
    CASE 
        WHEN ci.item_type = 'shield' THEN s.defense_bonus
        ELSE NULL
    END as defense_bonus,
    CASE 
        WHEN ci.item_type = 'weapon' THEN w.damage
        WHEN ci.item_type = 'ranged_weapon' THEN rw.damage
        ELSE NULL
    END as damage,
    CASE 
        WHEN ci.item_type = 'weapon' THEN w.attacks_per_round
        WHEN ci.item_type = 'ranged_weapon' THEN rw.rate_of_fire
        ELSE NULL
    END as attacks_per_round,
    CASE 
        WHEN ci.item_type = 'armor' THEN a.movement_rate
        ELSE NULL
    END as movement_rate,
    es.name as slot_name,
    CASE 
        WHEN ci.item_type = 'container' THEN c.capacity_weight
        ELSE NULL
    END as container_capacity,
    CASE 
        WHEN ci.item_type = 'container' THEN c.capacity_items
        ELSE NULL
    END as container_max_items
FROM 
    character_inventory ci
    LEFT JOIN equipment_slots es ON ci.equipment_slot_id = es.id
    LEFT JOIN equipment e ON ci.item_type = 'equipment' AND ci.item_id = e.id
    LEFT JOIN weapons w ON ci.item_type = 'weapon' AND ci.item_id = w.id
    LEFT JOIN armor a ON ci.item_type = 'armor' AND ci.item_id = a.id
    LEFT JOIN ammunition am ON ci.item_type = 'ammunition' AND ci.item_id = am.id
    LEFT JOIN containers c ON ci.item_type = 'container' AND ci.item_id = c.id
    LEFT JOIN shields s ON ci.item_type = 'shield' AND ci.item_id = s.id
    LEFT JOIN ranged_weapons rw ON ci.item_type = 'ranged_weapon' AND ci.item_id = rw.id
WHERE 
    ci.character_id = ?
ORDER BY 
    ci.equipment_slot_id IS NULL, 
    es.name,
    ci.container_id IS NOT NULL,
    item_name;

-- name: GetContainerContents :many
SELECT 
    ci.id, ci.character_id, ci.item_id, ci.item_type, ci.quantity,
    ci.container_id, ci.equipment_slot_id, ci.notes,
    ci.created_at, ci.updated_at,
    CASE 
        WHEN ci.item_type = 'equipment' THEN e.name
        WHEN ci.item_type = 'weapon' THEN w.name
        WHEN ci.item_type = 'armor' THEN a.name
        WHEN ci.item_type = 'ammunition' THEN am.name
        WHEN ci.item_type = 'container' THEN e.name  -- Use equipment name as fallback
        WHEN ci.item_type = 'shield' THEN s.name
        WHEN ci.item_type = 'ranged_weapon' THEN rw.name
    END as item_name,
    CASE 
        WHEN ci.item_type = 'equipment' THEN e.weight
        WHEN ci.item_type = 'weapon' THEN w.weight
        WHEN ci.item_type = 'armor' THEN a.weight
        WHEN ci.item_type = 'ammunition' THEN am.weight
        WHEN ci.item_type = 'container' THEN e.weight  -- Use equipment weight as fallback
        WHEN ci.item_type = 'shield' THEN s.weight
        WHEN ci.item_type = 'ranged_weapon' THEN rw.weight
        ELSE 0
    END as item_weight
FROM 
    character_inventory ci
    LEFT JOIN equipment e ON ci.item_type = 'equipment' AND ci.item_id = e.id
    LEFT JOIN weapons w ON ci.item_type = 'weapon' AND ci.item_id = w.id
    LEFT JOIN armor a ON ci.item_type = 'armor' AND ci.item_id = a.id
    LEFT JOIN ammunition am ON ci.item_type = 'ammunition' AND ci.item_id = am.id
    LEFT JOIN containers c ON ci.item_type = 'container' AND ci.item_id = c.id
    LEFT JOIN shields s ON ci.item_type = 'shield' AND ci.item_id = s.id
    LEFT JOIN ranged_weapons rw ON ci.item_type = 'ranged_weapon' AND ci.item_id = rw.id
WHERE 
    ci.container_id = ?
    AND ci.character_id = ?
ORDER BY 
    item_name;

-- name: GetContainerWeight :one
SELECT 
    COALESCE(SUM(
        CASE 
            WHEN ci.item_type = 'equipment' THEN e.weight * ci.quantity
            WHEN ci.item_type = 'weapon' THEN w.weight * ci.quantity
            WHEN ci.item_type = 'armor' THEN a.weight * ci.quantity
            WHEN ci.item_type = 'ammunition' THEN COALESCE(am.weight, 0) * ci.quantity
            WHEN ci.item_type = 'container' THEN COALESCE(e.weight, 0) * ci.quantity  -- Use equipment weight as fallback
            WHEN ci.item_type = 'shield' THEN s.weight * ci.quantity
            WHEN ci.item_type = 'ranged_weapon' THEN rw.weight * ci.quantity
            ELSE 0
        END
    ), 0) as total_weight
FROM 
    character_inventory ci
    LEFT JOIN equipment e ON ci.item_type = 'equipment' AND ci.item_id = e.id
    LEFT JOIN weapons w ON ci.item_type = 'weapon' AND ci.item_id = w.id
    LEFT JOIN armor a ON ci.item_type = 'armor' AND ci.item_id = a.id
    LEFT JOIN ammunition am ON ci.item_type = 'ammunition' AND ci.item_id = am.id
    LEFT JOIN containers c ON ci.item_type = 'container' AND ci.item_id = c.id
    LEFT JOIN shields s ON ci.item_type = 'shield' AND ci.item_id = s.id
    LEFT JOIN ranged_weapons rw ON ci.item_type = 'ranged_weapon' AND ci.item_id = rw.id
WHERE 
    ci.container_id = ?;

-- name: GetContainerCapacity :one
SELECT 
    c.capacity_weight,
    c.capacity_items
FROM 
    character_inventory ci
JOIN 
    containers c ON ci.item_id = c.base_item_id
WHERE 
    ci.id = ?
    AND ci.item_type = 'container';

-- name: AddItemToInventory :one
INSERT INTO
    character_inventory (
        character_id,
        item_id,
        item_type,
        quantity,
        container_id,
        equipment_slot_id,
        notes
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?) RETURNING id,
    character_id,
    item_id,
    item_type,
    quantity,
    container_id,
    equipment_slot_id,
    notes,
    created_at,
    updated_at;

-- name: UpdateInventoryItem :one
UPDATE character_inventory
SET
    quantity = ?,
    container_id = ?,
    equipment_slot_id = ?,
    notes = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = ?
    AND character_id = ? RETURNING id, character_id, item_id, item_type, quantity, container_id, equipment_slot_id, notes, created_at, updated_at;

-- name: RemoveItemFromInventory :exec
DELETE FROM character_inventory
WHERE
    id = ?
    AND character_id = ?;

-- name: IsSlotOccupied :one
SELECT 
    COUNT(*) > 0 as is_occupied
FROM 
    character_inventory
WHERE 
    character_id = ? 
    AND equipment_slot_id = ?;

-- name: MoveItemToContainer :exec
UPDATE character_inventory
SET 
    container_id = ?,
    equipment_slot_id = NULL
WHERE 
    id = ?
    AND character_id = ?;

-- name: EquipItem :exec
UPDATE character_inventory
SET 
    equipment_slot_id = ?,
    container_id = NULL
WHERE 
    id = ?
    AND character_id = ?;

-- name: UnequipItem :exec
UPDATE character_inventory
SET 
    equipment_slot_id = NULL
WHERE 
    id = ?
    AND character_id = ?;
    
-- name: UpdateItemQuantity :exec
UPDATE character_inventory
SET 
    quantity = ?
WHERE 
    id = ?
    AND character_id = ?;

-- name: FindStackableItemInInventory :one
SELECT 
    id, quantity
FROM 
    character_inventory
WHERE 
    character_id = ?
    AND item_id = ?
    AND item_type = ?
    AND container_id IS NULL
    AND equipment_slot_id IS NULL
LIMIT 1;

-- name: FindStackableItemInContainer :one
SELECT 
    id, quantity
FROM 
    character_inventory
WHERE 
    character_id = ?
    AND item_id = ?
    AND item_type = ?
    AND container_id = ?
LIMIT 1;

-- name: SplitStack :exec
INSERT INTO character_inventory (
    character_id, 
    item_id,
    item_type,
    quantity,
    container_id,
    equipment_slot_id,
    notes,
    created_at,
    updated_at
)
SELECT 
    ci.character_id,
    ci.item_id,
    ci.item_type,
    ?,
    ci.container_id,
    NULL,
    ci.notes,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM 
    character_inventory ci
WHERE 
    ci.id = ?;
    
-- name: ReduceStackQuantity :exec
UPDATE character_inventory
SET quantity = quantity - ?
WHERE id = ? AND character_id = ?;

-- name: GetEquippedItems :many
SELECT
    ci.id, ci.character_id, ci.item_id, ci.item_type, ci.quantity,
    ci.container_id, ci.equipment_slot_id, ci.notes, ci.created_at, ci.updated_at,
    es.name as slot_name,
    CASE 
        WHEN ci.item_type = 'equipment' THEN e.name
        WHEN ci.item_type = 'weapon' THEN w.name
        WHEN ci.item_type = 'armor' THEN a.name
        WHEN ci.item_type = 'ammunition' THEN am.name
        WHEN ci.item_type = 'container' THEN c.name
        WHEN ci.item_type = 'shield' THEN s.name
        WHEN ci.item_type = 'ranged_weapon' THEN rw.name
    END as item_name,
    CASE 
        WHEN ci.item_type = 'equipment' THEN e.weight
        WHEN ci.item_type = 'weapon' THEN w.weight
        WHEN ci.item_type = 'armor' THEN a.weight
        WHEN ci.item_type = 'ammunition' THEN am.weight
        WHEN ci.item_type = 'container' THEN c.weight
        WHEN ci.item_type = 'shield' THEN s.weight
        WHEN ci.item_type = 'ranged_weapon' THEN rw.weight
        ELSE 0
    END as item_weight
FROM
    character_inventory ci
    JOIN equipment_slots es ON ci.equipment_slot_id = es.id
    LEFT JOIN equipment e ON ci.item_type = 'equipment' AND ci.item_id = e.id
    LEFT JOIN weapons w ON ci.item_type = 'weapon' AND ci.item_id = w.id
    LEFT JOIN armor a ON ci.item_type = 'armor' AND ci.item_id = a.id
    LEFT JOIN ammunition am ON ci.item_type = 'ammunition' AND ci.item_id = am.id
    LEFT JOIN containers c ON ci.item_type = 'container' AND ci.item_id = c.id
    LEFT JOIN shields s ON ci.item_type = 'shield' AND ci.item_id = s.id
    LEFT JOIN ranged_weapons rw ON ci.item_type = 'ranged_weapon' AND ci.item_id = rw.id
WHERE
    ci.character_id = ?
    AND ci.equipment_slot_id IS NOT NULL
ORDER BY
    es.name;


-- name: GetAllWeapons :many
SELECT 
    id, name, weight, cost_gp
FROM 
    weapons
ORDER BY 
    name;

-- name: GetAllArmor :many
SELECT 
    id, name, weight, cost_gp
FROM 
    armor
ORDER BY 
    name;

-- name: GetAllShields :many
SELECT 
    id, name, weight, cost_gp
FROM 
    shields
ORDER BY 
    name;

-- name: GetAllEquipment :many
SELECT 
    id, name, weight, cost_gp
FROM 
    equipment
ORDER BY 
    name;

-- name: GetAllAmmunition :many
SELECT 
    id, name, weight, cost_gp
FROM 
    ammunition
ORDER BY 
    name;

-- name: GetAllRangedWeapons :many
SELECT 
    id, name, weight, cost_gp
FROM 
    ranged_weapons
ORDER BY 
    name;

-- name: GetEquipmentSlots :many
SELECT 
    id, name, description
FROM 
    equipment_slots
ORDER BY 
    name;