CREATE TABLE IF NOT EXISTS identities
(
    id            UUID               NOT NULL PRIMARY KEY,
    username      VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT                NOT NULL,
    created_at    TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS refresh_tokens
(
    token       BYTEA       NOT NULL PRIMARY KEY,
    identity_id UUID       NOT NULL REFERENCES identities (id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at  TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_identities_username ON identities (username);
CREATE UNIQUE INDEX IF NOT EXISTS idx_refresh_tokens_identity_id ON refresh_tokens (identity_id);