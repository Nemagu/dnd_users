package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/google/uuid"
)

type NewEmailUserRepository interface {
	ByID(ctx context.Context, id uuid.UUID) (*appdto.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user *appdto.User) error
}

type NewEmailDecrypter interface {
	DecryptEmailUserID(token string) (string, uuid.UUID, error)
}

type NewEmailUseCase struct {
	userRepo         NewEmailUserRepository
	emailDecrypter   NewEmailDecrypter
	emailValidator   EmailValidator
	passwordComparer PasswordComparer
}

func MustNewEmailUseCase(
	userRepo NewEmailUserRepository,
	emailDecrypter NewEmailDecrypter,
	emailValidator EmailValidator,
	passwordComparer PasswordComparer,
) *NewEmailUseCase {
	if userRepo == nil {
		panic("change email use case does not get user repository")
	}
	if emailDecrypter == nil {
		panic("change email use case does not get email decrypter")
	}
	if emailValidator == nil {
		panic("change email use case does not get email validator")
	}
	if passwordComparer == nil {
		panic("change email use case does not get password comparer")
	}
	return &NewEmailUseCase{
		userRepo:         userRepo,
		emailDecrypter:   emailDecrypter,
		emailValidator:   emailValidator,
		passwordComparer: passwordComparer,
	}
}

func (u *NewEmailUseCase) Execute(ctx context.Context, input *appdto.NewEmailCommand) error {
	email, userID, err := u.emailDecrypter.DecryptEmailUserID(input.Token)
	if err != nil {
		return err
	}

	if err := u.emailValidator.Validate(email); err != nil {
		return err
	}

	appUser, err := u.userRepo.ByID(ctx, userID)
	if err != nil {
		return err
	}

	compare, err := u.passwordComparer.Compare(input.Password, appUser.PasswordHash)
	if err != nil {
		return err
	}
	if !compare {
		return fmt.Errorf("%w: не верные учетные данные", application.ErrCredential)
	}

	domainUser, err := restoreDomainUser(appUser)
	if err != nil {
		return err
	}

	domainEmail, err := toDomainEmail(email)
	if err != nil {
		return err
	}

	if err = domainUser.ChangeEmail(domainEmail); err != nil {
		return handleError(err)
	}

	if err = u.userRepo.Save(ctx, toModifyAppUser(domainUser)); err != nil {
		return err
	}

	return nil
}
