package models

import "time"

type RefreshToken struct {
	ID           int
	AdminOrUserID      int
	RefreshToken string
	ExpiresAt    time.Time
}