package repository

import (
	"context"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
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
		email string,
	) (bool, error)
	All(
		ctx context.Context,
		limit, offset int,
	) ([]*appdto.User, error)
	ByID(
		ctx context.Context,
		id uuid.UUID,
	) (*appdto.User, error)
	ByEmail(
		ctx context.Context,
		email string,
	) (*appdto.User, error)
	Filter(
		ctx context.Context,
		searchByEmail string,
		filterByState []string,
		filterByStatus []string,
		limit, offset int,
	) ([]*appdto.User, error)
	Save(
		ctx context.Context,
		user *appdto.User,
	) error
}
