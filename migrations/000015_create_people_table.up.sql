CREATE TYPE person_case_status AS ENUM (
    'new',
    'active',
    'closed',
    'archived'
);

CREATE TYPE person_sex AS ENUM (
    'male',
    'female',
    'other',
    'unknown'
);

-- Age ranges (used when birth_date is not available):
--   infant            0–1
--   toddler           1–3
--   pre_school        3–6
--   middle_childhood  6–12
--   young_teen        12–14
--   teenager          14–18
--   young_adult       18–25
--   early_adult       25–35
--   middle_aged_adult 35–55
--   old_adult         55+
--   unknown           not specified
CREATE TYPE person_age_group AS ENUM (
    'infant',
    'toddler',
    'pre_school',
    'middle_childhood',
    'young_teen',
    'teenager',
    'young_adult',
    'early_adult',
    'middle_aged_adult',
    'old_adult',
    'unknown'
);

CREATE TABLE people (
    id               TEXT               PRIMARY KEY,
    project_id       TEXT               NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    parent_id        TEXT               REFERENCES people (id) ON DELETE SET NULL,
    category_id      TEXT               REFERENCES categories (id) ON DELETE SET NULL,
    consultant_id    TEXT               REFERENCES users (id) ON DELETE SET NULL,
    office_id        TEXT               REFERENCES offices (id) ON DELETE SET NULL,
    current_place_id TEXT               REFERENCES places (id) ON DELETE SET NULL,
    origin_place_id  TEXT               REFERENCES places (id) ON DELETE SET NULL,
    external_id      TEXT,
    first_name       TEXT               NOT NULL,
    last_name        TEXT,
    patronymic       TEXT,
    email            CITEXT,
    birth_date       DATE,
    sex              person_sex         NOT NULL DEFAULT 'unknown',
    age_group        person_age_group,
    phone_numbers    JSONB              NOT NULL DEFAULT '[]',
    idp_status       TEXT               CHECK (idp_status IN ('crimea', 'east_ukraine', 'non_idp')),
    case_status      person_case_status NOT NULL DEFAULT 'new',
    registered_at    DATE,
    created_at       TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    -- birth_date and age_group are mutually exclusive input methods
    CONSTRAINT chk_people_age_xor CHECK (birth_date IS NULL OR age_group IS NULL)
);

CREATE INDEX ix_people_project_id          ON people (project_id);
CREATE INDEX ix_people_parent_id           ON people (parent_id);
CREATE INDEX ix_people_consultant_id       ON people (consultant_id);
CREATE INDEX ix_people_category_id         ON people (category_id);
CREATE INDEX ix_people_current_place_id    ON people (current_place_id);
CREATE INDEX ix_people_case_status         ON people (case_status);
CREATE INDEX ix_people_project_case_status ON people (project_id, case_status);
CREATE INDEX ix_people_project_consultant  ON people (project_id, consultant_id);
CREATE INDEX ix_people_idp_status          ON people (idp_status) WHERE idp_status IS NOT NULL;
CREATE INDEX ix_people_registered_at       ON people (registered_at) WHERE registered_at IS NOT NULL;
CREATE INDEX ix_people_first_name          ON people USING gin (first_name gin_trgm_ops);
CREATE INDEX ix_people_last_name           ON people USING gin (last_name  gin_trgm_ops) WHERE last_name IS NOT NULL;
CREATE INDEX ix_people_email               ON people (email) WHERE email IS NOT NULL;
