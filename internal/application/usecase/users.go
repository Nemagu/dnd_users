package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

type UsersRepository interface {
	ByID(ctx context.Context, id uuid.UUID) (*appdto.User, error)
	Filter(
		ctx context.Context,
		searchByEmail string,
		filterByState []string,
		filterByStatus []string,
		limit, offset int,
	) ([]*appdto.User, error)
}

type UsersUseCase struct {
	userRepo      UsersRepository
	policyService *duser.PolicyService
}

func MustNewUsersUseCase(
	userRepo UsersRepository, policyService *duser.PolicyService,
) *UsersUseCase {
	if userRepo == nil {
		panic("users use case does not get user repository")
	}
	if policyService == nil {
		panic("users use case does not get policy service")
	}
	return &UsersUseCase{
		userRepo:      userRepo,
		policyService: policyService,
	}
}

func (u *UsersUseCase) Execute(ctx context.Context, input *appdto.UsersQuery) ([]*appdto.User, error) {
	appInitiator, err := u.userRepo.ByID(ctx, input.InitiatorID)
	if err != nil {
		return nil, err
	}

	domainInitiator, err := restoreDomainUser(appInitiator)
	if err != nil {
		return nil, err
	}

	if !u.policyService.CanReadAll(domainInitiator) {
		return nil, fmt.Errorf(
			"%w: вы не можете просматривать других пользователей",
			application.ErrNotAllowed,
		)
	}

	return u.userRepo.Filter(
		ctx,
		input.SearchByEmail,
		input.FilterByState,
		input.FilterByStatus,
		input.Limit,
		input.Offset,
	)
}
