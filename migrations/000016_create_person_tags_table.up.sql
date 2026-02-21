CREATE TABLE person_tags (
    person_id TEXT NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    tag_id    TEXT NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (person_id, tag_id)
);

CREATE INDEX ix_person_tags_tag_id ON person_tags (tag_id);
