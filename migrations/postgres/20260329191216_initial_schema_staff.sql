-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS staff.employees (
    id                   BIGSERIAL PRIMARY KEY,
    full_name            TEXT                     NOT NULL,
    position             INTEGER                  NOT NULL,
    created_by_id        BIGINT                   NOT NULL,
    created_by_name      TEXT                     NOT NULL,
    created_at           TIMESTAMPTZ              NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS staff.offices (
    id                   BIGSERIAL PRIMARY KEY,
    name                 TEXT                     NOT NULL,
    address              TEXT                     NOT NULL,
    created_by_id        BIGINT                   NOT NULL,
    created_by_name      TEXT                     NOT NULL,
    created_at           TIMESTAMPTZ              NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS staff.offices;
DROP TABLE IF EXISTS staff.employees;

-- +goose StatementEnd
