CREATE TABLE credentials (
    user_id       TEXT        PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    password_hash TEXT        NOT NULL,
    salt          TEXT        NOT NULL,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
