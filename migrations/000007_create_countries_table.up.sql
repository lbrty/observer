CREATE TABLE countries (
    id         TEXT        PRIMARY KEY,
    name       TEXT        NOT NULL,
    code       CITEXT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_countries_code ON countries (code);
CREATE INDEX        ix_countries_name ON countries (name);
