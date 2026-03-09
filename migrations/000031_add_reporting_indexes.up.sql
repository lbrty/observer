-- Critical: registration-window reports filter project_id + registered_at on every Group 3/4/5/9/10 query.
-- The existing partial index on registered_at alone is not used when project_id leads the filter.
CREATE INDEX ix_people_project_registered_at
    ON people (project_id, registered_at)
    WHERE registered_at IS NOT NULL;

-- Critical: all consultation-based reports (Groups 1/2/6/7/8) filter project_id + provided_at range.
-- The existing solo ix_support_records_provided_at is not used with a leading project_id condition.
CREATE INDEX ix_support_records_project_provided_at
    ON support_records (project_id, provided_at);

-- Important: Groups 1/6/7/8 further filter AND type = 'legal'/'social' inside the project scope.
-- The existing solo ix_support_records_type is not composite with project_id.
CREATE INDEX ix_support_records_project_type
    ON support_records (project_id, type);

-- Important: IDP status derivation joins origin_place_id → places → states.conflict_zone (Group 3, reports 5–11).
CREATE INDEX ix_people_project_origin_place_id
    ON people (project_id, origin_place_id)
    WHERE origin_place_id IS NOT NULL;

-- Important: region-of-stay reports join current_place_id → places → states.name (Group 5, reports 21–23).
CREATE INDEX ix_people_project_current_place_id
    ON people (project_id, current_place_id)
    WHERE current_place_id IS NOT NULL;

-- Important: sex breakdown reports (Group 2, reports 12–17) and person list sex filter.
CREATE INDEX ix_people_project_sex
    ON people (project_id, sex)
    WHERE sex IS NOT NULL;
