CREATE TABLE IF NOT EXISTS user_roles(
    "user_id" bigint NOT NULL,
    "role_id" bigint NOT NULL,
    PRIMARY KEY ("user_id", "role_id"),
    FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE CASCADE,
    FOREIGN KEY ("role_id") REFERENCES roles("id") ON DELETE CASCADE
);
