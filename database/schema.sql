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

CREATE TABLE temp_keys (
  temp_key_email VARCHAR(63) NOT NULL,
  temp_key VARCHAR(63) NOT NULL,
  temp_key_expires_unix INTEGER NOT NULL,

  CONSTRAINT pk_temp_keys PRIMARY KEY (temp_key_email),
  CONSTRAINT uk_temp_keys UNIQUE (temp_key)
);
