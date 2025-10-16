package command

import (
	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
)

type AppointOrdinaryCommand struct {
	userRepository  application.UserRepository
	eventRepository application.EventRepository
}

func NewAppointOrdinaryCommand(
	userRepository application.UserRepository,
	eventRepository application.EventRepository,
) *AppointOrdinaryCommand {
	return &AppointOrdinaryCommand{
		userRepository:  userRepository,
		eventRepository: eventRepository,
	}
}

func (c *AppointOrdinaryCommand) Execute(initiatorID, userID domain.UserID) error {
	initiator, err := c.userRepository.GetOfID(initiatorID)
	if err != nil {
		return err
	}
	if !initiator.Status().IsAdmin() {
		return application.NoAccessError("")
	}
	if initiatorID == userID {
		return application.InvalidDataError("вы не можете сами себя сделать обычным пользователем")
	}
	user, err := c.userRepository.GetOfID(userID)
	if err != nil {
		return err
	}
	err = user.AppointOrdinary()
	if err != nil {
		return application.InvalidDataError(err.Error())
	}
	err = c.userRepository.Save(user)
	if err != nil {
		return err
	}
	err = c.eventRepository.StatusChanged(user.UserID(), user.Status())
	if err != nil {
		return err
	}
	return nil
}
