CREATE TABLE pet_tags (
    pet_id TEXT NOT NULL REFERENCES pets (id) ON DELETE CASCADE,
    tag_id TEXT NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (pet_id, tag_id)
);

CREATE INDEX ix_pet_tags_tag_id ON pet_tags (tag_id);
