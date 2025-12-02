package service

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, event any)
}
