package command

import (
	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
)

type UserFreezeCommand struct {
	userRepository  application.UserRepository
	eventRepository application.EventRepository
}

func NewUserFreezeCommand(
	userRepository application.UserRepository,
	eventRepository application.EventRepository,
) *UserFreezeCommand {
	return &UserFreezeCommand{
		userRepository:  userRepository,
		eventRepository: eventRepository,
	}
}

func (c *UserFreezeCommand) Execute(initiatorID, userID domain.UserID) error {
	initiator, err := c.userRepository.GetOfID(initiatorID)
	if err != nil {
		return err
	}
	if !initiator.Status().IsAdmin() {
		return application.NoAccessError("")
	}
	if initiatorID == userID {
		return application.InvalidDataError("вы не можете заморозить сами себя")
	}
	user, err := c.userRepository.GetOfID(userID)
	if err != nil {
		return err
	}
	err = user.Freeze()
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
