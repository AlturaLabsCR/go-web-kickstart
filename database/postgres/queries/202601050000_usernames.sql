-- name: UpsertUserName :exec
INSERT INTO user_names ("user", name)
VALUES ($2, $1)
ON CONFLICT ("user") DO UPDATE
SET name = EXCLUDED.name;

-- name: GetUserName :one
SELECT name FROM user_names WHERE "user" = $1;
