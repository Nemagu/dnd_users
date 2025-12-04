package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

type ChangeUserRepository interface {
	ByID(ctx context.Context, id uuid.UUID) (*appdto.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user *appdto.User) error
}

type ChangeUserUseCase struct {
	userRepo          ChangeUserRepository
	policyService     *duser.PolicyService
	passwordHasher    PasswordHasher
	passwordValidator PasswordValidator
}

func MustNewChangeUserUseCase(
	userRepo ChangeUserRepository,
	policyService *duser.PolicyService,
	passwordHasher PasswordHasher,
	passwordValidator PasswordValidator,
) *ChangeUserUseCase {
	if userRepo == nil {
		panic("change user use case does not get user repository")
	}
	if policyService == nil {
		panic("change user use case does not get policy service")
	}
	if passwordHasher == nil {
		panic("change user use case does not get password hasher")
	}
	if passwordValidator == nil {
		panic("change user use case does not get password validator")
	}
	return &ChangeUserUseCase{
		userRepo:          userRepo,
		policyService:     policyService,
		passwordHasher:    passwordHasher,
		passwordValidator: passwordValidator,
	}
}

func (u *ChangeUserUseCase) Execute(ctx context.Context, input *appdto.ChangeUserCommand) error {
	if input.InitiatorID == input.UserID {
		return fmt.Errorf("%w: вы не можете редактировать сами себя", application.ErrNotAllowed)
	}

	appUser, err := u.userRepo.ByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	appInitiator, err := u.userRepo.ByID(ctx, input.InitiatorID)
	if err != nil {
		return err
	}

	domainInitiator, err := restoreDomainUser(appInitiator)
	if err != nil {
		return err
	}
	if !u.policyService.CanEditOther(domainInitiator) {
		return fmt.Errorf(
			"%w: вы не можете редактировать других пользователей", application.ErrNotAllowed,
		)
	}

	domainUser, err := restoreDomainUser(appUser)
	if err != nil {
		return err
	}

	if err := u.change(ctx, domainUser, input); err != nil {
		return err
	}

	if err := u.userRepo.Save(ctx, toModifyAppUser(domainUser)); err != nil {
		return err
	}

	return nil
}

func (u *ChangeUserUseCase) change(
	ctx context.Context, user *duser.User, input *appdto.ChangeUserCommand,
) error {
	if err := u.changeEmail(ctx, user, input.Email); err != nil {
		return err
	}

	if err := u.changeState(user, input.State); err != nil {
		return err
	}

	if err := u.changeStatus(user, input.Status); err != nil {
		return err
	}

	if err := u.changePassword(user, input.Password); err != nil {
		return err
	}

	return nil
}

func (u *ChangeUserUseCase) changeEmail(
	ctx context.Context, user *duser.User, email string,
) error {
	if email == "" {
		return nil
	}

	exists, err := u.userRepo.EmailExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf(
			"%w: пользователь с таким email уже существует", application.ErrValidation,
		)
	}

	e, err := toDomainEmail(email)
	if err != nil {
		return err
	}

	if err := user.ChangeEmail(e); err != nil {
		return handleError(err)
	}

	return nil
}

func (u *ChangeUserUseCase) changeState(user *duser.User, state string) error {
	if state == "" {
		return nil
	}

	s, err := toDomainState(state)
	if err != nil {
		return err
	}

	if err := user.ChangeState(s); err != nil {
		return handleError(err)
	}

	return nil
}

func (u *ChangeUserUseCase) changeStatus(user *duser.User, status string) error {
	if status == "" {
		return nil
	}

	s, err := toDomainStatus(status)
	if err != nil {
		return err
	}

	if err := user.ChangeStatus(s); err != nil {
		return handleError(err)
	}

	return nil
}

func (u *ChangeUserUseCase) changePassword(user *duser.User, password string) error {
	if password == "" {
		return nil
	}

	if err := u.passwordValidator.Validate(password, user.Email().String()); err != nil {
		return err
	}

	passHash, err := u.passwordHasher.Hash(password)
	if err != nil {
		return err
	}

	if err := user.ChangePassword(passHash); err != nil {
		return handleError(err)
	}

	return nil
}
