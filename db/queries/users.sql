-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT ? OFFSET ?;

-- name: CreateUser :execresult
INSERT INTO users (
    email, password_hash, first_name, last_name, status
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: UpdateUser :exec
UPDATE users
SET 
    email = ?,
    first_name = ?,
    last_name = ?,
    status = ?
WHERE id = ?;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: GetUsersByStatus :many
SELECT * FROM users
WHERE status = ?
ORDER BY id
LIMIT ? OFFSET ?; 