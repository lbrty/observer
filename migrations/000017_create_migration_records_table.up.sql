-- Movement events are immutable historical facts.
-- destination_place_id: where the person moved to in this event.
-- from_place_id:        where they were before this move.
-- project_id is omitted — derivable via person_id → people.project_id.
-- updated_at is omitted — movement records are never edited; corrections are new records.
CREATE TABLE migration_records (
    id                   TEXT        PRIMARY KEY,
    person_id            TEXT        NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    destination_place_id TEXT        REFERENCES places (id) ON DELETE SET NULL,
    from_place_id        TEXT        REFERENCES places (id) ON DELETE SET NULL,
    migration_date       DATE,
    notes                TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_migration_records_person_id            ON migration_records (person_id);
CREATE INDEX ix_migration_records_destination_place_id ON migration_records (destination_place_id);
