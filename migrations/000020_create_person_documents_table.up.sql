CREATE TABLE person_documents (
    id             TEXT        PRIMARY KEY,
    person_id      TEXT        NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    project_id     TEXT        NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    uploaded_by    TEXT        REFERENCES users (id) ON DELETE SET NULL,
    encryption_key TEXT,
    name           TEXT        NOT NULL,
    path           TEXT        NOT NULL,
    mime_type      TEXT        NOT NULL,
    size           BIGINT      NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_person_documents_person_id  ON person_documents (person_id);
CREATE INDEX ix_person_documents_project_id ON person_documents (project_id);
