CREATE TABLE places (
    id         TEXT          PRIMARY KEY,
    state_id   TEXT          NOT NULL REFERENCES states (id) ON DELETE CASCADE,
    name       TEXT          NOT NULL,
    lat        NUMERIC(10,7),
    lon        NUMERIC(10,7),
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_places_state_id ON places (state_id);
CREATE INDEX ix_places_name     ON places (name);
