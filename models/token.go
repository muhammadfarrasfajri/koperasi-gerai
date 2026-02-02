package models

import "time"

type RefreshToken struct {
	ID           int
	UserID      int
	TokenHash string 
	RevokedAt  time.Time
	ExpiresAt    time.Time
}