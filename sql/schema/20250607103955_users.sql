-- +goose Up
-- +goose StatementBegin
SELECT
    'up SQL query';

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';

DROP TABLE users;

-- +goose StatementEnd
