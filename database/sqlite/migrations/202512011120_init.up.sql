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
  "object_key" VARCHAR(255) PRIMARY KEY NOT NULL,
  "object_bucket" VARCHAR(255) NOT NULL,
  "object_mime" VARCHAR(64) NOT NULL,
  "object_size" INTEGER NOT NULL,
  "object_created" INTEGER NOT NULL DEFAULT (unixepoch('now')),
  "object_modified" INTEGER NOT NULL DEFAULT (unixepoch('now'))
);
