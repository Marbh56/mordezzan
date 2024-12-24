-- name: ViewCharacter :one
SELECT *
FROM characters
WHERE id = $1;
