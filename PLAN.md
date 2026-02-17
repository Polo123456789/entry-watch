# Superadmin: Condo & Admin Management

## Overview

Two separate management sections under `/super/`:
1. **Condominiums** (`/super/condos`) - List, create, edit condos
2. **Admin Users** (`/super/admins`) - List, create, edit admin accounts

---

## 1. SQL Queries (`db/sqlc/`)

### `condos.sql` (new file)

```sql
-- name: CondoList :many
SELECT * FROM condominiums ORDER BY name;

-- name: CondoGetByID :one
SELECT * FROM condominiums WHERE id = ?;

-- name: CondoCreate :one
INSERT INTO condominiums (name, address, created_at, updated_at, created_by, updated_by)
VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: CondoUpdate :one
UPDATE condominiums SET name = ?, address = ?, updated_at = ?, updated_by = ?
WHERE id = ? RETURNING *;

-- name: CondoDelete :exec
DELETE FROM condominiums WHERE id = ?;
```

### `users.sql` (add queries)

```sql
-- name: UserListByRole :many
SELECT u.*, c.name as condo_name FROM users u
LEFT JOIN condominiums c ON u.condominium_id = c.id
WHERE u.role = ? ORDER BY u.last_name, u.first_name;

-- name: UserUpdate :one
UPDATE users SET first_name = ?, last_name = ?, email = ?, phone = ?,
  condominium_id = ?, enabled = ?, updated_at = ?, updated_by = ?
WHERE id = ? RETURNING *;

-- name: UserDelete :exec
DELETE FROM users WHERE id = ?;
```

---

## 2. Store Layer (`internal/sqlc/`)

### `condo-store.go` (modify - currently has panics)

- `CondoList(ctx) ([]*entry.Condominium, error)`
- `CondoGetByID(ctx, id) (*entry.Condominium, error)`
- `CondoCreate(ctx, condo) (*entry.Condominium, error)`
- `CondoUpdate(ctx, id, updateFn) error`
- `CondoDelete(ctx, id) error`

### `user-store.go` (extend)

- `UserListByRole(ctx, role) ([]*auth.User, error)` - returns admins with condo name
- `UserUpdate(ctx, id, user, updatedBy) (*auth.User, error)`
- `UserDelete(ctx, id) error`

---

## 3. Domain Layer (`internal/entry/`)

### `condominium.go` - extend interface

```go
type CondominiumStore interface {
    CondoList(ctx) ([]*Condominium, error)
    CondoGetByID(ctx, id) (*Condominium, error)
    CondoCreate(ctx, condo) (*Condominium, error)
    CondoUpdate(ctx, id, updateFn) error
    CondoDelete(ctx, id) error
}
```

---

## 4. HTTP Handlers (`internal/http/superadmin/`)

| File | Routes |
|------|--------|
| `condos.go` | GET/POST `/condos`, GET/POST `/condos/{id}`, POST `/condos/{id}/delete` |
| `admins.go` | GET/POST `/admins`, GET/POST `/admins/{id}`, POST `/admins/{id}/delete` |

### Routes Detail

| Route | Handler | Purpose |
|-------|---------|---------|
| `GET /super/condos` | `hCondosList` | Table of condos |
| `GET /super/condos/new` | `hCondosNew` | New condo form |
| `POST /super/condos` | `hCondosCreate` | Create condo |
| `GET /super/condos/{id}/edit` | `hCondosEdit` | Edit condo form |
| `POST /super/condos/{id}` | `hCondosUpdate` | Update condo |
| `POST /super/condos/{id}/delete` | `hCondosDelete` | Delete condo |
| `GET /super/admins` | `hAdminsList` | Table of admins |
| `GET /super/admins/new` | `hAdminsNew` | New admin form |
| `POST /super/admins` | `hAdminsCreate` | Create admin |
| `GET /super/admins/{id}/edit` | `hAdminsEdit` | Edit admin form |
| `POST /super/admins/{id}` | `hAdminsUpdate` | Update admin |
| `POST /super/admins/{id}/delete` | `hAdminsDelete` | Delete admin |

---

## 5. Templates (`internal/templates/superadmin/`)

| File | Purpose |
|------|---------|
| `condos.templ` | Condos table (Name, Address, Edit) |
| `condo-form.templ` | Create/Edit condo form |
| `admins.templ` | Admins table (Name, Email, Condo, Status, Edit) |
| `admin-form.templ` | Create/Edit admin form with condo dropdown |

---

## 6. Execution Order

1. `db/sqlc/condos.sql` - create
2. `db/sqlc/users.sql` - add admin queries
3. `make sqlc`
4. `internal/sqlc/condo-store.go` - implement
5. `internal/sqlc/user-store.go` - extend
6. `internal/entry/condominium.go` - extend interface
7. `internal/http/superadmin/condos.go` - create handlers
8. `internal/http/superadmin/admins.go` - create handlers
9. `internal/templates/superadmin/*.templ` - create templates
10. `make templates`
11. `internal/http/superadmin/routes.go` - wire routes
12. `make tidy && make audit`

---

## Notes

- Naming convention: `ObjectAction` format for store methods (e.g., `CondoList`, `UserUpdate`)
- Admin management extends existing `UserStore` rather than creating a separate store
- Condos table shows: Name, Address
- Admins table shows: Name, Email, Condo, Status
