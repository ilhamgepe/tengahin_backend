CREATE TABLE IF NOT EXISTS roles (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(50) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);