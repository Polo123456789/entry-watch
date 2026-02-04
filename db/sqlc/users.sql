-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?;

-- name: CreateUser :exec
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
);

-- name: CountSuperAdmins :one
SELECT COUNT(*) AS count
FROM users
WHERE role = 'superadmin' AND enabled = 1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = ?, updated_at = ?, updated_by = ?
WHERE id = ?;
