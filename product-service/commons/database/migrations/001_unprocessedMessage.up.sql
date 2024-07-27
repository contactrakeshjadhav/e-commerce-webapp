CREATE TABLE IF NOT EXISTS unprocessed_messages
(
    event_id   TEXT
        CONSTRAINT unprocessed_messages_pkey
            PRIMARY KEY,
    delivered  BIGINT,
    service    TEXT,
    consumer   TEXT,
    payload    TEXT,
    outcome    TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);
