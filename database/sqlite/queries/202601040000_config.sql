-- name: GetConfigs :many
SELECT * FROM config;

-- name: SetConfig :exec
INSERT INTO config (name, value)
VALUES (?, ?)
ON CONFLICT (name)
DO UPDATE SET value = EXCLUDED.value;

-- name: GetConfig :one
SELECT value FROM config WHERE name = ?;
