CREATE TABLE verification_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('verification', 'password_reset')),
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_verification_tokens_token UNIQUE (token)
);

CREATE INDEX ix_verification_tokens_token ON verification_tokens(token);
CREATE INDEX ix_verification_tokens_user_id ON verification_tokens(user_id);
