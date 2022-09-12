-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS events (
    id uuid NOT NULL,
    title VARCHAR(128) NOT NULL,
    datetime TIMESTAMPTZ NOT NULL,
    duration BIGINT NOT NULL,
    description TEXT NOT NULL,
    owner_id uuid NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ,
    CONSTRAINT events_pkey PRIMARY KEY (id),
    CONSTRAINT events_datetime_owner_id_key UNIQUE (datetime, owner_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
