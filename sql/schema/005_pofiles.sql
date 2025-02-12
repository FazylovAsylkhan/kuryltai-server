-- +goose Up
CREATE TABLE profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slug VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    avatar_image VARCHAR(255) DEFAULT '',
    cover_image VARCHAR(255) DEFAULT '',
    head_line VARCHAR(255) DEFAULT ''
);

-- +goose Down
DROP TABLE IF EXISTS profiles;