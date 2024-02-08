--migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS meta
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    filename TEXT NOT NULL UNIQUE,
    blob_sequence TEXT NOT NULL,
    mime_type TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS meta_filename ON meta (filename);

-- migrate:down
DROP TABLE IF EXISTS meta;