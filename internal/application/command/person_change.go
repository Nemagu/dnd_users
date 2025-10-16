package command

import (
	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
)

type PersonChangeCommand struct {
	userRepository  application.UserRepository
	eventRepository application.EventRepository
}

func NewPersonChangeCommand(
	userRepository application.UserRepository,
	eventRepository application.EventRepository,
) *PersonChangeCommand {
	return &PersonChangeCommand{
		userRepository:  userRepository,
		eventRepository: eventRepository,
	}
}

func (c *PersonChangeCommand) Execute(
	initiatorID, userID domain.UserID,
	person domain.Person,
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
	err = user.ChangePerson(person)
	if err != nil {
		return application.InvalidDataError(err.Error())
	}
	err = c.userRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}
