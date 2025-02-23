-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES (?, ?, ?) RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = ? 
AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?
AND deleted_at IS NULL
LIMIT 1;

-- name: GetSession :one
SELECT s.*, u.username, u.email
FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.token = ?
AND s.expires_at > CURRENT_TIMESTAMP
AND u.deleted_at IS NULL
LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token = ?;

-- name: CreateSession :one
INSERT INTO sessions (token, user_id, expires_at)
VALUES (?, ?, ?) RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users 
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ? AND deleted_at IS NULL;