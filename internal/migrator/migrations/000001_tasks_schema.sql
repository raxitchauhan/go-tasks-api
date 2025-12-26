-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS tasks;

CREATE TABLE tasks.tasks (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'todo',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS tasks CASCADE;

-- +goose StatementEnd