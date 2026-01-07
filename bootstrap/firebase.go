package bootstrap

import (
	"context"
	"log"

	"firebase.google.com/go/auth"
	"github.com/muhammadfarrasfajri/login-google/config"
)

func InitFirebase() (adminAuth, userAuth *auth.Client) {
	config.InitFirebase()

	adminApp, err := config.FirebaseAppAdmin.Auth(context.Background())
	if err != nil {
		log.Fatal("Failed to init Firebase Admin:", err)
	}

	userApp, err := config.FirebaseAppUser.Auth(context.Background())
	if err != nil {
		log.Fatal("Failed to init Firebase User:", err)
	}

	return adminApp, userApp
}
