package app

import (
	"context"
	"fmt"
)

type ResetPasswordUseCase struct {
	repo              resetPasswordRepository
	store             resetPasswordCodeStore
	passwordValidator passwordValidator
	passwordHasher    passwordHasher
}

type ResetPasswordCommand struct {
	Email       string
	Code        string
	NewPassword string
}

type resetPasswordRepository interface {
	ByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, user *User) error
}

type resetPasswordCodeStore interface {
	GetResetPassword(ctx context.Context, key string) (string, error)
	DelResetPassword(ctx context.Context, key string) error
}

func MustResetPasswordUseCase(
	repo resetPasswordRepository,
	store resetPasswordCodeStore,
	passwordValidator passwordValidator,
	passwordHasher passwordHasher,
) *ResetPasswordUseCase {
	if repo == nil {
		panic("reset password use case did not get user repository")
	}
	if store == nil {
		panic("reset password use case did not get code store")
	}
	if passwordValidator == nil {
		panic("reset password use case did not get password validator")
	}
	if passwordHasher == nil {
		panic("reset password use case did not get password hasher")
	}
	return &ResetPasswordUseCase{
		repo:              repo,
		store:             store,
		passwordValidator: passwordValidator,
		passwordHasher:    passwordHasher,
	}
}

func (u *ResetPasswordUseCase) Execute(ctx context.Context, command *ResetPasswordCommand) error {
	user, err := u.repo.ByEmail(ctx, command.Email)
	if err != nil {
		return err
	}

	if err = u.passwordValidator.Validate(command.NewPassword, user.Email); err != nil {
		return err
	}

	code, err := u.store.GetResetPassword(ctx, user.ID.String()+user.Email)
	if err != nil {
		return err
	}
	if code != command.Code {
		return fmt.Errorf("%w: не верный код сброса пароля", ErrInvalidData)
	}

	hashedPassword, err := u.passwordHasher.Hash(command.NewPassword)
	if err != nil {
		return err
	}

	domainUser, err := domainUser(user)
	if err != nil {
		return handleDomainError(err)
	}

	if err = domainUser.NewPasswordHash(hashedPassword); err != nil {
		return handleDomainError(err)
	}

	if err = u.store.DelResetPassword(ctx, user.ID.String()+user.Email); err != nil {
		return err
	}

	newUser, err := modifiedUser(domainUser)
	if err != nil {
		return err
	}
	if err = u.repo.Save(ctx, newUser); err != nil {
		return err
	}

	return nil
}
