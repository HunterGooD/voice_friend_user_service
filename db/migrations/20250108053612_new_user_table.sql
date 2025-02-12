-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, 
    login VARCHAR(160) NOT NULL UNIQUE, 
    name VARCHAR(160) NOT NULL, 
    uid UUID DEFAULT gen_random_uuid(), 
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL, 
    is_active BOOLEAN DEFAULT TRUE, 
    last_login TIMESTAMP DEFAULT NULL, 
    role VARCHAR(20) DEFAULT 'user', 
    profile_picture VARCHAR(255), 
    phone VARCHAR(20),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    deleted_at TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_users_uid_hash ON users USING hash (uid);
CREATE INDEX IF NOT EXISTS idx_users_login_hash ON users USING hash (login);

CREATE INDEX IF NOT EXISTS idx_users_is_active ON users (is_active);
CREATE INDEX IF NOT EXISTS idx_users_is_delete ON users (deleted_at);
-- CREATE INDEX idx_users_last_login ON users (last_login DESC);
-- CREATE INDEX idx_users_phone ON users (phone);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
