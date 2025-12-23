CREATE TABLE IF NOT EXISTS "users" (
  "user_id" TEXT PRIMARY KEY,
  "user_created_at" INTEGER NOT NULL DEFAULT (unixepoch('now'))
);

CREATE TABLE IF NOT EXISTS "sessions" (
  "session_id" TEXT PRIMARY KEY,
  "session_user" TEXT NOT NULL REFERENCES "users"("user_id"),
  "session_created_at" INTEGER NOT NULL DEFAULT (unixepoch('now')),
  "session_last_used_at" INTEGER NOT NULL DEFAULT (unixepoch('now')),
  "session_csrf_token" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "objects" (
  "object_key" VARCHAR(64) PRIMARY KEY NOT NULL,
  "object_bucket" VARCHAR(64) NOT NULL,
  "object_public_url" VARCHAR(256) NOT NULL,
  "object_mime" VARCHAR(64) NOT NULL,
  "object_size" INTEGER NOT NULL,
  "object_created" INTEGER NOT NULL DEFAULT (unixepoch('now')),
  "object_modified" INTEGER NOT NULL DEFAULT (unixepoch('now'))
);

CREATE TABLE IF NOT EXISTS "permissions" (
  "permission_name" VARCHAR(32) PRIMARY KEY NOT NULL,
  "permission_description" VARCHAR(512) NOT NULL
);

CREATE TABLE IF NOT EXISTS "roles" (
  "role_name" VARCHAR(32) PRIMARY KEY NOT NULL,
  "role_description" VARCHAR(512) NOT NULL
);

CREATE TABLE IF NOT EXISTS "role_permissions" (
  "role_permission_role" VARCHAR(32) NOT NULL REFERENCES "roles"("role_name") ON DELETE CASCADE,
  "role_permission_permission" VARCHAR(32) NOT NULL REFERENCES "permissions"("permission_name") ON DELETE CASCADE,
  PRIMARY KEY ("role_permission_role", "role_permission_permission")
);

CREATE TABLE IF NOT EXISTS "user_roles" (
  "user_role_user" TEXT NOT NULL REFERENCES "users"("user_id") ON DELETE CASCADE,
  "user_role_role" VARCHAR(32) NOT NULL REFERENCES "roles"("role_name") ON DELETE CASCADE,
  PRIMARY KEY ("user_role_user", "user_role_role")
);
