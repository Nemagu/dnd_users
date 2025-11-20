package repository

import (
	"context"

	"github.com/Nemagu/dnd/internal/domain"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

type UserRepository interface {
	NextID(ctx context.Context) uuid.UUID
	IDExists(
		ctx context.Context,
		id uuid.UUID,
	) (bool, error)
	EmailExists(
		ctx context.Context,
		email domain.Email,
	) (bool, error)
	All(
		ctx context.Context,
		limit, offset int,
	) []*duser.User
	ByID(
		ctx context.Context,
		id uuid.UUID,
	) (*duser.User, error)
	ByEmail(
		ctx context.Context,
		email domain.Email,
	) (*duser.User, error)
	Filter(
		ctx context.Context,
		searchByEmail string,
		filterByState []duser.UserState,
		filterByStatus []duser.UserStatus,
		limit, offset int,
	) ([]*duser.User, error)
	Save(
		ctx context.Context,
		user *duser.User,
	) error
}
