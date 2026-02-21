CREATE TABLE migration_records (
    id               TEXT        PRIMARY KEY,
    person_id        TEXT        NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    project_id       TEXT        NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    current_place_id TEXT        REFERENCES places (id) ON DELETE SET NULL,
    from_place_id    TEXT        REFERENCES places (id) ON DELETE SET NULL,
    migration_date   DATE,
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_migration_records_person_id        ON migration_records (person_id);
CREATE INDEX ix_migration_records_project_id       ON migration_records (project_id);
CREATE INDEX ix_migration_records_current_place_id ON migration_records (current_place_id);
