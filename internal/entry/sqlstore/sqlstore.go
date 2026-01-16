package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"log/slog"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

// Ensure sqlstore implements entry.Store via embedding the required interfaces
type Store struct {
	db     *sql.DB
	logger *slog.Logger
}

func New(db *sql.DB, logger *slog.Logger) *Store {
	return &Store{db: db, logger: logger}
}

var ErrNotFound = entry.ErrNotFound

// helper to scan common user fields
func scanUser(row scannable) (*entry.StoreUser, error) {
	var u entry.StoreUser
	var createdAt int64
	var updatedAt int64
	// columns: id, condominium_id, first_name, last_name, email, password, role, enabled, hidden, created_at, updated_at, created_by, updated_by
	if err := row.Scan(&u.ID, &u.CondominiumID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.Role, &u.Enabled, &u.Hidden, &createdAt, &updatedAt, &u.CreatedBy, &u.UpdatedBy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.CreatedAt = time.Unix(createdAt, 0)
	u.UpdatedAt = time.Unix(updatedAt, 0)
	return &u, nil
}

// UserGetByID implements entry.UserStore
func (s *Store) UserGetByID(ctx context.Context, id int64) (*entry.StoreUser, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, condominium_id, first_name, last_name, email, password, role, enabled, hidden, created_at, updated_at, created_by, updated_by FROM users WHERE id = ?`, id)
	return scanUser(row)
}

// UserGetByEmail implements entry.UserStore
func (s *Store) UserGetByEmail(ctx context.Context, email string) (*entry.StoreUser, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, condominium_id, first_name, last_name, email, password, role, enabled, hidden, created_at, updated_at, created_by, updated_by FROM users WHERE email = ?`, email)
	return scanUser(row)
}

// UserCreate implements entry.UserStore
func (s *Store) UserCreate(ctx context.Context, u *entry.StoreUser) (*entry.StoreUser, error) {
	now := time.Now().Unix()
	var condoID interface{}
	if u.CondominiumID == 0 {
		condoID = nil
	} else {
		condoID = u.CondominiumID
	}
	res, err := s.db.ExecContext(ctx, `INSERT INTO users (condominium_id, first_name, last_name, email, password, role, enabled, hidden, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, condoID, u.FirstName, u.LastName, u.Email, u.PasswordHash, u.Role, u.Enabled, u.Hidden, now, now, u.CreatedBy, u.UpdatedBy)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.UserGetByID(ctx, id)
}

// UserCountByRole implements entry.UserStore
func (s *Store) UserCountByRole(ctx context.Context, role entry.UserRole) (int64, error) {
	var c int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE role = ?`, role).Scan(&c)
	if err != nil {
		return 0, err
	}
	return c, nil
}

// Condominium methods
func (s *Store) CondoGetByID(ctx context.Context, id int64) (*entry.Condominium, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, name, address, created_at, updated_at, created_by, updated_by FROM condominiums WHERE id = ?`, id)
	var c entry.Condominium
	var createdAt int64
	var updatedAt int64
	if err := row.Scan(&c.ID, &c.Name, &c.Address, &createdAt, &updatedAt, &c.CreatedBy, &c.UpdatedBy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	c.CreatedAt = time.Unix(createdAt, 0)
	c.UpdatedAt = time.Unix(updatedAt, 0)
	return &c, nil
}

func (s *Store) CondoCreate(ctx context.Context, condo *entry.Condominium) (*entry.Condominium, error) {
	now := time.Now().Unix()
	res, err := s.db.ExecContext(ctx, `INSERT INTO condominiums (name, address, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)`, condo.Name, condo.Address, now, now, condo.CreatedBy, condo.UpdatedBy)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.CondoGetByID(ctx, id)
}

func (s *Store) CondoUpdate(ctx context.Context, id int64, updateFn func(condo *entry.Condominium) (*entry.Condominium, error)) error {
	// load
	c, err := s.CondoGetByID(ctx, id)
	if err != nil {
		return err
	}
	updated, err := updateFn(c)
	if err != nil {
		return err
	}
	updated.UpdatedAt = time.Now()
	_, err = s.db.ExecContext(ctx, `UPDATE condominiums SET name = ?, address = ?, updated_at = ?, created_by = ?, updated_by = ? WHERE id = ?`, updated.Name, updated.Address, updated.UpdatedAt.Unix(), updated.CreatedBy, updated.UpdatedBy, id)
	return err
}

// Visit methods
func (s *Store) VisitGetByID(ctx context.Context, id string) (*entry.Visit, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, condominium_id, user_id, visitor_name, max_uses, uses, valid_from, valid_to, created_at, updated_at FROM visits WHERE id = ?`, id)
	var v entry.Visit
	var validFrom int64
	var validTo int64
	var createdAt int64
	var updatedAt int64
	if err := row.Scan(&v.ID, &v.CondominiumID, &v.UserID, &v.VisitorName, &v.MaxUses, &v.Uses, &validFrom, &validTo, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	v.ValidFrom = time.Unix(validFrom, 0)
	v.ValidTo = time.Unix(validTo, 0)
	v.CreatedAt = time.Unix(createdAt, 0)
	v.UpdatedAt = time.Unix(updatedAt, 0)
	return &v, nil
}

func (s *Store) VisitCreate(ctx context.Context, visit *entry.Visit) (*entry.Visit, error) {
	if visit.ID == "" {
		return nil, errors.New("visit id required")
	}
	now := time.Now().Unix()
	_, err := s.db.ExecContext(ctx, `INSERT INTO visits (id, condominium_id, user_id, visitor_name, max_uses, uses, valid_from, valid_to, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, visit.ID, visit.CondominiumID, visit.UserID, visit.VisitorName, visit.MaxUses, visit.Uses, visit.ValidFrom.Unix(), visit.ValidTo.Unix(), now, now)
	if err != nil {
		return nil, err
	}
	return s.VisitGetByID(ctx, visit.ID)
}

func (s *Store) VisitUpdate(ctx context.Context, id string, updateFn func(visit *entry.Visit) (*entry.Visit, error)) error {
	v, err := s.VisitGetByID(ctx, id)
	if err != nil {
		return err
	}
	updated, err := updateFn(v)
	if err != nil {
		return err
	}
	updated.UpdatedAt = time.Now()
	_, err = s.db.ExecContext(ctx, `UPDATE visits SET uses = ?, max_uses = ?, valid_from = ?, valid_to = ?, updated_at = ? WHERE id = ?`, updated.Uses, updated.MaxUses, updated.ValidFrom.Unix(), updated.ValidTo.Unix(), updated.UpdatedAt.Unix(), id)
	return err
}

// scannable interface to accept *sql.Row and *sql.Rows
type scannable interface {
	Scan(dest ...interface{}) error
}
