CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL,
    user_agent TEXT,
    ip VARCHAR(45),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_sessions_refresh_token UNIQUE (refresh_token)
);

CREATE INDEX ix_sessions_user_id ON sessions(user_id);
CREATE INDEX ix_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX ix_sessions_expires_at ON sessions(expires_at);
