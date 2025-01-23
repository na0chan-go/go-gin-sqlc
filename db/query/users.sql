-- name: SearchUsers :many
SELECT id, email, password_hash, first_name, last_name, status, created_at, updated_at
FROM users
WHERE (
    email LIKE CONCAT('%', ?, '%') OR
    first_name LIKE CONCAT('%', ?, '%') OR
    last_name LIKE CONCAT('%', ?, '%')
)
AND (? IS NULL OR status = ?)
ORDER BY id
LIMIT ? OFFSET ?;

-- name: CreatePasswordReset :execresult
INSERT INTO password_resets (
    user_id, token, expires_at
) VALUES (
    ?, ?, ?
);

-- name: GetPasswordResetByToken :one
SELECT user_id, token, expires_at, created_at
FROM password_resets
WHERE token = ? AND expires_at > NOW()
LIMIT 1;

-- name: DeletePasswordReset :exec
DELETE FROM password_resets
WHERE token = ?; 