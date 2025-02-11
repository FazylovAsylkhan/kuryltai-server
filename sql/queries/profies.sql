-- name: CreateProfile :one
INSERT INTO profiles (id, user_id, username, slug, avatar_image, cover_image, head_line, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetProfile :one
SELECT * FROM profiles
WHERE user_id = $1;

-- name: GetProfileBySlug :one
SELECT * FROM profiles
WHERE slug = $1;

-- name: UpdateProfile :one
UPDATE profiles 
SET slug = $2, username = $3, avatar_image = $4, cover_image = $5
WHERE user_id = $1
RETURNING *;