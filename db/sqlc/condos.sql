-- name: CondoList :many
SELECT *
FROM condominiums
ORDER BY name;

-- name: CondoGetByID :one
SELECT *
FROM condominiums
WHERE id = ?;

-- name: CondoCreate :one
INSERT INTO condominiums (
    name,
    address,
    created_at,
    updated_at,
    created_by,
    updated_by
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: CondoUpdate :one
UPDATE condominiums
SET name = ?, address = ?, updated_at = ?, updated_by = ?
WHERE id = ?
RETURNING *;

-- name: CondoDelete :exec
DELETE FROM condominiums
WHERE id = ?;
