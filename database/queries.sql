-- name: AllDogs :many
SELECT * FROM "dogs";

-- name: InsertOwner :one
INSERT INTO "owners" (
       "owner_id",
       "owner_name"
) VALUES (
    ?, ?
) RETURNING *;

-- name: InsertDog :one
INSERT INTO "dogs" (
       "dog_id",
       "dog_name",
       "dog_owner"
) VALUES (
    ?, ?, ?
) RETURNING *;
