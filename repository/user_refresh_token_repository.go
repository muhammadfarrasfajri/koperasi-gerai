package repository

import (
	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
)

type UserRefreshTokenRepository interface {
	FindRefreshTokenUser(userID int) (*models.RefreshToken, error)
	UpsertRefreshToken(refreshToken models.RefreshToken) error
	DeleteRefreshToken(token string) error
}