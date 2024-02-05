CREATE TABLE IF NOT EXISTS meta
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL UNIQUE,
    blob_sequence TEXT NOT NULL,
    mime_type TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS meta_filename ON meta (filename);
