-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdatePassword :one
UPDATE users 
SET password = $2
WHERE id = $1
RETURNING *;