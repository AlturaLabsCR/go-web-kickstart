-- DDL

CREATE TABLE "owners" (
  "owner_id" INTEGER PRIMARY KEY,
  "owner_created_at" INTEGER NOT NULL DEFAULT (unixepoch('now')),
  "owner_email" VARCHAR(64) NOT NULL UNIQUE
);

CREATE TABLE "dogs" (
  "dog_id" INTEGER PRIMARY KEY,
  "dog_owner" INTEGER NOT NULL REFERENCES "owners"("owner_id"),
  "dog_created_at" INTEGER NOT NULL DEFAULT (unixepoch('now')),
  "dog_name" VARCHAR(64) NOT NULL,
  "dog_weight" REAL NOT NULL
);
