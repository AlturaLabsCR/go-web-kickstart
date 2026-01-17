-- name: GetUserMeta :one
SELECT * FROM user_meta WHERE id = ?;
