-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS auth.users (
    id             BIGSERIAL PRIMARY KEY,
    employee_id    BIGINT                   NOT NULL,
    password_hash  TEXT                     NOT NULL,
    created_at     TIMESTAMPTZ              NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS auth.users;

-- +goose StatementEnd
