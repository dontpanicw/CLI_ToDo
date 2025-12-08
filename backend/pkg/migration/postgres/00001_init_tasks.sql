-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT ''
);

-- +goose Down
DROP TABLE IF EXISTS tasks;

