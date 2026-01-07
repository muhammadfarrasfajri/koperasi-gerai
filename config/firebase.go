package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirebaseAppAdmin *firebase.App
var FirebaseAppUser *firebase.App

func InitFirebase() {
	optAdmin := option.WithCredentialsFile("firebase-key-admin.json")
	appAdmin, err := firebase.NewApp(context.Background(), nil, optAdmin)
	if err != nil {
		log.Fatalf("Failed to init Firebase: %v", err)
	}
	
	optUser := option.WithCredentialsFile("firebase-key-user.json")
	appUser, err := firebase.NewApp(context.Background(), nil, optUser)
	if err != nil {
		log.Fatalf("Failed to init Firebase: %v", err)
	}
	FirebaseAppAdmin = appAdmin
	FirebaseAppUser = appUser

}
