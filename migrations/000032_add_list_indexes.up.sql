-- Pattern C joins support_records -> people on person_id with a project_id scope.
-- Covers the join key and project scope in one index.
CREATE INDEX ix_support_records_person_project
    ON support_records (person_id, project_id);

-- age group filter in Group 8 reports and person list age_group filter.
-- Only indexes rows where age_group is explicitly set; birth_date-derived bucketing is handled in the query.
CREATE INDEX ix_people_project_age_group
    ON people (project_id, age_group)
    WHERE age_group IS NOT NULL;

-- Low: notes list always orders by created_at DESC per person.
-- The existing ix_person_notes_person_id covers the lookup but not the sort.
CREATE INDEX ix_person_notes_person_created_at
    ON person_notes (person_id, created_at DESC);

-- Low: document list orders by created_at DESC per person. Same reasoning as notes.
CREATE INDEX ix_documents_person_created_at
    ON documents (person_id, created_at DESC);

-- Low: migration records listed per person; migration_date will be used for temporal queries.
CREATE INDEX ix_migration_records_person_date
    ON migration_records (person_id, migration_date);
