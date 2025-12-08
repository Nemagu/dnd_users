package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type NewPasswordUseCase struct {
	repo              newPasswordRepository
	store             newPasswordCodeStore
	passwordComparer  passwordComparer
	passwordValidator passwordValidator
	passwordHasher    passwordHasher
}

type NewPasswordCommand struct {
	InitiatorID uuid.UUID
	UserID      uuid.UUID
	OldPassword string
	NewPassword string
	Code        string
}

type newPasswordRepository interface {
	ByID(ctx context.Context, id uuid.UUID) (*User, error)
	Save(ctx context.Context, user *User) error
}

type newPasswordCodeStore interface {
	GetNewPassword(ctx context.Context, key string) (string, error)
	DelNewPassword(ctx context.Context, key string) error
}

func MustNewPasswordUseCase(
	repo newPasswordRepository,
	store newPasswordCodeStore,
	passwordComparer passwordComparer,
	passwordValidator passwordValidator,
	passwordHasher passwordHasher,
) *NewPasswordUseCase {
	if repo == nil {
		panic("new password use case did not get user repository")
	}
	if store == nil {
		panic("new password use case did not get code store")
	}
	if passwordComparer == nil {
		panic("new password use case did not get password comparer")
	}
	if passwordValidator == nil {
		panic("new password use case did not get password validator")
	}
	if passwordHasher == nil {
		panic("new password use case did not get password hasher")
	}
	return &NewPasswordUseCase{
		repo:              repo,
		store:             store,
		passwordComparer:  passwordComparer,
		passwordValidator: passwordValidator,
		passwordHasher:    passwordHasher,
	}
}

func (u *NewPasswordUseCase) Execute(ctx context.Context, command *NewPasswordCommand) error {
	if command.InitiatorID != command.UserID {
		return fmt.Errorf(
			"%w: вы не можете изменять пароль другим пользователям",
			ErrNotAllowed,
		)
	}

	if command.NewPassword == command.OldPassword {
		return fmt.Errorf(
			"%w: новый пароль не должен совпадать со старым",
			ErrInvalidData,
		)
	}

	appUser, err := u.repo.ByID(ctx, command.UserID)
	if err != nil {
		return err
	}

	code, err := u.store.GetNewPassword(ctx, appUser.ID.String()+appUser.Email)
	if err != nil {
		return err
	}
	if code != command.Code {
		return fmt.Errorf("%w: неверный код для смены пароля", ErrInvalidData)
	}

	compare, err := u.passwordComparer.Compare(command.OldPassword, appUser.PasswordHash)
	if err != nil {
		return err
	}
	if !compare {
		return fmt.Errorf("%w: неверный пароль", ErrInvalidData)
	}

	if err = u.passwordValidator.Validate(command.NewPassword, appUser.Email); err != nil {
		return err
	}

	hashedPassword, err := u.passwordHasher.Hash(command.NewPassword)
	if err != nil {
		return err
	}

	domainUser, err := domainUser(appUser)
	if err != nil {
		return err
	}

	if err = domainUser.NewPasswordHash(hashedPassword); err != nil {
		return handleDomainError(err)
	}

	newAppUser, err := modifiedUser(domainUser)
	if err != nil {
		return err
	}

	if err = u.store.DelNewPassword(ctx, appUser.ID.String()+appUser.Email); err != nil {
		return err
	}

	if err = u.repo.Save(ctx, newAppUser); err != nil {
		return err
	}

	return nil
}
