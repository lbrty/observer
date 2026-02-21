CREATE TYPE project_role AS ENUM ('owner', 'manager', 'consultant', 'viewer');

CREATE TABLE project_permissions (
    id                  TEXT         PRIMARY KEY,
    project_id          TEXT         NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    user_id             TEXT         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role                project_role NOT NULL DEFAULT 'viewer',
    can_view_contact    BOOLEAN      NOT NULL DEFAULT FALSE,
    can_view_personal   BOOLEAN      NOT NULL DEFAULT FALSE,
    can_view_documents  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_project_permissions_user_project ON project_permissions (user_id, project_id);
CREATE INDEX        ix_project_permissions_project_id   ON project_permissions (project_id);
CREATE INDEX        ix_project_permissions_user_id      ON project_permissions (user_id);
