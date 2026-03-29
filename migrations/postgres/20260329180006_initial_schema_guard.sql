-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS guard.tasks (
    id              BIGSERIAL PRIMARY KEY,
    damage_amount   BIGINT                   NOT NULL,
    priority        INTEGER                  NOT NULL,
    status          INTEGER                  NOT NULL,
    start_date      TIMESTAMPTZ              NOT NULL DEFAULT NOW(),
    end_date        TIMESTAMPTZ,
    executor_id     BIGINT                   NOT NULL,
    executor_name   TEXT                     NOT NULL,
    created_by_id   BIGINT                   NOT NULL,
    created_by_name TEXT                     NOT NULL,
    created_at      TIMESTAMPTZ              NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS guard.vud_decisions (
    id                   BIGSERIAL PRIMARY KEY,
    task_id              BIGINT                   NOT NULL REFERENCES guard.tasks (id),
    criminal_case_opened BOOLEAN,
    comment              TEXT,
    kusp                 VARCHAR(10)              NOT NULL,
    ud                   VARCHAR(25),
    created_by_id        BIGINT                   NOT NULL,
    created_by_name      TEXT                     NOT NULL,
    created_at           TIMESTAMPTZ              NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS guard.refunds (
    id                   BIGSERIAL PRIMARY KEY,
    task_id              BIGINT                   NOT NULL REFERENCES guard.tasks (id),
    amount               INTEGER                  NOT NULL,
    comment              TEXT,
    created_by_id        BIGINT                   NOT NULL,
    created_by_name      TEXT                     NOT NULL,
    created_at           TIMESTAMPTZ              NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS guard.comments (
    id                   BIGSERIAL PRIMARY KEY,
    task_id              BIGINT                   NOT NULL REFERENCES guard.tasks (id),
    comment              TEXT,
    created_by_id        BIGINT                   NOT NULL,
    created_by_name      TEXT                     NOT NULL,
    created_at           TIMESTAMPTZ              NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS guard.violators (
    id                   BIGSERIAL PRIMARY KEY,
    task_id              BIGINT                   NOT NULL REFERENCES guard.tasks (id),
    type                 INTEGER                  NOT NULL,
    full_name            TEXT                     NOT NULL,
    phone                TEXT
);

CREATE TABLE IF NOT EXISTS guard.history_changes (
    id              BIGSERIAL PRIMARY KEY,
    task_id         BIGINT                   NOT NULL REFERENCES guard.tasks (id),
    event           INTEGER                  NOT NULL,
    metadata_json   JSONB,
    created_by_id   BIGINT                   NOT NULL,
    created_by_name TEXT                     NOT NULL,
    created_at      TIMESTAMPTZ              NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS guard.history_changes;
DROP TABLE IF EXISTS guard.violators;
DROP TABLE IF EXISTS guard.comments;
DROP TABLE IF EXISTS guard.refunds;
DROP TABLE IF EXISTS guard.vud_decisions;
DROP TABLE IF EXISTS guard.tasks;

-- +goose StatementEnd