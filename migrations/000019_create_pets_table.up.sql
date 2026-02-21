CREATE TYPE pet_status AS ENUM (
    'registered',
    'adopted',
    'owner_found',
    'needs_shelter',
    'unknown'
);

CREATE TABLE pets (
    id              TEXT        PRIMARY KEY,
    project_id      TEXT        NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    owner_id        TEXT        REFERENCES people (id) ON DELETE SET NULL,
    name            TEXT        NOT NULL,
    status          pet_status  NOT NULL DEFAULT 'unknown',
    registration_id TEXT,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_pets_project_id ON pets (project_id);
CREATE INDEX ix_pets_owner_id   ON pets (owner_id);
CREATE INDEX ix_pets_status     ON pets (status);
