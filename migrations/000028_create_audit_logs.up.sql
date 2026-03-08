CREATE TABLE audit_logs (
    id          TEXT        PRIMARY KEY,
    project_id  TEXT        REFERENCES projects(id) ON DELETE SET NULL,
    user_id     TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action      TEXT        NOT NULL,
    entity_type TEXT        NOT NULL,
    entity_id   TEXT,
    summary     TEXT        NOT NULL,
    ip          TEXT,
    user_agent  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_audit_logs_project_time ON audit_logs (project_id, created_at DESC);
CREATE INDEX ix_audit_logs_user_time    ON audit_logs (user_id, created_at DESC);
CREATE INDEX ix_audit_logs_action_time  ON audit_logs (action, created_at DESC);
