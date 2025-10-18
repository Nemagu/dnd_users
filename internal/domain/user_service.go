package domain

type UserRepository interface {
	GetAll() []*User
	GetOfID(id UserID) (*User, error)
	IsExistsEmail(email Email) bool
	Save(user *User) error
}

type UserService struct {
	userRepository UserRepository
	eventPublisher *EventPublisher
}

func NewUserService(
	userRepository UserRepository,
	eventPublisher *EventPublisher,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		eventPublisher: eventPublisher,
	}
}

func (s *UserService) ChangeEmail(
	initiator, user *User,
	email Email,
) error {
	err := s.assertAccessOfChange(initiator, user)
	if err != nil {
		return err
	}
	if s.userRepository.IsExistsEmail(email) {
		return InvalidDataError("такой email уже существует")
	}
	err = user.ChangeEmail(email)
	if err != nil {
		return err
	}
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}
	err = s.eventPublisher.Publish(NewEmailChangedFromUser(user))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ChangePassword(
	initiator, user *User,
	passwordHash PasswordHash,
) error {
	err := s.assertAccessOfChange(initiator, user)
	if err != nil {
		return err
	}
	err = user.ChangePassword(passwordHash)
	if err != nil {
		return err
	}
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) AppointAdmin(initiator, user *User) error {
	err := s.assertNotIsAdmin(initiator)
	if err != nil {
		return err
	}
	err = user.AppointAdmin()
	if err != nil {
		return err
	}
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) AppointOrdinary(initiator, user *User) error {
	err := s.assertNotIsAdmin(initiator)
	if err != nil {
		return err
	}
	err = user.AppointOrdinary()
	if err != nil {
		return err
	}
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Activate(initiator, user *User) error {
	err := s.assertNotIsAdmin(initiator)
	if err != nil {
		return err
	}
	err = user.Activate()
	if err != nil {
		return err
	}
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Freeze(initiator, user *User) error {
	err := s.assertNotIsAdmin(initiator)
	if err != nil {
		return err
	}
	err = user.Freeze()
	if err != nil {
		return err
	}
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Delete(initiator, user *User) error {
	err := s.assertNotIsAdmin(initiator)
	if err != nil {
		return err
	}
	err = user.Delete()
	if err != nil {
		return err
	}
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) assertAccessOfChange(initiator, user *User) error {
	if initiator.Status().IsOrdinary() && initiator.UserID() != user.UserID() {
		return NoAccessError("вы не можете изменять данные других пользователей")
	}
	return nil
}

func (s *UserService) assertNotIsAdmin(initiator *User) error {
	if initiator.Status().IsOrdinary() {
		return NoAccessError("вы не можете изменять данные других пользователей")
	}
	return nil
}
