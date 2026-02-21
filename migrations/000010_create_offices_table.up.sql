CREATE TABLE offices (
    id         TEXT        PRIMARY KEY,
    name       TEXT        NOT NULL,
    place_id   TEXT        REFERENCES places (id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_offices_name     ON offices (name);
CREATE INDEX        ix_offices_place_id ON offices (place_id);
