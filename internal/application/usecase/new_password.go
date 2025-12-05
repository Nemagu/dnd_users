package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/google/uuid"
)

type NewPasswordUserRepository interface {
	ByID(ctx context.Context, id uuid.UUID) (*appdto.User, error)
	Save(ctx context.Context, user *appdto.User) error
}

type NewPasswordUseCase struct {
	userRepo          NewPasswordUserRepository
	passwordValidator PasswordValidator
	passwordComparer  PasswordComparer
	passwordHasher    PasswordHasher
}

func MustNewPasswordUseCase(
	userRepo NewPasswordUserRepository,
	passwordValidator PasswordValidator,
	passwordComparer PasswordComparer,
	passwordHasher PasswordHasher,
) *NewPasswordUseCase {
	if userRepo == nil {
		panic("change password use case does not get user repository")
	}
	if passwordValidator == nil {
		panic("change password use case does not get password validator")
	}
	if passwordComparer == nil {
		panic("change password use case does not get password comparer")
	}
	return &NewPasswordUseCase{
		userRepo:          userRepo,
		passwordValidator: passwordValidator,
		passwordComparer:  passwordComparer,
		passwordHasher:    passwordHasher,
	}
}

func (u *NewPasswordUseCase) Execute(
	ctx context.Context, input *appdto.ChangePasswordCommand,
) error {
	appUser, err := u.userRepo.ByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	domainUser, err := restoreDomainUser(appUser)
	if err != nil {
		return err
	}

	compare, err := u.passwordComparer.Compare(input.OldPassword, domainUser.PasswordHash())
	if err != nil {
		return err
	}

	if !compare {
		return fmt.Errorf("%w: не верный пароль", application.ErrValidation)
	}

	if err = u.passwordValidator.Validate(
		input.NewPassword, domainUser.Email().String(),
	); err != nil {
		return err
	}

	hashedPassword, err := u.passwordHasher.Hash(input.NewPassword)
	if err != nil {
		return err
	}

	if err = domainUser.ChangePassword(hashedPassword); err != nil {
		return handleError(err)
	}

	if err = u.userRepo.Save(ctx, toModifyAppUser(domainUser)); err != nil {
		return err
	}

	return nil
}
