ALTER TABLE users ADD COLUMN office_id TEXT REFERENCES offices (id) ON DELETE SET NULL;

CREATE INDEX ix_users_office_id ON users (office_id);
