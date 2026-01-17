-- name: GetUserMeta :one
SELECT * FROM user_meta WHERE id = ?;

-- name: GetUsersMeta :many
SELECT * FROM user_meta;
