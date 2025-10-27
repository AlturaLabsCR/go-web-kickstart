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

-- name: InsertTempKey :exec
INSERT INTO temp_keys (
  temp_key_email,
  temp_key,
  temp_key_expires_unix
) VALUES (?, ?, ?);

-- name: GetTempKey :one
SELECT * FROM temp_keys WHERE temp_key_email = ?;

-- name: UpdateTempKey :exec
UPDATE temp_keys SET
temp_key = ?,
temp_key_expires_unix = ?
WHERE temp_key_email = ?;

-- name: SetTempKeyUsed :exec
UPDATE temp_keys SET temp_key_expires_unix = 0 WHERE temp_key_email = ?;
