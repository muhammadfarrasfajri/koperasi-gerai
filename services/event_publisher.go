package services

import (
	"context"

	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
)

type EventPublisher interface {
	UserRegistered(ctx context.Context, event models.UserRegisteredEvent)
}
