-- name: SetCache :exec
INSERT INTO cache (key, value) VALUES (?, ?);

-- name: GetCache :one
SELECT value FROM cache WHERE key = ?;

-- name: DelCache :exec
DELETE FROM cache WHERE key = ?;

-- name: GetAllCache :many
SELECT * FROM cache;
