package command

import (
	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
)

type AppointAdminCommand struct {
	userRepository  application.UserRepository
	eventRepository application.EventRepository
}

func NewAppointAdminCommand(
	userRepository application.UserRepository,
	eventRepository application.EventRepository,
) *AppointAdminCommand {
	return &AppointAdminCommand{
		userRepository:  userRepository,
		eventRepository: eventRepository,
	}
}

func (c *AppointAdminCommand) Execute(initiatorID, userID domain.UserID) error {
	initiator, err := c.userRepository.GetOfID(initiatorID)
	if err != nil {
		return err
	}
	if !initiator.Status().IsAdmin() {
		return application.NoAccessError("")
	}
	user, err := c.userRepository.GetOfID(userID)
	if err != nil {
		return err
	}
	err = user.AppointAdmin()
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
