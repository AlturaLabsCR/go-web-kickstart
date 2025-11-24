-- name: InsertOwner :one
INSERT INTO "owners" ("owner_name") VALUES (?) RETURNING owner_id;

-- name: InsertDog :one
INSERT INTO "dogs" (
  "dog_name",
  "dog_owner"
) VALUES (?, ?) RETURNING "dog_id";
