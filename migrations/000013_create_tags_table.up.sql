-- Tags are project-scoped: each project maintains its own tag vocabulary.
-- A tag name is unique within a project but not globally.
CREATE TABLE tags (
    id         TEXT        PRIMARY KEY,
    project_id TEXT        NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_tags_project_name ON tags (project_id, name);
CREATE INDEX        ix_tags_project_id   ON tags (project_id);
