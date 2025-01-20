-- name: CreateCharacter :one
INSERT INTO
    characters (
        user_id,
        name,
        max_hp,
        current_hp,
        strength,
        dexterity,
        constitution,
        intelligence,
        wisdom,
        charisma
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: GetCharacter :one
SELECT
    *
FROM
    characters
WHERE
    id = ?
    AND user_id = ?
LIMIT
    1;

-- name: ListCharactersByUser :many
SELECT
    *
FROM
    characters
WHERE
    user_id = ?
ORDER BY
    name;

-- name: UpdateCharacter :one
UPDATE characters
SET
    name = ?,
    max_hp = ?,
    current_hp = ?,
    strength = ?,
    dexterity = ?,
    constitution = ?,
    intelligence = ?,
    wisdom = ?,
    charisma = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = ?
    AND user_id = ? RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE
    id = ?
    AND user_id = ?;
