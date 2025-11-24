-- name: InsertUser :one
INSERT INTO "users" ("user_email") VALUES (?) RETURNING user_id;

-- name: SelectUserEmails :many
SELECT "user_email" FROM "users";

-- name: InsertDog :one
INSERT INTO "dogs" (
  "dog_name",
  "dog_owner"
) VALUES (?, ?) RETURNING "dog_id";
