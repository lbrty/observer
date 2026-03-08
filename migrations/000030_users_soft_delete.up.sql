ALTER TABLE users ADD COLUMN deactivated_at TIMESTAMPTZ;
CREATE INDEX ix_users_deactivated ON users (deactivated_at) WHERE deactivated_at IS NOT NULL;
