package command

import (
	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
)

type UsernameChangeCommand struct {
	userRepository  application.UserRepository
	eventRepository application.EventRepository
}

func NewUsernameChangeCommand(
	userRepository application.UserRepository,
	eventRepository application.EventRepository,
) *UsernameChangeCommand {
	return &UsernameChangeCommand{
		userRepository:  userRepository,
		eventRepository: eventRepository,
	}
}

func (c *UsernameChangeCommand) Execute(
	initiatorID, userID domain.UserID,
	username domain.Username,
) error {
	initiator, err := c.userRepository.GetOfID(initiatorID)
	if err != nil {
		return err
	}
	if !initiator.Status().IsAdmin() && initiatorID != userID {
		return application.NoAccessError("")
	}
	var user *domain.User
	if initiatorID == userID {
		user = initiator
	} else {
		user, err = c.userRepository.GetOfID(userID)
		if err != nil {
			return err
		}
	}
	err = user.ChangeUsername(username)
	if err != nil {
		return application.InvalidDataError(err.Error())
	}
	err = c.userRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}
