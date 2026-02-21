CREATE TYPE person_status AS ENUM (
    'active',
    'inactive',
    'registered',
    'needs_help',
    'consulted',
    'helped',
    'unknown'
);

CREATE TYPE person_sex AS ENUM (
    'male',
    'female',
    'other',
    'unknown'
);

CREATE TYPE person_age_group AS ENUM (
    'infant',
    'toddler',
    'preschool',
    'early_school',
    'preteen',
    'teen',
    'older_teen',
    'young_adult',
    'adult',
    'senior',
    'elderly'
);

CREATE TABLE people (
    id               TEXT             PRIMARY KEY,
    project_id       TEXT             NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    parent_id        TEXT             REFERENCES people (id) ON DELETE SET NULL,
    category_id      TEXT             REFERENCES categories (id) ON DELETE SET NULL,
    consultant_id    TEXT             REFERENCES users (id) ON DELETE SET NULL,
    office_id        TEXT             REFERENCES offices (id) ON DELETE SET NULL,
    current_place_id TEXT             REFERENCES places (id) ON DELETE SET NULL,
    origin_place_id  TEXT             REFERENCES places (id) ON DELETE SET NULL,
    external_id      TEXT,
    full_name        TEXT             NOT NULL,
    email            CITEXT,
    birth_date       DATE,
    sex              person_sex       NOT NULL DEFAULT 'unknown',
    age_group        person_age_group,
    phone_numbers    JSONB            NOT NULL DEFAULT '[]',
    status           person_status    NOT NULL DEFAULT 'unknown',
    created_at       TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_people_project_id        ON people (project_id);
CREATE INDEX ix_people_parent_id         ON people (parent_id);
CREATE INDEX ix_people_consultant_id     ON people (consultant_id);
CREATE INDEX ix_people_category_id       ON people (category_id);
CREATE INDEX ix_people_current_place_id  ON people (current_place_id);
CREATE INDEX ix_people_status            ON people (status);
CREATE INDEX ix_people_project_status    ON people (project_id, status);
CREATE INDEX ix_people_project_consultant ON people (project_id, consultant_id);
CREATE INDEX ix_people_full_name         ON people USING gin (full_name gin_trgm_ops);
CREATE INDEX ix_people_email             ON people (email) WHERE email IS NOT NULL;
