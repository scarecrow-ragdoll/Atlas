-- name: CreateAiExport :one
INSERT INTO ai_exports (user_id, date_range_start, date_range_end, include_photos, include_nutrition, include_cardio, include_measurements, user_comment, generated_prompt, export_file_path)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, user_id, date_range_start, date_range_end, include_photos, include_nutrition, include_cardio, include_measurements, user_comment, generated_prompt, export_file_path, created_at, updated_at;

-- name: GetAiExportByID :one
SELECT id, user_id, date_range_start, date_range_end, include_photos, include_nutrition, include_cardio, include_measurements, user_comment, generated_prompt, export_file_path, created_at, updated_at
FROM ai_exports
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListAiExportsByUserID :many
SELECT id, user_id, date_range_start, date_range_end, include_photos, include_nutrition, include_cardio, include_measurements, user_comment, generated_prompt, export_file_path, created_at, updated_at
FROM ai_exports
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateAiExportFilePath :one
UPDATE ai_exports
SET export_file_path = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, date_range_start, date_range_end, include_photos, include_nutrition, include_cardio, include_measurements, user_comment, generated_prompt, export_file_path, created_at, updated_at;

-- name: DeleteAiExport :one
DELETE FROM ai_exports
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date_range_start, date_range_end, include_photos, include_nutrition, include_cardio, include_measurements, user_comment, generated_prompt, export_file_path, created_at, updated_at;

-- name: ListStaleAiExports :many
SELECT id, user_id, date_range_start, date_range_end, include_photos, include_nutrition, include_cardio, include_measurements, user_comment, generated_prompt, export_file_path, created_at, updated_at
FROM ai_exports
WHERE created_at < NOW() - $1::INTERVAL
ORDER BY created_at ASC;