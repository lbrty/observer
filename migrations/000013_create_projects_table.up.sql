CREATE TABLE projects (
    id          TEXT        PRIMARY KEY,
    name        TEXT        NOT NULL,
    description TEXT,
    owner_id    TEXT        NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_projects_owner_id ON projects (owner_id);
CREATE INDEX ix_projects_name     ON projects (name);
