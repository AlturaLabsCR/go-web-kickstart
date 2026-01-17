-- name: UpsertUserName :exec
INSERT INTO user_names ("user", name)
VALUES (?, ?)
ON CONFLICT ("user") DO UPDATE
SET name = EXCLUDED.name;

-- name: GetUserName :one
SELECT name FROM user_names WHERE "user" = ?;
