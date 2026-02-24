-- Informal case-worker notes attached to a person record.
-- Distinct from support_records (formal, quantified, reportable service events).
-- Notes are internal observations, follow-up reminders, and context.
CREATE TABLE person_notes (
    id         TEXT        PRIMARY KEY,
    person_id  TEXT        NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    author_id  TEXT        REFERENCES users (id) ON DELETE SET NULL,
    body       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_person_notes_person_id ON person_notes (person_id);
CREATE INDEX ix_person_notes_author_id ON person_notes (author_id);
