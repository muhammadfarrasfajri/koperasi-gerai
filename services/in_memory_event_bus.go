package services

import (
	"context"

	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
)

type InMemoryEventBus struct {
	notficationService *NotificationService
}

func NewInMemoryEventBus(
	n *NotificationService,
) *InMemoryEventBus {
	return &InMemoryEventBus{
		notficationService: n,
	}
}

func (b *InMemoryEventBus) UserRegistered(
	ctx context.Context,
	event models.UserRegisteredEvent,
) {
	go b.notficationService.OnUserRegistered(ctx, event)
}
