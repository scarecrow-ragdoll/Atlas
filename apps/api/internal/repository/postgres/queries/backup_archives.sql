-- name: CreateBackupArchive :one
INSERT INTO backup_archives (user_id, include_media, size_bytes, entity_counts)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, include_media, size_bytes, entity_counts, archive_path, created_at, updated_at;

-- name: GetBackupArchiveByID :one
SELECT id, user_id, include_media, size_bytes, entity_counts, archive_path, created_at, updated_at
FROM backup_archives
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: UpdateBackupArchiveFilePath :one
UPDATE backup_archives
SET archive_path = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, include_media, size_bytes, entity_counts, archive_path, created_at, updated_at;