-- name: GetSpellByID :one
SELECT * FROM spells WHERE id = ?;

-- name: GetSpellLevels :many
SELECT class, level
FROM spell_levels
WHERE spell_id = ?
ORDER BY class;

-- name: ListSpellsByClass :many
SELECT 
    s.id, s.name, sl.level
FROM 
    spells s
JOIN 
    spell_levels sl ON s.id = sl.spell_id
WHERE 
    sl.class = ?
ORDER BY 
    sl.level, s.name;

-- name: ListSpellsByClassAndLevel :many
SELECT 
    s.id, s.name
FROM 
    spells s
JOIN 
    spell_levels sl ON s.id = sl.spell_id
WHERE 
    sl.class = ? AND sl.level = ?
ORDER BY 
    s.name;

-- name: AddSpell :one
INSERT INTO spells (
    name, description, range, duration
) VALUES (
    ?, ?, ?, ?
) RETURNING *;

-- name: AddSpellLevel :exec
INSERT INTO spell_levels (
    spell_id, class, level
) VALUES (
    ?, ?, ?
);

-- name: UpdateSpell :one
UPDATE spells
SET 
    name = ?,
    description = ?,
    range = ?,
    duration = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = ?
RETURNING *;

-- name: DeleteSpell :exec
DELETE FROM spells
WHERE id = ?;

-- name: DeleteSpellLevels :exec
DELETE FROM spell_levels
WHERE spell_id = ?;

-- name: ListAllSpells :many
SELECT id, name FROM spells ORDER BY name;