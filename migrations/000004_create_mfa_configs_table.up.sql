CREATE TABLE mfa_configs (
    user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    method VARCHAR(10) NOT NULL CHECK (method IN ('totp', 'sms')),
    secret TEXT,
    phone VARCHAR(20),
    is_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);
