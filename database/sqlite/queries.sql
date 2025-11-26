-- name: UpsertSession :exec
INSERT INTO "sessions" (
  "session_id",
  "session_user",
  "session_last_used_at",
  "session_csrf_token"
) VALUES (
  ?,
  ?,
  (unixepoch('now')),
  ?
) ON CONFLICT ("session_id") DO UPDATE SET
  "session_last_used_at" = EXCLUDED."session_last_used_at",
  "session_csrf_token" = EXCLUDED."session_csrf_token"
;

-- name: SelectSession :one
SELECT * FROM "sessions" WHERE "session_id" = ?;

-- name: DeleteSession :exec
DELETE FROM "sessions" WHERE "session_id" = ?;

-- name: UpsertUser :exec
INSERT INTO "users" ("user_id")
VALUES (?)
ON CONFLICT DO NOTHING;
