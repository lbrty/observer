ALTER TABLE migration_records ADD COLUMN updated_at TIMESTAMPTZ;
ALTER TABLE person_notes ADD COLUMN updated_at TIMESTAMPTZ;
ALTER TABLE documents ADD COLUMN updated_at TIMESTAMPTZ;
