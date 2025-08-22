-- DDL

CREATE TABLE "owners" (
       "owner_id"    INTEGER PRIMARY KEY NOT NULL,
       "owner_name"  TEXT    NOT NULL
);

CREATE TABLE "dogs" (
       "dog_id"    INTEGER PRIMARY KEY NOT NULL,
       "dog_name"  TEXT    NOT NULL,
       "dog_owner" INTEGER NOT NULL REFERENCES "owners"("owner_id")
);
