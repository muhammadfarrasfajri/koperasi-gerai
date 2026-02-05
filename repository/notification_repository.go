package repository

import (
	"context"

	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
)

type NotificationRepository interface {
	Send(ctx context.Context, notif models.Notification) error
}
