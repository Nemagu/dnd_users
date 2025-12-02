package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

type UserRepository interface {
	ByID(ctx context.Context, id uuid.UUID) (*appdto.User, error)
}

type UserUseCase struct {
	userRepo      UserRepository
	policyService *duser.PolicyService
}

func MustNewUserUseCase(
	userRepo UserRepository, policyService *duser.PolicyService,
) *UserUseCase {
	if userRepo == nil {
		panic("user use case does not get user repository")
	}
	if policyService == nil {
		panic("user use case does not get policy service")
	}
	return &UserUseCase{
		userRepo:      userRepo,
		policyService: policyService,
	}
}

func (u *UserUseCase) Execute(
	ctx context.Context, input *appdto.UserQuery,
) (*appdto.User, error) {
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

	return u.userRepo.ByID(ctx, input.UserID)
}
