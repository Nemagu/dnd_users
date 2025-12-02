package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/repository"
	"github.com/Nemagu/dnd/internal/application/service"
	"github.com/Nemagu/dnd/internal/domain"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

type RegisterUserUseCase struct {
	userRepo          repository.UserRepository
	passwordValidator service.PasswordValidator
	passwordHasher    service.PasswordHasher
	emailCrypter      service.EmailCrypter
	emailValidator    service.EmailValidator
}

func NewRegisterUserUseCase(
	userRepo repository.UserRepository,
	passwordValidator service.PasswordValidator,
	passwordHasher service.PasswordHasher,
	emailCrypter service.EmailCrypter,
) (*RegisterUserUseCase, error) {
	return &RegisterUserUseCase{
		userRepo:          userRepo,
		passwordValidator: passwordValidator,
		passwordHasher:    passwordHasher,
		emailCrypter:      emailCrypter,
	}, nil
}

func MustNewRegisterUserUseCase(
	userRepo repository.UserRepository,
	passwordValidator service.PasswordValidator,
	passwordHasher service.PasswordHasher,
	emailCrypter service.EmailCrypter,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:          userRepo,
		passwordValidator: passwordValidator,
		passwordHasher:    passwordHasher,
		emailCrypter:      emailCrypter,
	}
}

func (u *RegisterUserUseCase) Execute(
	ctx context.Context,
	input appdto.RegisterUserCommand,
) (uuid.UUID, error) {
	email, err := u.emailCrypter.Decrypt(input.Token)
	if err != nil {
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

	domainEmail, err := domain.NewEmail(email)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%w: %s", application.ErrValidation, err)
	}

	domainUser, err := duser.New(
		u.userRepo.NextID(ctx),
		domainEmail,
		passwordHash,
	)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%w: %s", application.ErrValidation, err)
	}

	dtoUser := appdto.User{
		UserID:       domainUser.ID(),
		Email:        domainUser.Email().String(),
		State:        domainUser.State().String(),
		Status:       domainUser.Status().String(),
		PasswordHash: domainUser.PasswordHash(),
		Version:      domainUser.ModifyVersion(),
	}

	if err := u.userRepo.Save(ctx, &dtoUser); err != nil {
		return uuid.UUID{}, err
	}

	return dtoUser.UserID, nil
}
