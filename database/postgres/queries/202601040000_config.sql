-- name: GetConfigs :many
SELECT * FROM config;

-- name: SetConfig :exec
INSERT INTO config (name, value)
VALUES ($1, $2)
ON CONFLICT (name)
DO UPDATE SET value = EXCLUDED.value;

-- name: GetConfig :one
SELECT value FROM config WHERE name = $1;
