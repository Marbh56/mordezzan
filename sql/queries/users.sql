-- name: CreateUser :one
INSERT INTO
    users (username, email, password_hash)
VALUES
    (?, ?, ?) RETURNING *;

-- name: GetUserByUsername :one
SELECT
    *
FROM
    users
WHERE
    username = ?
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

-- name: CreateSession :one
INSERT INTO
    sessions (token, user_id, expires_at)
VALUES
    (?, ?, ?) RETURNING *;

-- name: GetSession :one
SELECT
    s.*,
    u.username,
    u.email
FROM
    sessions s
    JOIN users u ON s.user_id = u.id
WHERE
    s.token = ?
    AND s.expires_at > CURRENT_TIMESTAMP
LIMIT
    1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE
    token = ?;
