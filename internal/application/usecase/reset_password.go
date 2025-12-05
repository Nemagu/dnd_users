package usecase

import (
	"context"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
)

type ResetPasswordUserRepository interface {
	ByEmail(ctx context.Context, email string) (*appdto.User, error)
	Save(ctx context.Context, user *appdto.User) error
}

type ResetPasswordUseCase struct {
	userRepo          ResetPasswordUserRepository
	emailDecrypter    EmailDecrypter
	emailValidator    EmailValidator
	passwordValidator PasswordValidator
	passwordHasher    PasswordHasher
}

func MustNewResetPasswordUseCase(
	userRepo ResetPasswordUserRepository,
	emailDecrypter EmailDecrypter,
	emailValidator EmailValidator,
	passwordValidator PasswordValidator,
	passwordHasher PasswordHasher,
) *ResetPasswordUseCase {
	if userRepo == nil {
		panic("reset password use case does not get user repository")
	}
	if emailDecrypter == nil {
		panic("reset password use case does not get email decrypter")
	}
	if emailValidator == nil {
		panic("reset password use case does not get email validator")
	}
	if passwordValidator == nil {
		panic("reset password use case does not get password validator")
	}
	if passwordHasher == nil {
		panic("reset password use case does not get password hasher")
	}
	return &ResetPasswordUseCase{
		userRepo:          userRepo,
		emailDecrypter:    emailDecrypter,
		emailValidator:    emailValidator,
		passwordValidator: passwordValidator,
		passwordHasher:    passwordHasher,
	}
}

func (u *ResetPasswordUseCase) Execute(
	ctx context.Context, input *appdto.ResetPasswordCommand,
) error {
	email, err := u.emailDecrypter.DecryptEmail(input.Token)
	if err != nil {
		return err
	}

	if err := u.emailValidator.Validate(email); err != nil {
		return err
	}

	appUser, err := u.userRepo.ByEmail(ctx, email)
	if err != nil {
		return err
	}

	if err := u.passwordValidator.Validate(input.NewPassword, email); err != nil {
		return err
	}

	hashedPassword, err := u.passwordHasher.Hash(input.NewPassword)
	if err != nil {
		return err
	}

	domainUser, err := restoreDomainUser(appUser)
	if err != nil {
		return err
	}

	if err := domainUser.ChangePassword(hashedPassword); err != nil {
		return handleError(err)
	}

	appUser = toModifyAppUser(domainUser)
	if err := u.userRepo.Save(ctx, appUser); err != nil {
		return err
	}

	return nil
}
