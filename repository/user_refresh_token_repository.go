package repository

import (
	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
)

type UserRefreshTokenRepository interface {
	FindRefreshTokenUser(userID int) (*models.RefreshToken, error)
	CreateRefreshTokenUser(historyLogin models.RefreshToken) error
	UpdateRefreshTokenUser(historyLogin models.RefreshToken) error
}