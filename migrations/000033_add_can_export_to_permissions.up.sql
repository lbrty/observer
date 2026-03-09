-- Adds explicit export permission flag to project_permissions.
-- Default FALSE: existing project members lose export access until explicitly granted.
-- Platform admins and project owners are unaffected (they bypass permission checks entirely).
-- Staff platform role receives export access automatically in the middleware regardless of this flag.
ALTER TABLE project_permissions
    ADD COLUMN can_export BOOLEAN NOT NULL DEFAULT FALSE;
