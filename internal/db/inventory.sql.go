// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: inventory.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const addItemToInventory = `-- name: AddItemToInventory :one
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
    updated_at
`

type AddItemToInventoryParams struct {
	CharacterID     int64          `json:"character_id"`
	ItemID          int64          `json:"item_id"`
	ItemType        string         `json:"item_type"`
	Quantity        int64          `json:"quantity"`
	ContainerID     sql.NullInt64  `json:"container_id"`
	EquipmentSlotID sql.NullInt64  `json:"equipment_slot_id"`
	Notes           sql.NullString `json:"notes"`
}

type AddItemToInventoryRow struct {
	ID              int64          `json:"id"`
	CharacterID     int64          `json:"character_id"`
	ItemID          int64          `json:"item_id"`
	ItemType        string         `json:"item_type"`
	Quantity        int64          `json:"quantity"`
	ContainerID     sql.NullInt64  `json:"container_id"`
	EquipmentSlotID sql.NullInt64  `json:"equipment_slot_id"`
	Notes           sql.NullString `json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (q *Queries) AddItemToInventory(ctx context.Context, arg AddItemToInventoryParams) (AddItemToInventoryRow, error) {
	row := q.db.QueryRowContext(ctx, addItemToInventory,
		arg.CharacterID,
		arg.ItemID,
		arg.ItemType,
		arg.Quantity,
		arg.ContainerID,
		arg.EquipmentSlotID,
		arg.Notes,
	)
	var i AddItemToInventoryRow
	err := row.Scan(
		&i.ID,
		&i.CharacterID,
		&i.ItemID,
		&i.ItemType,
		&i.Quantity,
		&i.ContainerID,
		&i.EquipmentSlotID,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const addMagicalItemToInventory = `-- name: AddMagicalItemToInventory :one
INSERT INTO
    character_inventory (
        character_id,
        item_id,
        item_type,
        quantity,
        container_id,
        equipment_slot_id,
        charges,
        notes
    )
VALUES
    (?, ?, 'magical_item', 1, ?, ?, ?, ?) 
RETURNING id,
    character_id,
    item_id,
    item_type,
    quantity,
    container_id,
    equipment_slot_id,
    charges,
    notes,
    created_at,
    updated_at
`

type AddMagicalItemToInventoryParams struct {
	CharacterID     int64          `json:"character_id"`
	ItemID          int64          `json:"item_id"`
	ContainerID     sql.NullInt64  `json:"container_id"`
	EquipmentSlotID sql.NullInt64  `json:"equipment_slot_id"`
	Charges         sql.NullInt64  `json:"charges"`
	Notes           sql.NullString `json:"notes"`
}

type AddMagicalItemToInventoryRow struct {
	ID              int64          `json:"id"`
	CharacterID     int64          `json:"character_id"`
	ItemID          int64          `json:"item_id"`
	ItemType        string         `json:"item_type"`
	Quantity        int64          `json:"quantity"`
	ContainerID     sql.NullInt64  `json:"container_id"`
	EquipmentSlotID sql.NullInt64  `json:"equipment_slot_id"`
	Charges         sql.NullInt64  `json:"charges"`
	Notes           sql.NullString `json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (q *Queries) AddMagicalItemToInventory(ctx context.Context, arg AddMagicalItemToInventoryParams) (AddMagicalItemToInventoryRow, error) {
	row := q.db.QueryRowContext(ctx, addMagicalItemToInventory,
		arg.CharacterID,
		arg.ItemID,
		arg.ContainerID,
		arg.EquipmentSlotID,
		arg.Charges,
		arg.Notes,
	)
	var i AddMagicalItemToInventoryRow
	err := row.Scan(
		&i.ID,
		&i.CharacterID,
		&i.ItemID,
		&i.ItemType,
		&i.Quantity,
		&i.ContainerID,
		&i.EquipmentSlotID,
		&i.Charges,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const equipItem = `-- name: EquipItem :exec
UPDATE character_inventory
SET 
    equipment_slot_id = ?,
    container_id = NULL
WHERE 
    id = ?
    AND character_id = ?
`

type EquipItemParams struct {
	EquipmentSlotID sql.NullInt64 `json:"equipment_slot_id"`
	ID              int64         `json:"id"`
	CharacterID     int64         `json:"character_id"`
}

func (q *Queries) EquipItem(ctx context.Context, arg EquipItemParams) error {
	_, err := q.db.ExecContext(ctx, equipItem, arg.EquipmentSlotID, arg.ID, arg.CharacterID)
	return err
}

const findStackableItemInContainer = `-- name: FindStackableItemInContainer :one
SELECT 
    id, quantity
FROM 
    character_inventory
WHERE 
    character_id = ?
    AND item_id = ?
    AND item_type = ?
    AND container_id = ?
LIMIT 1
`

type FindStackableItemInContainerParams struct {
	CharacterID int64         `json:"character_id"`
	ItemID      int64         `json:"item_id"`
	ItemType    string        `json:"item_type"`
	ContainerID sql.NullInt64 `json:"container_id"`
}

type FindStackableItemInContainerRow struct {
	ID       int64 `json:"id"`
	Quantity int64 `json:"quantity"`
}

func (q *Queries) FindStackableItemInContainer(ctx context.Context, arg FindStackableItemInContainerParams) (FindStackableItemInContainerRow, error) {
	row := q.db.QueryRowContext(ctx, findStackableItemInContainer,
		arg.CharacterID,
		arg.ItemID,
		arg.ItemType,
		arg.ContainerID,
	)
	var i FindStackableItemInContainerRow
	err := row.Scan(&i.ID, &i.Quantity)
	return i, err
}

const findStackableItemInInventory = `-- name: FindStackableItemInInventory :one
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
LIMIT 1
`

type FindStackableItemInInventoryParams struct {
	CharacterID int64  `json:"character_id"`
	ItemID      int64  `json:"item_id"`
	ItemType    string `json:"item_type"`
}

type FindStackableItemInInventoryRow struct {
	ID       int64 `json:"id"`
	Quantity int64 `json:"quantity"`
}

func (q *Queries) FindStackableItemInInventory(ctx context.Context, arg FindStackableItemInInventoryParams) (FindStackableItemInInventoryRow, error) {
	row := q.db.QueryRowContext(ctx, findStackableItemInInventory, arg.CharacterID, arg.ItemID, arg.ItemType)
	var i FindStackableItemInInventoryRow
	err := row.Scan(&i.ID, &i.Quantity)
	return i, err
}

const getAllAmmunition = `-- name: GetAllAmmunition :many
SELECT 
    id, name, weight, cost_gp
FROM 
    ammunition
ORDER BY 
    name
`

type GetAllAmmunitionRow struct {
	ID     int64         `json:"id"`
	Name   string        `json:"name"`
	Weight sql.NullInt64 `json:"weight"`
	CostGp float64       `json:"cost_gp"`
}

func (q *Queries) GetAllAmmunition(ctx context.Context) ([]GetAllAmmunitionRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllAmmunition)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllAmmunitionRow
	for rows.Next() {
		var i GetAllAmmunitionRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Weight,
			&i.CostGp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllArmor = `-- name: GetAllArmor :many
SELECT 
    id, name, weight, cost_gp
FROM 
    armor
ORDER BY 
    name
`

type GetAllArmorRow struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Weight int64  `json:"weight"`
	CostGp int64  `json:"cost_gp"`
}

func (q *Queries) GetAllArmor(ctx context.Context) ([]GetAllArmorRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllArmor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllArmorRow
	for rows.Next() {
		var i GetAllArmorRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Weight,
			&i.CostGp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllEquipment = `-- name: GetAllEquipment :many
SELECT 
    id, name, weight, cost_gp
FROM 
    equipment
ORDER BY 
    name
`

type GetAllEquipmentRow struct {
	ID     int64         `json:"id"`
	Name   string        `json:"name"`
	Weight sql.NullInt64 `json:"weight"`
	CostGp float64       `json:"cost_gp"`
}

func (q *Queries) GetAllEquipment(ctx context.Context) ([]GetAllEquipmentRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllEquipment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllEquipmentRow
	for rows.Next() {
		var i GetAllEquipmentRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Weight,
			&i.CostGp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllMagicalItems = `-- name: GetAllMagicalItems :many
SELECT 
    id, name, description, weight, cost_gp, max_charges, category, effect_description
FROM 
    magical_items
ORDER BY 
    name
`

type GetAllMagicalItemsRow struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Weight            int64  `json:"weight"`
	CostGp            int64  `json:"cost_gp"`
	MaxCharges        int64  `json:"max_charges"`
	Category          string `json:"category"`
	EffectDescription string `json:"effect_description"`
}

func (q *Queries) GetAllMagicalItems(ctx context.Context) ([]GetAllMagicalItemsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllMagicalItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllMagicalItemsRow
	for rows.Next() {
		var i GetAllMagicalItemsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Weight,
			&i.CostGp,
			&i.MaxCharges,
			&i.Category,
			&i.EffectDescription,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllRangedWeapons = `-- name: GetAllRangedWeapons :many
SELECT 
    id, name, weight, cost_gp
FROM 
    ranged_weapons
ORDER BY 
    name
`

type GetAllRangedWeaponsRow struct {
	ID     int64         `json:"id"`
	Name   string        `json:"name"`
	Weight int64         `json:"weight"`
	CostGp sql.NullInt64 `json:"cost_gp"`
}

func (q *Queries) GetAllRangedWeapons(ctx context.Context) ([]GetAllRangedWeaponsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllRangedWeapons)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllRangedWeaponsRow
	for rows.Next() {
		var i GetAllRangedWeaponsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Weight,
			&i.CostGp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllShields = `-- name: GetAllShields :many
SELECT 
    id, name, weight, cost_gp
FROM 
    shields
ORDER BY 
    name
`

type GetAllShieldsRow struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Weight int64  `json:"weight"`
	CostGp int64  `json:"cost_gp"`
}

func (q *Queries) GetAllShields(ctx context.Context) ([]GetAllShieldsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllShields)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllShieldsRow
	for rows.Next() {
		var i GetAllShieldsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Weight,
			&i.CostGp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllWeapons = `-- name: GetAllWeapons :many
SELECT 
    id, name, weight, cost_gp
FROM 
    weapons
ORDER BY 
    name
`

type GetAllWeaponsRow struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Weight int64  `json:"weight"`
	CostGp int64  `json:"cost_gp"`
}

func (q *Queries) GetAllWeapons(ctx context.Context) ([]GetAllWeaponsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllWeapons)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllWeaponsRow
	for rows.Next() {
		var i GetAllWeaponsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Weight,
			&i.CostGp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCharacterInventoryItems = `-- name: GetCharacterInventoryItems :many
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
    CASE 
        WHEN ci.item_type = 'armor' THEN a.armor_class
        ELSE NULL
    END as armor_class,
    CASE
        WHEN ci.item_type = 'weapon' THEN w.enhancement_bonus
        WHEN ci.item_type = 'armor' THEN a.enhancement_bonus
        WHEN ci.item_type = 'shield' THEN s.enhancement_bonus
        WHEN ci.item_type = 'ranged_weapon' THEN rw.enhancement_bonus
        WHEN ci.item_type = 'ammunition' THEN am.enhancement_bonus
        ELSE NULL
    END as enhancement_bonus,
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
    item_name
`

type GetCharacterInventoryItemsRow struct {
	ID                int64          `json:"id"`
	CharacterID       int64          `json:"character_id"`
	ItemID            int64          `json:"item_id"`
	ItemType          string         `json:"item_type"`
	Quantity          int64          `json:"quantity"`
	ContainerID       sql.NullInt64  `json:"container_id"`
	EquipmentSlotID   sql.NullInt64  `json:"equipment_slot_id"`
	Notes             sql.NullString `json:"notes"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	ItemName          interface{}    `json:"item_name"`
	ItemWeight        int64          `json:"item_weight"`
	DefenseBonus      interface{}    `json:"defense_bonus"`
	Damage            interface{}    `json:"damage"`
	AttacksPerRound   interface{}    `json:"attacks_per_round"`
	MovementRate      interface{}    `json:"movement_rate"`
	ArmorClass        interface{}    `json:"armor_class"`
	EnhancementBonus  interface{}    `json:"enhancement_bonus"`
	SlotName          sql.NullString `json:"slot_name"`
	ContainerCapacity interface{}    `json:"container_capacity"`
	ContainerMaxItems interface{}    `json:"container_max_items"`
}

func (q *Queries) GetCharacterInventoryItems(ctx context.Context, characterID int64) ([]GetCharacterInventoryItemsRow, error) {
	rows, err := q.db.QueryContext(ctx, getCharacterInventoryItems, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCharacterInventoryItemsRow
	for rows.Next() {
		var i GetCharacterInventoryItemsRow
		if err := rows.Scan(
			&i.ID,
			&i.CharacterID,
			&i.ItemID,
			&i.ItemType,
			&i.Quantity,
			&i.ContainerID,
			&i.EquipmentSlotID,
			&i.Notes,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ItemName,
			&i.ItemWeight,
			&i.DefenseBonus,
			&i.Damage,
			&i.AttacksPerRound,
			&i.MovementRate,
			&i.ArmorClass,
			&i.EnhancementBonus,
			&i.SlotName,
			&i.ContainerCapacity,
			&i.ContainerMaxItems,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getChargedItem = `-- name: GetChargedItem :one
SELECT 
    ci.id, ci.character_id, ci.charges, 
    mi.name, mi.description, mi.effect_description, mi.category
FROM 
    character_inventory ci
JOIN 
    magical_items mi ON ci.item_id = mi.id
WHERE 
    ci.id = ?
    AND ci.character_id = ?
    AND ci.item_type = 'magical_item'
`

type GetChargedItemParams struct {
	ID          int64 `json:"id"`
	CharacterID int64 `json:"character_id"`
}

type GetChargedItemRow struct {
	ID                int64         `json:"id"`
	CharacterID       int64         `json:"character_id"`
	Charges           sql.NullInt64 `json:"charges"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	EffectDescription string        `json:"effect_description"`
	Category          string        `json:"category"`
}

func (q *Queries) GetChargedItem(ctx context.Context, arg GetChargedItemParams) (GetChargedItemRow, error) {
	row := q.db.QueryRowContext(ctx, getChargedItem, arg.ID, arg.CharacterID)
	var i GetChargedItemRow
	err := row.Scan(
		&i.ID,
		&i.CharacterID,
		&i.Charges,
		&i.Name,
		&i.Description,
		&i.EffectDescription,
		&i.Category,
	)
	return i, err
}

const getContainerCapacity = `-- name: GetContainerCapacity :one
SELECT 
    c.capacity_weight,
    c.capacity_items
FROM 
    character_inventory ci
JOIN 
    containers c ON ci.item_id = c.base_item_id
WHERE 
    ci.id = ?
    AND ci.item_type = 'container'
`

type GetContainerCapacityRow struct {
	CapacityWeight float64       `json:"capacity_weight"`
	CapacityItems  sql.NullInt64 `json:"capacity_items"`
}

func (q *Queries) GetContainerCapacity(ctx context.Context, id int64) (GetContainerCapacityRow, error) {
	row := q.db.QueryRowContext(ctx, getContainerCapacity, id)
	var i GetContainerCapacityRow
	err := row.Scan(&i.CapacityWeight, &i.CapacityItems)
	return i, err
}

const getContainerContents = `-- name: GetContainerContents :many
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
    item_name
`

type GetContainerContentsParams struct {
	ContainerID sql.NullInt64 `json:"container_id"`
	CharacterID int64         `json:"character_id"`
}

type GetContainerContentsRow struct {
	ID              int64          `json:"id"`
	CharacterID     int64          `json:"character_id"`
	ItemID          int64          `json:"item_id"`
	ItemType        string         `json:"item_type"`
	Quantity        int64          `json:"quantity"`
	ContainerID     sql.NullInt64  `json:"container_id"`
	EquipmentSlotID sql.NullInt64  `json:"equipment_slot_id"`
	Notes           sql.NullString `json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	ItemName        interface{}    `json:"item_name"`
	ItemWeight      int64          `json:"item_weight"`
}

func (q *Queries) GetContainerContents(ctx context.Context, arg GetContainerContentsParams) ([]GetContainerContentsRow, error) {
	rows, err := q.db.QueryContext(ctx, getContainerContents, arg.ContainerID, arg.CharacterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetContainerContentsRow
	for rows.Next() {
		var i GetContainerContentsRow
		if err := rows.Scan(
			&i.ID,
			&i.CharacterID,
			&i.ItemID,
			&i.ItemType,
			&i.Quantity,
			&i.ContainerID,
			&i.EquipmentSlotID,
			&i.Notes,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ItemName,
			&i.ItemWeight,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getContainerWeight = `-- name: GetContainerWeight :one
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
    ci.container_id = ?
`

func (q *Queries) GetContainerWeight(ctx context.Context, containerID sql.NullInt64) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getContainerWeight, containerID)
	var total_weight interface{}
	err := row.Scan(&total_weight)
	return total_weight, err
}

const getEnhancedArmor = `-- name: GetEnhancedArmor :many
SELECT 
    id, name, weight, cost_gp
FROM 
    armor
WHERE
    enhancement_bonus = ?
ORDER BY 
    name
`

type GetEnhancedArmorRow struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Weight int64  `json:"weight"`
	CostGp int64  `json:"cost_gp"`
}

func (q *Queries) GetEnhancedArmor(ctx context.Context, enhancementBonus sql.NullInt64) ([]GetEnhancedArmorRow, error) {
	rows, err := q.db.QueryContext(ctx, getEnhancedArmor, enhancementBonus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEnhancedArmorRow
	for rows.Next() {
		var i GetEnhancedArmorRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Weight,
			&i.CostGp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEnhancedRangedWeapons = `-- name: GetEnhancedRangedWeapons :many
SELECT 
    id, name, weight, cost_gp
FROM 
    ranged_weapons
WHERE
    enhancement_bonus = ?
ORDER BY 
    name
`

type GetEnhancedRangedWeaponsRow struct {
	ID     int64         `json:"id"`
	Name   string        `json:"name"`
	Weight int64         `json:"weight"`
	CostGp sql.NullInt64 `json:"cost_gp"`
}

func (q *Queries) GetEnhancedRangedWeapons(ctx context.Context, enhancementBonus sql.NullInt64) ([]GetEnhancedRangedWeaponsRow, error) {
	rows, err := q.db.QueryContext(ctx, getEnhancedRangedWeapons, enhancementBonus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEnhancedRangedWeaponsRow
	for rows.Next() {
		var i GetEnhancedRangedWeaponsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Weight,
			&i.CostGp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEnhancedShields = `-- name: GetEnhancedShields :many
SELECT 
    id, name, weight, cost_gp
FROM 
    shields
WHERE
    enhancement_bonus = ?
ORDER BY 
    name
`

type GetEnhancedShieldsRow struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Weight int64  `json:"weight"`
	CostGp int64  `json:"cost_gp"`
}

func (q *Queries) GetEnhancedShields(ctx context.Context, enhancementBonus sql.NullInt64) ([]GetEnhancedShieldsRow, error) {
	rows, err := q.db.QueryContext(ctx, getEnhancedShields, enhancementBonus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEnhancedShieldsRow
	for rows.Next() {
		var i GetEnhancedShieldsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Weight,
			&i.CostGp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEnhancedWeapons = `-- name: GetEnhancedWeapons :many
SELECT 
    id, name, weight, cost_gp
FROM 
    weapons
WHERE
    enhancement_bonus = ?
ORDER BY 
    name
`

type GetEnhancedWeaponsRow struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Weight int64  `json:"weight"`
	CostGp int64  `json:"cost_gp"`
}

func (q *Queries) GetEnhancedWeapons(ctx context.Context, enhancementBonus sql.NullInt64) ([]GetEnhancedWeaponsRow, error) {
	rows, err := q.db.QueryContext(ctx, getEnhancedWeapons, enhancementBonus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEnhancedWeaponsRow
	for rows.Next() {
		var i GetEnhancedWeaponsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Weight,
			&i.CostGp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEquipmentSlots = `-- name: GetEquipmentSlots :many
SELECT 
    id, name, description
FROM 
    equipment_slots
ORDER BY 
    name
`

func (q *Queries) GetEquipmentSlots(ctx context.Context) ([]EquipmentSlot, error) {
	rows, err := q.db.QueryContext(ctx, getEquipmentSlots)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EquipmentSlot
	for rows.Next() {
		var i EquipmentSlot
		if err := rows.Scan(&i.ID, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEquippedItems = `-- name: GetEquippedItems :many
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
    es.name
`

type GetEquippedItemsRow struct {
	ID              int64          `json:"id"`
	CharacterID     int64          `json:"character_id"`
	ItemID          int64          `json:"item_id"`
	ItemType        string         `json:"item_type"`
	Quantity        int64          `json:"quantity"`
	ContainerID     sql.NullInt64  `json:"container_id"`
	EquipmentSlotID sql.NullInt64  `json:"equipment_slot_id"`
	Notes           sql.NullString `json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	SlotName        string         `json:"slot_name"`
	ItemName        interface{}    `json:"item_name"`
	ItemWeight      int64          `json:"item_weight"`
}

func (q *Queries) GetEquippedItems(ctx context.Context, characterID int64) ([]GetEquippedItemsRow, error) {
	rows, err := q.db.QueryContext(ctx, getEquippedItems, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEquippedItemsRow
	for rows.Next() {
		var i GetEquippedItemsRow
		if err := rows.Scan(
			&i.ID,
			&i.CharacterID,
			&i.ItemID,
			&i.ItemType,
			&i.Quantity,
			&i.ContainerID,
			&i.EquipmentSlotID,
			&i.Notes,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SlotName,
			&i.ItemName,
			&i.ItemWeight,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMagicalItemByID = `-- name: GetMagicalItemByID :one
SELECT 
    id, name, description, weight, cost_gp, max_charges, category, effect_description
FROM 
    magical_items
WHERE 
    id = ?
`

type GetMagicalItemByIDRow struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Weight            int64  `json:"weight"`
	CostGp            int64  `json:"cost_gp"`
	MaxCharges        int64  `json:"max_charges"`
	Category          string `json:"category"`
	EffectDescription string `json:"effect_description"`
}

func (q *Queries) GetMagicalItemByID(ctx context.Context, id int64) (GetMagicalItemByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getMagicalItemByID, id)
	var i GetMagicalItemByIDRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Weight,
		&i.CostGp,
		&i.MaxCharges,
		&i.Category,
		&i.EffectDescription,
	)
	return i, err
}

const isSlotOccupied = `-- name: IsSlotOccupied :one
SELECT 
    COUNT(*) > 0 as is_occupied
FROM 
    character_inventory
WHERE 
    character_id = ? 
    AND equipment_slot_id = ?
`

type IsSlotOccupiedParams struct {
	CharacterID     int64         `json:"character_id"`
	EquipmentSlotID sql.NullInt64 `json:"equipment_slot_id"`
}

func (q *Queries) IsSlotOccupied(ctx context.Context, arg IsSlotOccupiedParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, isSlotOccupied, arg.CharacterID, arg.EquipmentSlotID)
	var is_occupied bool
	err := row.Scan(&is_occupied)
	return is_occupied, err
}

const moveItemToContainer = `-- name: MoveItemToContainer :exec
UPDATE character_inventory
SET 
    container_id = ?,
    equipment_slot_id = NULL
WHERE 
    id = ?
    AND character_id = ?
`

type MoveItemToContainerParams struct {
	ContainerID sql.NullInt64 `json:"container_id"`
	ID          int64         `json:"id"`
	CharacterID int64         `json:"character_id"`
}

func (q *Queries) MoveItemToContainer(ctx context.Context, arg MoveItemToContainerParams) error {
	_, err := q.db.ExecContext(ctx, moveItemToContainer, arg.ContainerID, arg.ID, arg.CharacterID)
	return err
}

const reduceStackQuantity = `-- name: ReduceStackQuantity :exec
UPDATE character_inventory
SET quantity = quantity - ?
WHERE id = ? AND character_id = ?
`

type ReduceStackQuantityParams struct {
	Quantity    int64 `json:"quantity"`
	ID          int64 `json:"id"`
	CharacterID int64 `json:"character_id"`
}

func (q *Queries) ReduceStackQuantity(ctx context.Context, arg ReduceStackQuantityParams) error {
	_, err := q.db.ExecContext(ctx, reduceStackQuantity, arg.Quantity, arg.ID, arg.CharacterID)
	return err
}

const removeItemFromInventory = `-- name: RemoveItemFromInventory :exec
DELETE FROM character_inventory
WHERE
    id = ?
    AND character_id = ?
`

type RemoveItemFromInventoryParams struct {
	ID          int64 `json:"id"`
	CharacterID int64 `json:"character_id"`
}

func (q *Queries) RemoveItemFromInventory(ctx context.Context, arg RemoveItemFromInventoryParams) error {
	_, err := q.db.ExecContext(ctx, removeItemFromInventory, arg.ID, arg.CharacterID)
	return err
}

const splitStack = `-- name: SplitStack :exec
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
    ci.id = ?
`

type SplitStackParams struct {
	Quantity int64 `json:"quantity"`
	ID       int64 `json:"id"`
}

func (q *Queries) SplitStack(ctx context.Context, arg SplitStackParams) error {
	_, err := q.db.ExecContext(ctx, splitStack, arg.Quantity, arg.ID)
	return err
}

const unequipItem = `-- name: UnequipItem :exec
UPDATE character_inventory
SET 
    equipment_slot_id = NULL
WHERE 
    id = ?
    AND character_id = ?
`

type UnequipItemParams struct {
	ID          int64 `json:"id"`
	CharacterID int64 `json:"character_id"`
}

func (q *Queries) UnequipItem(ctx context.Context, arg UnequipItemParams) error {
	_, err := q.db.ExecContext(ctx, unequipItem, arg.ID, arg.CharacterID)
	return err
}

const updateInventoryItem = `-- name: UpdateInventoryItem :one
UPDATE character_inventory
SET
    quantity = ?,
    container_id = ?,
    equipment_slot_id = ?,
    notes = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = ?
    AND character_id = ? RETURNING id, character_id, item_id, item_type, quantity, container_id, equipment_slot_id, notes, created_at, updated_at
`

type UpdateInventoryItemParams struct {
	Quantity        int64          `json:"quantity"`
	ContainerID     sql.NullInt64  `json:"container_id"`
	EquipmentSlotID sql.NullInt64  `json:"equipment_slot_id"`
	Notes           sql.NullString `json:"notes"`
	ID              int64          `json:"id"`
	CharacterID     int64          `json:"character_id"`
}

type UpdateInventoryItemRow struct {
	ID              int64          `json:"id"`
	CharacterID     int64          `json:"character_id"`
	ItemID          int64          `json:"item_id"`
	ItemType        string         `json:"item_type"`
	Quantity        int64          `json:"quantity"`
	ContainerID     sql.NullInt64  `json:"container_id"`
	EquipmentSlotID sql.NullInt64  `json:"equipment_slot_id"`
	Notes           sql.NullString `json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (q *Queries) UpdateInventoryItem(ctx context.Context, arg UpdateInventoryItemParams) (UpdateInventoryItemRow, error) {
	row := q.db.QueryRowContext(ctx, updateInventoryItem,
		arg.Quantity,
		arg.ContainerID,
		arg.EquipmentSlotID,
		arg.Notes,
		arg.ID,
		arg.CharacterID,
	)
	var i UpdateInventoryItemRow
	err := row.Scan(
		&i.ID,
		&i.CharacterID,
		&i.ItemID,
		&i.ItemType,
		&i.Quantity,
		&i.ContainerID,
		&i.EquipmentSlotID,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateItemQuantity = `-- name: UpdateItemQuantity :exec
UPDATE character_inventory
SET 
    quantity = ?
WHERE 
    id = ?
    AND character_id = ?
`

type UpdateItemQuantityParams struct {
	Quantity    int64 `json:"quantity"`
	ID          int64 `json:"id"`
	CharacterID int64 `json:"character_id"`
}

func (q *Queries) UpdateItemQuantity(ctx context.Context, arg UpdateItemQuantityParams) error {
	_, err := q.db.ExecContext(ctx, updateItemQuantity, arg.Quantity, arg.ID, arg.CharacterID)
	return err
}

const useChargedItem = `-- name: UseChargedItem :exec
UPDATE character_inventory
SET 
    charges = charges - 1
WHERE 
    id = ?
    AND character_id = ?
    AND charges > 0
`

type UseChargedItemParams struct {
	ID          int64 `json:"id"`
	CharacterID int64 `json:"character_id"`
}

func (q *Queries) UseChargedItem(ctx context.Context, arg UseChargedItemParams) error {
	_, err := q.db.ExecContext(ctx, useChargedItem, arg.ID, arg.CharacterID)
	return err
}
