-- name: Register :exec
INSERT INTO identities (id, username, password_hash)
VALUES ($1, $2, $3);

-- name: GetIdentityById :one
SELECT * FROM identities
WHERE id = $1 LIMIT 1;

-- name: GetIdentityByUsername :one
SELECT * FROM identities
WHERE username = $1 LIMIT 1;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1 LIMIT 1;

-- name: GetRefreshTokenByIdentityId :one
SELECT * FROM refresh_tokens
WHERE identity_id = $1 LIMIT 1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, identity_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteRefreshToken :execrows
DELETE FROM refresh_tokens
WHERE token = $1;

-- name: DeleteIdentityById :exec
DELETE FROM identities
WHERE id = $1;