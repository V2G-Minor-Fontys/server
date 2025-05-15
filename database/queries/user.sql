-- name: CreateUser :exec
INSERT INTO users (id, username)
VALUES ($1, $2);

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = $1;