-- name: ListEquipmentSlots :many
SELECT
    id,
    name,
    description
FROM
    equipment_slots
ORDER BY
    name;
