-- name: CreateUser :one
INSERT INTO
    users (name, email, password_hash)
VALUES
    (?, ?, ?) RETURNING *;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = ?
LIMIT
    1;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = ?
LIMIT
    1;

-- name: ListUsers :many
SELECT
    *
FROM
    users
ORDER BY
    name;

-- name: UpdateUser :one
UPDATE users
SET
    name = ?,
    email = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
    id = ?;
