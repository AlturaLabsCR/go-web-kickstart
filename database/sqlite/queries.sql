-- name: InsertOwner :one
INSERT INTO "owners" ("owner_email") VALUES (?) RETURNING owner_id;

-- name: SelectOwnerEmails :many
SELECT "owner_email" FROM "owners";

-- name: InsertDog :one
INSERT INTO "dogs" (
  "dog_name",
  "dog_owner"
) VALUES (?, ?) RETURNING "dog_id";
