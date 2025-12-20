package app

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd_users/internal/domain"
	"github.com/google/uuid"
)

type ChangeUserUseCase struct {
	repo              changeUserRepository
	emailValidator    emailValidator
	passwordValidator passwordValidator
	passwordHasher    passwordHasher
	policy            *domain.PolicyService
}

type ChangeUserCommand struct {
	InitiatorID uuid.UUID
	UserID      uuid.UUID
	Email       string
	State       string
	Status      string
	Password    string
}

type changeUserRepository interface {
	IDExists(ctx context.Context, id uuid.UUID) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	ByID(ctx context.Context, id uuid.UUID) (*User, error)
	Save(ctx context.Context, user *User) error
}

func MustChangeUserUseCase(
	repo changeUserRepository,
	emailValidator emailValidator,
	passwordValidator passwordValidator,
	passwordHasher passwordHasher,
	policy *domain.PolicyService,
) *ChangeUserUseCase {
	if repo == nil {
		panic("change user use case did not get user repository")
	}
	if emailValidator == nil {
		panic("change user use case did not get email validator")
	}
	if passwordValidator == nil {
		panic("change user use case did not get password validator")
	}
	if passwordHasher == nil {
		panic("change user use case did not get password hasher")
	}
	if policy == nil {
		panic("change user use case did not get policy service")
	}
	return &ChangeUserUseCase{
		repo:              repo,
		emailValidator:    emailValidator,
		passwordValidator: passwordValidator,
		passwordHasher:    passwordHasher,
		policy:            policy,
	}
}

func (u *ChangeUserUseCase) Execute(ctx context.Context, command *ChangeUserCommand) error {
	exists, err := u.repo.IDExists(ctx, command.InitiatorID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: пользователь с id %s не найден", ErrNotFound, command.InitiatorID)
	}
	initiator, err := u.repo.ByID(ctx, command.InitiatorID)
	if err != nil {
		return err
	}

	domainInitiator, err := domainUser(initiator)
	if err != nil {
		return err
	}
	if !u.policy.CanEditOthers(domainInitiator) {
		return fmt.Errorf("%w: вы не можете редактировать других пользователей", ErrNotAllowed)
	}

	exists, err = u.repo.IDExists(ctx, command.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: пользователь с id %s не найден", ErrNotFound, command.UserID)
	}

	user, err := u.repo.ByID(ctx, command.UserID)
	if err != nil {
		return err
	}

	domainUser, err := domainUser(user)
	if err != nil {
		return err
	}
	if command.Email != "" {
		if err = u.emailValidator.Validate(command.Email); err != nil {
			return err
		}
		exists, err := u.repo.EmailExists(ctx, command.Email)
		if err != nil {
			return err
		}
		if exists {
			if command.Email != domainUser.Email() {
				return fmt.Errorf("%w: email %s уже существует", ErrInvalidData, command.Email)
			}
		}
		if err = domainUser.NewEmail(command.Email); err != nil {
			return handleDomainError(err)
		}
	}

	if command.State != "" {
		state, err := domainState(command.State)
		if err != nil {
			return err
		}
		if err = domainUser.NewState(state); err != nil {
			return handleDomainError(err)
		}
	}

	if command.Status != "" {
		status, err := domainStatus(command.Status)
		if err != nil {
			return err
		}
		if err = domainUser.NewStatus(status); err != nil {
			return handleDomainError(err)
		}
	}

	if command.Password != "" {
		if err = u.passwordValidator.Validate(command.Password, domainUser.Email()); err != nil {
			return err
		}
		hashedPassword, err := u.passwordHasher.Hash(command.Password)
		if err != nil {
			return err
		}
		if err = domainUser.NewPasswordHash(hashedPassword); err != nil {
			return handleDomainError(err)
		}
	}

	user, err = modifiedUser(domainUser)
	if err != nil {
		return err
	}
	if err = u.repo.Save(ctx, user); err != nil {
		return err
	}

	return nil
}
