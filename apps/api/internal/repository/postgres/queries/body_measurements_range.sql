-- name: ListBodyMeasurementsByUserTypeRange :many
SELECT m.id, m.check_in_id, m.measurement_type, m.side, m.value, c.date::date as check_in_date, m.created_at, m.updated_at
FROM body_measurements m
JOIN body_check_ins c ON c.id = m.check_in_id
WHERE c.user_id = $1
  AND m.measurement_type = $2
  AND c.date >= $3::date
  AND c.date <= $4::date
ORDER BY c.date ASC, m.created_at ASC;