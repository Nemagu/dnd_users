package inmemory

import (
	"fmt"
	"sync"

	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	store map[domain.UserID]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{}
}

func (r *InMemoryUserRepository) GetAll() []*domain.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.User, 0, len(r.store))
	for _, user := range r.store {
		result = append(result, user)
	}
	return result
}

func (r *InMemoryUserRepository) GetOfID(id domain.UserID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exist := r.store[id]
	if !exist {
		return user, application.NotFoundError(fmt.Sprintf("пользователь не существует (id: %s)", id))
	}
	return user, nil
}

func (r *InMemoryUserRepository) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[user.UserID()] = user
	return nil
}

type emailChanged struct {
	Email     domain.Email
	Published bool
}

type stateChanged struct {
	State     domain.UserState
	Published bool
}

type statusChanged struct {
	Status    domain.UserStatus
	Published bool
}

type InMemoryEventRepository struct {
	mu          sync.RWMutex
	emailStore  map[domain.UserID][]emailChanged
	stateStore  map[domain.UserID][]stateChanged
	statusStore map[domain.UserID][]statusChanged
}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{}
}

func (r *InMemoryEventRepository) EmailChanged(
	id domain.UserID,
	email domain.Email,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.emailStore[id] = append(r.emailStore[id], emailChanged{Email: email, Published: false})
	return nil
}

func (r *InMemoryEventRepository) StateChanged(
	id domain.UserID,
	state domain.UserState,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stateStore[id] = append(r.stateStore[id], stateChanged{State: state, Published: false})
	return nil
}

func (r *InMemoryEventRepository) StatusChanged(
	id domain.UserID,
	status domain.UserStatus,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.statusStore[id] = append(r.statusStore[id], statusChanged{Status: status, Published: false})
	return nil
}
