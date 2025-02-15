-- name: ListEquipmentSlots :many
SELECT
    id,
    name,
    description
FROM
    equipment_slots
ORDER BY
    name;


-- name: GetEquipmentItems :many
SELECT id, name, cost_gp as cost, weight
FROM equipment
ORDER BY name;

-- name: GetWeaponItems :many
SELECT id, name, cost_gp as cost, weight
FROM weapons
ORDER BY name;

-- name: GetArmorItems :many
SELECT id, name, cost_gp as cost, weight
FROM armor
ORDER BY name;

-- name: GetAmmunitionItems :many
SELECT id, name, cost_gp as cost, weight
FROM ammunition
ORDER BY name;

-- name: GetContainerItems :many
SELECT id, name, cost_gp as cost, weight
FROM containers
ORDER BY name;

-- name: GetShieldItems :many
SELECT id, name, cost_gp as cost, weight
FROM shields
ORDER BY name;

-- name: GetRangedWeaponItems :many
SELECT id, name, cost_gp as cost, weight
FROM ranged_weapons
ORDER BY name;