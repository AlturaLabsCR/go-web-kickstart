-- name: SetUser :exec
INSERT INTO users (id, created_at) VALUES (?, ?);

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: DelUser :exec
DELETE FROM users WHERE id = ?;
