-- conflict_zone tags an oblast/state with its conflict zone designation.
-- Used to derive IDP classification from origin_place_id without hardcoding
-- geopolitical categories on the person record.
CREATE TABLE states (
    id            TEXT        PRIMARY KEY,
    country_id    TEXT        NOT NULL REFERENCES countries (id) ON DELETE CASCADE,
    name          TEXT        NOT NULL,
    code          CITEXT,
    conflict_zone TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX        ix_states_country_id        ON states (country_id);
CREATE INDEX        ix_states_name              ON states (name);
CREATE UNIQUE INDEX uq_states_country_code      ON states (country_id, code) WHERE code IS NOT NULL;
