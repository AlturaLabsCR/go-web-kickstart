-- name: InsertOwner :one
INSERT INTO "owners" ("owner_name") VALUES ($1) RETURNING owner_id;

-- name: InsertDog :one
INSERT INTO "dogs" (
       "dog_name",
       "dog_owner"
) VALUES ($1, $2) RETURNING "dog_id";
