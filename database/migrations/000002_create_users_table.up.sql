CREATE TABLE IF NOT EXISTS users
(
    id         UUID PRIMARY KEY REFERENCES identities (id) ON DELETE CASCADE,
    username   VARCHAR(100) UNIQUE NOT NULL REFERENCES identities (username) ON DELETE CASCADE,
    created_at TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_username ON users (username)