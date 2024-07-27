CREATE TABLE IF NOT EXISTS object_has_tags
(
    id          TEXT NOT NULL
        PRIMARY KEY,
    object_id         TEXT,
    object_type         TEXT,
    tag_id         TEXT
        CONSTRAINT fk_object_has_tags_tags
        REFERENCES tags
);