-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS admins
(
    id INTEGER PRIMARY KEY
);
CREATE TABLE IF NOT EXISTS apps 
(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users(email);
INSERT INTO apps (id, name, secret) VALUES (1, 'test', 'secret');
INSERT INTO admins (id) VALUES (1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS apps;
-- +goose StatementEnd