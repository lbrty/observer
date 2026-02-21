-- person_categories replaces the single people.category_id FK.
-- A person may hold multiple concurrent vulnerability categories,
-- mirroring the UNHCR PSN multi-code model (11 categories, 71 subcodes).
-- This prevents combinatorial explosion in the categories table
-- and avoids losing secondary vulnerability dimensions.
CREATE TABLE person_categories (
    person_id   TEXT NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    category_id TEXT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    PRIMARY KEY (person_id, category_id)
);

CREATE INDEX ix_person_categories_category_id ON person_categories (category_id);
