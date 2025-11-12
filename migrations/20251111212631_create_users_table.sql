-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id uuid DEFAULT uuidv7() PRIMARY KEY,
    email text NOT NULL UNIQUE,
    password_hash text DEFAULT null,
    firebase_uid text DEFAULT null,
    revenuecat_id text DEFAULT null,
    timezone text DEFAULT null,
    locale text DEFAULT null,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
