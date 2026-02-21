CREATE TYPE support_type AS ENUM (
    'humanitarian',
    'legal',
    'medical',
    'general'
);

CREATE TABLE support_records (
    id            TEXT         PRIMARY KEY,
    person_id     TEXT         NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    project_id    TEXT         NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    consultant_id TEXT         REFERENCES users (id) ON DELETE SET NULL,
    owner_id      TEXT         REFERENCES users (id) ON DELETE SET NULL,
    type          support_type NOT NULL DEFAULT 'general',
    notes         TEXT,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_support_records_person_id     ON support_records (person_id);
CREATE INDEX ix_support_records_project_id    ON support_records (project_id);
CREATE INDEX ix_support_records_consultant_id ON support_records (consultant_id);
CREATE INDEX ix_support_records_type          ON support_records (type);
