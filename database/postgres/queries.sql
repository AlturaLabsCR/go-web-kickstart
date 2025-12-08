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
  "object_public_url",
  "object_mime",
  "object_size",
  "object_modified"
) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
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

-- name: InsertPermission :exec
INSERT INTO "permissions" ("permission_name") VALUES ($1);

-- name: InsertRole :exec
INSERT INTO "roles" ("role_name") VALUES ($1);

-- name: InsertRolePermission :exec
INSERT INTO "role_permissions" (
  "role_permission_role",
  "role_permission_permission"
) VALUES ($1, $2);

-- name: GetPermissions :many
SELECT DISTINCT rp."role_permission_permission"
FROM "user_roles" ur
JOIN "role_permissions" rp
ON rp."role_permission_role" = ur."user_role_role"
WHERE ur."user_role_user" = $1;
