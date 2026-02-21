CREATE TABLE credentials (
    user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    salt TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);
