package entry

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

// MemStore is an in-memory implementation of Store for local/dev usage.
type MemStore struct {
	mu sync.Mutex

	nextUserID   int64
	users        map[int64]*StoreUser
	usersByEmail map[string]*StoreUser

	nextCondoID int64
	condos      map[int64]*Condominium

	visits map[string]*Visit
}

func NewMemStore() *MemStore {
	return &MemStore{
		users:        make(map[int64]*StoreUser),
		usersByEmail: make(map[string]*StoreUser),
		condos:       make(map[int64]*Condominium),
		visits:       make(map[string]*Visit),
		nextUserID:   1,
		nextCondoID:  1,
	}
}

// UserStore implementation
func (m *MemStore) UserGetByID(ctx context.Context, id int64) (*StoreUser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneStoreUser(u), nil
}

func (m *MemStore) UserGetByEmail(ctx context.Context, email string) (*StoreUser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.usersByEmail[email]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneStoreUser(u), nil
}

func (m *MemStore) UserCreate(ctx context.Context, u *StoreUser) (*StoreUser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if u.ID == 0 {
		u.ID = m.nextUserID
		m.nextUserID++
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	copy := cloneStoreUser(u)
	m.users[u.ID] = copy
	m.usersByEmail[u.Email] = copy
	return cloneStoreUser(copy), nil
}

func (m *MemStore) UserCountByRole(ctx context.Context, role UserRole) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var c int64
	for _, u := range m.users {
		if u.Role == role {
			c++
		}
	}
	return c, nil
}

func cloneStoreUser(u *StoreUser) *StoreUser {
	if u == nil {
		return nil
	}
	c := *u
	return &c
}

// CondominiumStore implementation (minimal)
func (m *MemStore) CondoGetByID(ctx context.Context, id int64) (*Condominium, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.condos[id]
	if !ok {
		return nil, ErrNotFound
	}
	copy := *c
	return &copy, nil
}

func (m *MemStore) CondoCreate(ctx context.Context, condo *Condominium) (*Condominium, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if condo.ID == 0 {
		condo.ID = m.nextCondoID
		m.nextCondoID++
	}
	now := time.Now()
	condo.CreatedAt = now
	condo.UpdatedAt = now
	c := *condo
	m.condos[condo.ID] = &c
	return &c, nil
}

func (m *MemStore) CondoUpdate(ctx context.Context, id int64, updateFn func(condo *Condominium) (*Condominium, error)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.condos[id]
	if !ok {
		return ErrNotFound
	}
	updated, err := updateFn(cloneCondo(c))
	if err != nil {
		return err
	}
	updated.UpdatedAt = time.Now()
	m.condos[id] = updated
	return nil
}

func cloneCondo(c *Condominium) *Condominium {
	if c == nil {
		return nil
	}
	cc := *c
	return &cc
}

// VisitStore implementation (minimal)
func (m *MemStore) VisitGetByID(ctx context.Context, id string) (*Visit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.visits[id]
	if !ok {
		return nil, ErrNotFound
	}
	copy := *v
	return &copy, nil
}

func (m *MemStore) VisitCreate(ctx context.Context, visit *Visit) (*Visit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if visit.ID == "" {
		return nil, errors.New("visit id required")
	}
	visit.CreatedAt = time.Now()
	visit.UpdatedAt = time.Now()
	v := *visit
	m.visits[visit.ID] = &v
	return &v, nil
}

func (m *MemStore) VisitUpdate(ctx context.Context, id string, updateFn func(visit *Visit) (*Visit, error)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.visits[id]
	if !ok {
		return ErrNotFound
	}
	updated, err := updateFn(cloneVisit(v))
	if err != nil {
		return err
	}
	updated.UpdatedAt = time.Now()
	m.visits[id] = updated
	return nil
}

func cloneVisit(v *Visit) *Visit {
	if v == nil {
		return nil
	}
	vv := *v
	return &vv
}
