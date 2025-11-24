-- name: InsertUser :one
INSERT INTO "users" ("user_email") VALUES ($1) RETURNING user_id;

-- name: SelectUserEmails :many
SELECT "user_email" FROM "users";

-- name: InsertDog :one
INSERT INTO "dogs" (
  "dog_name",
  "dog_owner"
) VALUES ($1, $2) RETURNING "dog_id";
