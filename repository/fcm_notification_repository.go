package repository

import (
	"context"
	"log"

	"firebase.google.com/go/messaging"
	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
)

type FCMNotificationRepository struct {
	client *messaging.Client
}

func NewFCMNotificationRepository(c *messaging.Client) *FCMNotificationRepository {
	return &FCMNotificationRepository{client: c}
}

func (r *FCMNotificationRepository) Send(
	ctx context.Context,
	notif models.Notification,
) error {
	msg := &messaging.Message{
		Token: notif.Token,
		Notification: &messaging.Notification{
			Title: notif.Title,
			Body:  notif.Body,
		},
		Data: notif.Data,
	}

	resp, err := r.client.Send(ctx, msg)
	log.Println("FCM Response:", resp, err)
	return err
}
