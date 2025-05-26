-- name: CreateUser :exec
INSERT INTO users (id, username)
VALUES ($1, $2);

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;