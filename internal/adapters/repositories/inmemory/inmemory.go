package inmemory

import (
	"context"
	"fmt"
	"strings"

	"sync"

	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	store map[uuid.UUID]*duser.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{}
}

func (r *InMemoryUserRepository) NextID(ctx context.Context) uuid.UUID {
	return uuid.New()
}

func (r *InMemoryUserRepository) IDExists(
	ctx context.Context,
	id uuid.UUID,
) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.store[id]
	return exists
}

func (r *InMemoryUserRepository) EmailExists(
	ctx context.Context,
	email domain.Email,
) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.store {
		if u.Email() == email {
			return true
		}
	}
	return false
}

func (r *InMemoryUserRepository) GetAll(ctx context.Context) []*duser.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*duser.User, 0, len(r.store))
	for _, u := range r.store {
		result = append(result, u)
	}
	return result
}

func (r *InMemoryUserRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*duser.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, exists := r.store[id]
	if !exists {
		return nil, fmt.Errorf(
			"%w: пользователя с id %s не существует",
			application.ErrNotFound,
			id,
		)
	}
	return u, nil
}

func (r *InMemoryUserRepository) GetByEmail(
	ctx context.Context,
	email domain.Email,
) (*duser.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.store {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf(
		"%w: пользователя с email %s не существует",
		application.ErrNotFound,
		email,
	)
}

func (r *InMemoryUserRepository) Filter(
	ctx context.Context,
	searchByEmail string,
	filterByState []duser.UserState,
	filterByStatus []duser.UserStatus,
	limit, offset int,
) ([]*duser.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	searchFlag := len(searchByEmail) > 0
	filterByStateFlag := len(filterByState) > 0
	filterByStatusFlag := len(filterByStatus) > 0
	result := make([]*duser.User, 0)
	for _, u := range r.store {
		if searchFlag {
			if !strings.Contains(
				strings.ToLower(u.Email().String()),
				searchByEmail,
			) {
				continue
			}
		}
		if filterByStateFlag {
			contain := false
			for _, state := range filterByState {
				if u.State() == state {
					contain = true
				}
			}
			if !contain {
				continue
			}
		}
		if filterByStatusFlag {
			contain := false
			for _, status := range filterByStatus {
				if u.Status() == status {
					contain = true
				}
			}
			if !contain {
				continue
			}
		}
		result = append(result, u)
	}
	return result, nil
}

func (r *InMemoryUserRepository) Save(
	ctx context.Context,
	user *duser.User,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, err := duser.Restore(
		user.ID(),
		user.Email(),
		user.State(),
		user.Status(),
		user.PasswordHash(),
		user.ModifyVersion(),
	)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	r.store[user.ID()] = u
	return nil
}
