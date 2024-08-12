-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL, 
    license_status TEXT NOT NULL DEFAULT 'active',
    license_lvl TEXT DEFAULT 'basic',
    license_hash TEXT UNIQUE NOT NULL,
    license_expires_date TIMESTAMP,
    wazzup_uri TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
