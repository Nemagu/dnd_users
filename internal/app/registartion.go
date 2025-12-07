package app

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd_users/internal/domain"
	"github.com/google/uuid"
)

type registrationRepository interface {
	NextID(ctx context.Context) (uuid.UUID, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user *User) error
}

type registrationCodeStore interface {
	GetCEC(ctx context.Context, key string) (string, error)
	DelCEC(ctx context.Context, key string) error
}

type RegistrationCommand struct {
	Email    string
	Password string
	Code     string
}

type RegistrationUseCase struct {
	repo              registrationRepository
	store             registrationCodeStore
	emailValidator    emailValidator
	passwordValidator passwordValidator
	passwordHasher    passwordHasher
}

func MustRegistrationUseCase(
	repo registrationRepository,
	store registrationCodeStore,
	emailValidator emailValidator,
	passwordValidator passwordValidator,
	passwordHasher passwordHasher,
) *RegistrationUseCase {
	if repo == nil {
		panic("registration use case did not get user repository")
	}
	if store == nil {
		panic("registration use case did not get code store")
	}
	if emailValidator == nil {
		panic("registration use case did not get email validator")
	}
	if passwordValidator == nil {
		panic("registration use case did not get password validator")
	}
	if passwordHasher == nil {
		panic("registration use case did not get password hasher")
	}
	return &RegistrationUseCase{
		repo:              repo,
		store:             store,
		emailValidator:    emailValidator,
		passwordValidator: passwordValidator,
		passwordHasher:    passwordHasher,
	}
}

func (u *RegistrationUseCase) Execute(
	ctx context.Context,
	command *RegistrationCommand,
) (uuid.UUID, error) {
	code, err := u.store.GetCEC(ctx, command.Code)
	if err != nil {
		return uuid.Nil, err
	}

	if code != command.Code {
		return uuid.Nil, fmt.Errorf("%w: не верный код подтверждения", ErrInvalidData)
	}

	if err = u.emailValidator.Validate(command.Email); err != nil {
		return uuid.Nil, err
	}

	exists, err := u.repo.EmailExists(ctx, command.Email)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, fmt.Errorf(
			"%w: пользователь с email %s уже существует",
			ErrInvalidData,
			command.Email,
		)
	}

	if err = u.passwordValidator.Validate(command.Password, command.Email); err != nil {
		return uuid.Nil, err
	}

	hashedPassword, err := u.passwordHasher.Hash(command.Password)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := u.repo.NextID(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	domainUser, err := domain.NewUser(id, command.Email, hashedPassword)
	if err != nil {
		return uuid.Nil, handleDomainError(err)
	}

	appUser, err := modifiedUser(domainUser)
	if err != nil {
		return uuid.Nil, err
	}

	if err = u.repo.Save(ctx, appUser); err != nil {
		return uuid.Nil, err
	}

	if err = u.store.DelCEC(ctx, command.Code); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
