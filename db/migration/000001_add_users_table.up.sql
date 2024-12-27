CREATE TABLE IF NOT EXISTS users (
    "id" SERIAL PRIMARY KEY,
    "email" varchar UNIQUE NOT NULL,
    "username" VARCHAR(255) NOT NULL,
    "fullname" varchar NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "password_change_at" timestamptz NOT NULL DEFAULT NOW(),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");
CREATE INDEX ON "users" ("fullname");
