-- DDL

CREATE TABLE "owners" (
  "owner_id" BIGSERIAL PRIMARY KEY,
  "owner_created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "owner_name" VARCHAR(64) NOT NULL
);

CREATE TABLE "dogs" (
  "dog_id" BIGSERIAL PRIMARY KEY,
  "dog_owner" BIGINT NOT NULL REFERENCES "owners"("owner_id"),
  "dog_created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "dog_name" VARCHAR(64) NOT NULL,
  "dog_weight" REAL NOT NULL
);
