package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/google/uuid"
)

type RegisterUserRepository interface {
	NextID(ctx context.Context) uuid.UUID
	Save(ctx context.Context, user *appdto.User) error
}

type EmailValidator interface {
	Validate(email string) error
}

type PasswordValidator interface {
	Validate(password string, email string) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type EmailDecrypter interface {
	Decrypt(token string) (string, error)
}

type RegisterUserUseCase struct {
	userRepo          RegisterUserRepository
	passwordValidator PasswordValidator
	passwordHasher    PasswordHasher
	emailDecrypter    EmailDecrypter
	emailValidator    EmailValidator
}

func NewRegisterUserUseCase(
	userRepo RegisterUserRepository,
	passwordValidator PasswordValidator,
	passwordHasher PasswordHasher,
	emailDecrypter EmailDecrypter,
	emailValidator EmailValidator,
) (*RegisterUserUseCase, error) {
	return &RegisterUserUseCase{
		userRepo:          userRepo,
		passwordValidator: passwordValidator,
		passwordHasher:    passwordHasher,
		emailDecrypter:    emailDecrypter,
		emailValidator:    emailValidator,
	}, nil
}

func MustNewRegisterUserUseCase(
	userRepo RegisterUserRepository,
	passwordValidator PasswordValidator,
	passwordHasher PasswordHasher,
	emailDecrypter EmailDecrypter,
	emailValidator EmailValidator,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:          userRepo,
		passwordValidator: passwordValidator,
		passwordHasher:    passwordHasher,
		emailDecrypter:    emailDecrypter,
		emailValidator:    emailValidator,
	}
}

func (u *RegisterUserUseCase) Execute(
	ctx context.Context,
	input *appdto.RegisterUserCommand,
) (uuid.UUID, error) {
	email, err := u.emailDecrypter.Decrypt(input.Token)
	if err != nil {
		return uuid.UUID{}, err
	}

	if err := u.emailValidator.Validate(email); err != nil {
		return uuid.UUID{}, err
	}

	err = u.passwordValidator.Validate(input.Password, email)
	if err != nil {
		return uuid.UUID{}, err
	}

	passwordHash, err := u.passwordHasher.Hash(input.Password)
	if err != nil {
		return uuid.UUID{}, err
	}

	domainUser, err := newDomainUser(u.userRepo.NextID(ctx), email, passwordHash)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%w: %s", application.ErrValidation, err)
	}

	dtoUser := toAppUser(domainUser)

	if err := u.userRepo.Save(ctx, dtoUser); err != nil {
		return uuid.UUID{}, err
	}

	return dtoUser.UserID, nil
}
