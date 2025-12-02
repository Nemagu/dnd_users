package inmemory

import (
	"context"
	"fmt"
	"strings"

	"sync"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/google/uuid"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	store map[uuid.UUID]*appdto.User
}

func MustNewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		store: make(map[uuid.UUID]*appdto.User),
	}
}

func (r *InMemoryUserRepository) NextID(ctx context.Context) uuid.UUID {
	return uuid.New()
}

func (r *InMemoryUserRepository) IDExists(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.store[id]
	return exists, nil
}

func (r *InMemoryUserRepository) EmailExists(
	ctx context.Context,
	email string,
) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.store {
		if u.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (r *InMemoryUserRepository) All(
	ctx context.Context,
	limit, offset int,
) ([]*appdto.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*appdto.User, 0, len(r.store))
	for _, u := range r.store {
		result = append(result, u)
	}
	return result, nil
}

func (r *InMemoryUserRepository) ByID(
	ctx context.Context,
	id uuid.UUID,
) (*appdto.User, error) {
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

func (r *InMemoryUserRepository) ByEmail(
	ctx context.Context,
	email string,
) (*appdto.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.store {
		if u.Email == email {
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
	filterByState []string,
	filterByStatus []string,
	limit, offset int,
) ([]*appdto.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	searchFlag := len(searchByEmail) > 0
	filterByStateFlag := len(filterByState) > 0
	filterByStatusFlag := len(filterByStatus) > 0
	result := make([]*appdto.User, 0)
	for _, u := range r.store {
		if searchFlag {
			if !strings.Contains(
				strings.ToLower(u.Email),
				searchByEmail,
			) {
				continue
			}
		}
		if filterByStateFlag {
			contain := false
			for _, state := range filterByState {
				if u.State == state {
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
				if u.Status == status {
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
	user *appdto.User,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[user.UserID] = user
	return nil
}
