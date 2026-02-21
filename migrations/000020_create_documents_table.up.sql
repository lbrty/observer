-- Documents attached to a person record, scoped to a project.
-- path:               relative path from the storage root (not a full URL or absolute path).
-- encryption_key_ref: KMS key identifier used to decrypt this file; NULL for unencrypted files.
--                     The raw key is never stored here.
CREATE TABLE documents (
    id                TEXT        PRIMARY KEY,
    person_id         TEXT        NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    project_id        TEXT        NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    uploaded_by       TEXT        REFERENCES users (id) ON DELETE SET NULL,
    encryption_key_ref TEXT,
    name              TEXT        NOT NULL,
    path              TEXT        NOT NULL,
    mime_type         TEXT        NOT NULL,
    size              BIGINT      NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_documents_person_id  ON documents (person_id);
CREATE INDEX ix_documents_project_id ON documents (project_id);
