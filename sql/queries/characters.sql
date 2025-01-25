-- name: CreateCharacter :one
INSERT INTO
    characters (
        user_id,
        name,
        class,
        level,
        max_hp,
        current_hp,
        strength,
        dexterity,
        constitution,
        intelligence,
        wisdom,
        charisma,
        experience_points
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

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
    class = ?,
    level = ?,
    max_hp = ?,
    current_hp = ?,
    strength = ?,
    dexterity = ?,
    constitution = ?,
    intelligence = ?,
    wisdom = ?,
    charisma = ?,
    experience_points = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = ?
    AND user_id = ? RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE
    id = ?
    AND user_id = ?;
