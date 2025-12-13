-- name: CreateUser :exec
INSERT INTO users (
    id, email, password, first_name, last_name, role, is_active, created_at, updated_at
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9
         );

-- name: GetUserByID :one
SELECT id, email, password, first_name, last_name, role, is_active, created_at, updated_at
FROM users
WHERE id = $1 AND is_active = true;

-- name: GetUserByEmail :one
SELECT id, email, password, first_name, last_name, role, is_active, created_at, updated_at
FROM users
WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users
SET
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    email = COALESCE($4, email),
    updated_at = $5
WHERE id = $1;

-- name: DeleteUser :exec
UPDATE users
SET is_active = false, updated_at = $2
WHERE id = $1;

-- name: ListUsers :many
SELECT id, email, password, first_name, last_name, role, is_active, created_at, updated_at
FROM users
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountActiveUsers :one
SELECT COUNT(*) FROM users WHERE is_active = true;
