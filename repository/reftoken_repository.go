package repository

import (
	"time"

	"github.com/muhammadfarrasfajri/login-google/models"
)

type RefreshTokenRepository interface {
	RefreshToken(userID int, refreshToken string, exp time.Time) error
	FindRefreshToken(userID int) (*models.RefreshToken, error)
	UpdateRefreshToken(userID int, newRefreshToken string, exp time.Time) error
	DeleteRefreshToken(UserID int) error
}