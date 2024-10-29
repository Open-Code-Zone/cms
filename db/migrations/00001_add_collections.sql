-- +goose Up
CREATE TABLE collections (
    filename TEXT PRIMARY KEY,
    collection_name TEXT NOT NULL,
    content TEXT NOT NULL,
    metadata TEXT NOT NULL, -- JSON unstructured data
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

-- Create an index on the filename for faster lookups
CREATE INDEX idx_filename ON collections(filename, created_at);
CREATE INDEX idx_collection_type ON collections(collection_name, created_at);

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE collections;

-- +goose StatementBegin
-- +goose StatementEnd
