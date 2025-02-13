-- name: CreateSession :one
INSERT INTO sessions (id, user_id, user_email, refresh_token, expires_at, is_revoked, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1;

-- name: RevokeSession :one
UPDATE sessions
SET is_revoked = true
WHERE id = $1
RETURNING *;

-- name: DeleteSession :one
DELETE FROM sessions
WHERE user_id = $1
RETURNING *;