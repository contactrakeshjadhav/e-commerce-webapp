CREATE TABLE IF NOT EXISTS db_publisher_events
(
    id         BIGSERIAL
        CONSTRAINT db_publisher_events_pkey
            PRIMARY KEY,
    record     TEXT,
    type       TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);
