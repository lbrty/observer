CREATE TABLE person_status_history (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    person_id   TEXT NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    from_status TEXT NOT NULL,
    to_status   TEXT NOT NULL,
    changed_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_psh_person ON person_status_history(person_id);
CREATE INDEX idx_psh_changed ON person_status_history(changed_at);

CREATE FUNCTION track_person_status_change() RETURNS TRIGGER AS $$
BEGIN
    IF OLD.case_status IS DISTINCT FROM NEW.case_status THEN
        INSERT INTO person_status_history(person_id, from_status, to_status)
        VALUES (NEW.id, OLD.case_status, NEW.case_status);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_person_status_change
    AFTER UPDATE ON people
    FOR EACH ROW EXECUTE FUNCTION track_person_status_change();
