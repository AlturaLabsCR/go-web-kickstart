-- DDL

CREATE TABLE "users" (
  "user_id" INTEGER PRIMARY KEY,
  "user_created_at" INTEGER NOT NULL DEFAULT (unixepoch('now')),
  "user_email" VARCHAR(64) NOT NULL UNIQUE
);

CREATE TABLE "dogs" (
  "dog_id" INTEGER PRIMARY KEY,
  "dog_owner" INTEGER NOT NULL REFERENCES "owners"("owner_id"),
  "dog_created_at" INTEGER NOT NULL DEFAULT (unixepoch('now')),
  "dog_name" VARCHAR(64) NOT NULL,
  "dog_weight" REAL NOT NULL
);
