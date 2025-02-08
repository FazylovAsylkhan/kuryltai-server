-- +goose Up
ALTER TABLE users ADD COLUMN api_key VARCHAR(64) UNIQUE;
UPDATE users SET api_key = encode(sha256(random()::text::bytea), 'hex');
ALTER TABLE users ALTER COLUMN api_key SET NOT NULL;

-- +goose Down
ALTER TABLE users ALTER COLUMN api_key DROP NOT NULL;
ALTER TABLE users DROP COLUMN api_key;
