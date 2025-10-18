package domain

type EventRepository interface {
	Save(event any) error
}

type EventPublisher struct {
	repository EventRepository
}

func NewEventPublisher(eventRepository EventRepository) *EventPublisher {
	return &EventPublisher{repository: eventRepository}
}

func (ep *EventPublisher) Publish(event any) (err error) {
	err = ep.repository.Save(event)
	return
}
