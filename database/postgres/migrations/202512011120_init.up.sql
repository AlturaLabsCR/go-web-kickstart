CREATE TABLE IF NOT EXISTS "users" (
  "user_id" TEXT PRIMARY KEY,
  "user_created_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "sessions" (
  "session_id" TEXT PRIMARY KEY,
  "session_user" TEXT NOT NULL REFERENCES "users"("user_id"),
  "session_created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "session_last_used_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "session_csrf_token" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "objects" (
  "object_key" VARCHAR(255) PRIMARY KEY NOT NULL,
  "object_bucket" VARCHAR(255) NOT NULL,
  "object_mime" VARCHAR(64) NOT NULL,
  "object_size" BIGINT NOT NULL,
  "object_created" TIMESTAMP NOT NULL DEFAULT NOW(),
  "object_modified" TIMESTAMP NOT NULL DEFAULT NOW()
);
