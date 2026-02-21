-- support_type — the broad category of service delivered:
--   humanitarian  — material aid (food, clothing, hygiene kits)
--   legal         — legal consultations and document assistance
--   social        — social services, benefits navigation
--   psychological — psychological support and counselling
--   medical       — medical referrals and assistance
--   general       — unclassified or multi-category support
CREATE TYPE support_type AS ENUM (
    'humanitarian',
    'legal',
    'social',
    'psychological',
    'medical',
    'general'
);

-- support_sphere — the specific topic area within a consultation.
-- Used for "by sphere of appeal" report breakdowns.
-- Values follow snake_case convention.
CREATE TYPE support_sphere AS ENUM (
    'housing_assistance',   -- housing rights, eviction, social housing, temporary shelter
    'document_recovery',    -- passports, birth certificates, property documents
    'social_benefits',      -- IDP registration, social payments, benefit eligibility
    'property_rights',      -- property left in occupied territories, restitution
    'employment_rights',    -- labour law, dismissal, job placement
    'family_law',           -- divorce, custody, alimony, guardianship
    'healthcare_access',    -- medical coverage, disability documentation
    'education_access',     -- school enrolment, tutoring, educational rights
    'financial_aid',        -- emergency financial assistance, humanitarian grants
    'psychological_support',-- mental health referrals, counselling coordination
    'other'                 -- unlisted or cross-cutting topics
);

CREATE TABLE support_records (
    id            TEXT           PRIMARY KEY,
    person_id     TEXT           NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    project_id    TEXT           NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    consultant_id TEXT           REFERENCES users (id) ON DELETE SET NULL,
    recorded_by   TEXT           REFERENCES users (id) ON DELETE SET NULL,
    office_id     TEXT           REFERENCES offices (id) ON DELETE SET NULL,
    type          support_type   NOT NULL DEFAULT 'general',
    sphere        support_sphere,
    provided_at   DATE,
    notes         TEXT,
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_support_records_person_id     ON support_records (person_id);
CREATE INDEX ix_support_records_project_id    ON support_records (project_id);
CREATE INDEX ix_support_records_consultant_id ON support_records (consultant_id);
CREATE INDEX ix_support_records_office_id     ON support_records (office_id);
CREATE INDEX ix_support_records_type          ON support_records (type);
CREATE INDEX ix_support_records_provided_at   ON support_records (provided_at);
CREATE INDEX ix_support_records_sphere        ON support_records (sphere) WHERE sphere IS NOT NULL;
