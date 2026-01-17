-- name: SetCache :exec
INSERT INTO cache (key, value)
VALUES (?, ?)
ON CONFLICT(key) DO
UPDATE SET value = excluded.value;

-- name: GetCache :one
SELECT value FROM cache WHERE key = ?;

-- name: DelCache :exec
DELETE FROM cache WHERE key = ?;

-- name: GetAllCache :many
SELECT * FROM cache;
