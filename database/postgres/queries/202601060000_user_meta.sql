-- name: GetUserMeta :one
SELECT * FROM user_meta WHERE id = $1;

-- name: GetUsersMeta :many
SELECT * FROM user_meta;
