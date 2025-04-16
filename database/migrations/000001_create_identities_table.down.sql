DROP INDEX IF EXISTS idx_refresh_tokens_identity_id;
DROP INDEX IF EXISTS idx_identities_username;

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS identities;