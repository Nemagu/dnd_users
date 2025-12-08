package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type newEmailRepository interface {
	ByID(ctx context.Context, id uuid.UUID) (*User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user *User) error
}

type newEmailCodeStore interface {
	GetNEC(ctx context.Context, key string) (string, error)
}

type NewEmailCommand struct {
	InitiatorID  uuid.UUID
	UserID       uuid.UUID
	NewEmail     string
	NewEmailCode string
	OldEmailCode string
	Password     string
}

type NewEmailUseCase struct {
	repo             newEmailRepository
	store            newEmailCodeStore
	emailValidator   emailValidator
	passwordComparer passwordComparer
}

func MustNewEmailUseCase(
	repo newEmailRepository,
	store newEmailCodeStore,
	emailValidator emailValidator,
	passwordComparer passwordComparer,
) *NewEmailUseCase {
	if repo == nil {
		panic("new email use case did not get user repository")
	}
	if store == nil {
		panic("new email use case did not get code store")
	}
	if emailValidator == nil {
		panic("new email use case did not get email validator")
	}
	if passwordComparer == nil {
		panic("new email use case did not get password comparer")
	}
	return &NewEmailUseCase{
		repo:             repo,
		store:            store,
		emailValidator:   emailValidator,
		passwordComparer: passwordComparer,
	}
}

func (u *NewEmailUseCase) Execute(ctx context.Context, command *NewEmailCommand) error {
	if command.InitiatorID != command.UserID {
		return fmt.Errorf("%w: вы не можете изменять email другим пользователям", ErrNotAllowed)
	}

	if err := u.emailValidator.Validate(command.NewEmail); err != nil {
		return err
	}

	exists, err := u.repo.EmailExists(ctx, command.NewEmail)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: email %s уже существует", ErrAlreadyExists, command.NewEmail)
	}

	appUser, err := u.repo.ByID(ctx, command.UserID)
	if err != nil {
		return err
	}

	compare, err := u.passwordComparer.Compare(command.Password, appUser.PasswordHash)
	if err != nil {
		return err
	}
	if !compare {
		return fmt.Errorf("%w: неверный пароль", ErrInvalidData)
	}

	oldEmailCode, err := u.store.GetNEC(ctx, appUser.Email)
	if err != nil {
		return err
	}
	newEmailCode, err := u.store.GetNEC(ctx, command.NewEmail)
	if err != nil {
		return err
	}

	if oldEmailCode != command.OldEmailCode {
		return fmt.Errorf("%w: не верный код для текущего email", ErrInvalidData)
	}
	if newEmailCode != command.NewEmailCode {
		return fmt.Errorf("%w: не верный код для нового email", ErrInvalidData)
	}

	domainUser, err := domainUser(appUser)
	if err != nil {
		return err
	}

	if err = domainUser.NewEmail(command.NewEmail); err != nil {
		return handleDomainError(err)
	}

	newAppUser, err := modifiedUser(domainUser)
	if err != nil {
		return err
	}
	if err = u.repo.Save(ctx, newAppUser); err != nil {
		return err
	}

	return nil
}
