-- A household is a first-class entity representing a family unit.
-- It replaces the people.parent_id self-reference, which could not represent
-- peer relationships (e.g. spouses) or carry relationship semantics.
--
-- reference_number: human-readable case file ID for paper documents
--   (e.g. KYV-2024-00142). Optional.
-- head_person_id: the nominated head of household for entitlement attribution.
--   Nullable — a household can exist before a head is designated.
CREATE TABLE households (
    id               TEXT        PRIMARY KEY,
    project_id       TEXT        NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    reference_number TEXT,
    head_person_id   TEXT        REFERENCES people (id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- household_members links people to their household with an explicit relationship type.
-- A person may belong to at most one household per the PRIMARY KEY constraint.
-- relationship values follow UNHCR PSN/proGres conventions where applicable.
CREATE TABLE household_members (
    household_id TEXT NOT NULL REFERENCES households (id) ON DELETE CASCADE,
    person_id    TEXT NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    relationship TEXT NOT NULL CHECK (relationship IN (
        'head', 'spouse', 'child', 'parent', 'sibling',
        'grandchild', 'grandparent', 'other_relative', 'non_relative'
    )),
    PRIMARY KEY (household_id, person_id)
);

CREATE INDEX ix_households_project_id         ON households (project_id);
CREATE INDEX ix_households_head_person_id     ON households (head_person_id) WHERE head_person_id IS NOT NULL;
CREATE INDEX ix_household_members_person_id   ON household_members (person_id);
