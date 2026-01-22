package models

import "time"

type RefreshToken struct {
	ID           int
	UserID      int
	RefreshToken string
	RevokedAt  time.Time
	ExpiresAt    time.Time
}