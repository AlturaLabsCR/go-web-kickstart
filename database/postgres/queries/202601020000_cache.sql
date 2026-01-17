-- name: SetCache :exec
INSERT INTO cache (key, value)
VALUES ($1, $2)
ON CONFLICT(key) DO
UPDATE SET value = excluded.value;

-- name: GetCache :one
SELECT value FROM cache WHERE key = $1;

-- name: DelCache :exec
DELETE FROM cache WHERE key = $1;

-- name: GetAllCache :many
SELECT * FROM cache;
