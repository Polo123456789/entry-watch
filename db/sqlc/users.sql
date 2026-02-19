-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?;

-- name: CreateUser :one
INSERT INTO users (
    condominium_id,
	first_name,
	last_name,
	email,
	phone,
	role,
	password,
	enabled,
	hidden,
	created_at,
	updated_at,
	created_by,
	updated_by
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: CountSuperAdmins :one
SELECT COUNT(*) AS count
FROM users
WHERE role = 'superadmin' AND enabled = 1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = ?, updated_at = ?, updated_by = ?
WHERE id = ?;

-- name: UserListByRole :many
SELECT u.*, c.name AS condo_name
FROM users u
LEFT JOIN condominiums c ON u.condominium_id = c.id
WHERE u.role = ?
ORDER BY u.last_name, u.first_name;

-- name: UserUpdate :one
UPDATE users
SET first_name = ?, last_name = ?, email = ?, phone = ?,
    condominium_id = ?, enabled = ?, updated_at = ?, updated_by = ?
WHERE id = ?
RETURNING *;

-- name: UserDelete :exec
DELETE FROM users
WHERE id = ?;

-- name: CountUsersByCondo :one
SELECT COUNT(*) AS count
FROM users
WHERE condominium_id = ?;
