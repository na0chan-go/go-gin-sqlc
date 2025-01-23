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