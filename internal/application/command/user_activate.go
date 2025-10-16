package command

import (
	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
)

type UserActivateCommand struct {
	userRepository  application.UserRepository
	eventRepository application.EventRepository
}

func NewUserActivateCommand(
	userRepository application.UserRepository,
	eventRepository application.EventRepository,
) *UserActivateCommand {
	return &UserActivateCommand{
		userRepository:  userRepository,
		eventRepository: eventRepository,
	}
}

func (c *UserActivateCommand) Execute(initiatorID, userID domain.UserID) error {
	initiator, err := c.userRepository.GetOfID(initiatorID)
	if err != nil {
		return err
	}
	if !initiator.Status().IsAdmin() {
		return application.NoAccessError("")
	}
	if initiatorID == userID {
		return application.InvalidDataError("вы не можете активировать сами себя")
	}
	user, err := c.userRepository.GetOfID(userID)
	if err != nil {
		return err
	}
	err = user.Activate()
	if err != nil {
		return application.InvalidDataError(err.Error())
	}
	err = c.userRepository.Save(user)
	if err != nil {
		return err
	}
	err = c.eventRepository.StateChanged(user.UserID(), user.State())
	if err != nil {
		return err
	}
	return nil
}
