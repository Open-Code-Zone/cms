-- name: CreateCollectionItem :one
INSERT INTO collections (filename, collection_name,  content, metadata)
VALUES (?, ?, ?, ?)
RETURNING filename, collection_name, content, metadata, created_at;

-- name: GetCollectionItem :one
SELECT filename, collection_name, content, metadata, created_at
FROM collections
WHERE filename = ? AND collection_name = ?;

-- name: UpdateCollectionItem :exec
UPDATE collections
SET content = ?, metadata = ?
WHERE filename = ? AND collection_name = ?;

-- name: DeleteCollectionItem :exec
DELETE FROM collections
WHERE filename = ? AND collection_name = ?;

-- name: ListAllCollectionItems :many
SELECT filename, content, metadata, created_at
FROM collections
WHERE collection_name = ?
ORDER BY created_at DESC;
