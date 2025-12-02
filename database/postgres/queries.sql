-- name: UpsertSession :exec
INSERT INTO "sessions" (
  "session_id",
  "session_user",
  "session_last_used_at",
  "session_csrf_token"
) VALUES (
  $1,
  $2,
  CURRENT_TIMESTAMP,
  $3
) ON CONFLICT ("session_id") DO UPDATE SET
  "session_last_used_at" = EXCLUDED."session_last_used_at",
  "session_csrf_token" = EXCLUDED."session_csrf_token"
;

-- name: SelectSession :one
SELECT * FROM "sessions" WHERE "session_id" = $1;

-- name: GetSessions :many
SELECT * FROM "sessions";

-- name: DeleteSession :exec
DELETE FROM "sessions" WHERE "session_id" = $1;

-- name: UpsertObject :exec
INSERT INTO "objects" (
  "object_key",
  "object_bucket",
  "object_mime",
  "object_size",
  "object_modified"
) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
ON CONFLICT ("object_key") DO UPDATE SET
  "object_mime" = EXCLUDED."object_mime",
  "object_size" = EXCLUDED."object_size",
  "object_modified" = EXCLUDED."object_modified"
;

-- name: SelectObject :one
SELECT * FROM "objects" WHERE "object_key" = $1;

-- name: GetObjects :many
SELECT * FROM "objects";

-- name: DeleteObject :exec
DELETE FROM "objects" WHERE "object_key" = $1;

-- name: UpsertUser :exec
INSERT INTO "users" ("user_id")
VALUES ($1)
ON CONFLICT DO NOTHING;
