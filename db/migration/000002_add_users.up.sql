CREATE TABLE "users"
(
    "username"             varchar PRIMARY KEY,
    "hashed_password"      varchar        NOT NULL,
    "first_name"           varchar        NOT NULL,
    "last_name"            varchar        NOT NULL,
    "email"                varchar UNIQUE NOT NULL,
    "password_modified_at" timestamptz    NOT NULL DEFAULT (now()),
    "created_at"           timestamptz    NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

CREATE INDEX ON "users" ("username");
