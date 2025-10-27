-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS url_mappings (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    short_code varchar(10) NOT NULL UNIQUE,
    long_url text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    click_count bigint NOT NULL DEFAULT 0,
    last_accessed_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_url_mappings_short_code ON url_mappings(short_code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url_mappings;
-- +goose StatementEnd
