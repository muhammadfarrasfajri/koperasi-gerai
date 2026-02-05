package services

import (
	"context"
	"log"

	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
	"github.com/muhammadfarrasfajri/koperasi-gerai/repository"
)

type AdminDeviceTokenRepository interface {
	GetAllTokens(ctx context.Context) ([]string, error)
}

type NotificationService struct {
	notifRepo repository.NotificationRepository
	tokenRepo AdminDeviceTokenRepository
}

func NewNotificationService(
	repo repository.NotificationRepository,
	tokenRepo AdminDeviceTokenRepository,
) *NotificationService {
	return &NotificationService{
		notifRepo: repo,
		tokenRepo: tokenRepo,
	}
}

func (s *NotificationService) OnUserRegistered(
	ctx context.Context,
	event models.UserRegisteredEvent,
) {
	// 1️⃣ Ambil token admin dari database
	tokens, err := s.tokenRepo.GetAllTokens(ctx)
	if err != nil {
		log.Println("❌ failed to get admin tokens:", err)
		return
	}

	if len(tokens) == 0 {
		log.Println("⚠️  no admin device token found, skip notification")
		return
	}

	// 2️⃣ Kirim notif ke setiap token
	for _, token := range tokens {
		notif := models.Notification{
			Title: "New User Registered",
			Body:  "Klik untuk verifikasi KTP",
			Token: token,
			Data: map[string]string{
				"user_id": event.IDMember,
				"action":  "verify_ktp",
			},
		}

		log.Println("📨 Sending notification to admin token:", token[:20]+"...")

		err := s.notifRepo.Send(ctx, notif)
		if err != nil {
			log.Println("❌ failed to send notification:", err)
			continue
		}

		log.Println("✅ Notification sent")
	}
}
