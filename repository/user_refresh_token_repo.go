package repository

import (
	"database/sql"
	"errors"

	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
)

type UserRefreshTokenRepo struct {
	DB *sql.DB
}

func NewUserRefreshTokenRepo(db *sql.DB) *UserRefreshTokenRepo{
	return &UserRefreshTokenRepo{
		DB: db,
	}
}

func (r *UserRefreshTokenRepo) FindRefreshTokenUser(userID int) (*models.RefreshToken, error){
	sqlQuery := `SELECT id, token_hash, expires_at FROM users WHERE user_id = ? LIMIT 1`
	row := r.DB.QueryRow(sqlQuery, userID)
	user := models.RefreshToken{}
	err := row.Scan(&user.ID, &user.RefreshToken, &user.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("refresh token not found")
	}
	return &user, err
}

func (r *UserRefreshTokenRepo) CreateRefreshTokenUser(refreshToken models.RefreshToken) error{
	sqlQuery := `INSERT INTO user_refresh_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE token_hash = VALUES(token_hash), expires_at = VALUES(expires_at)`
	formatted := refreshToken.ExpiresAt.Format("2006-01-02 15:04:05")
	_, err := r.DB.Exec(sqlQuery, refreshToken.UserID, refreshToken.RefreshToken, formatted)
	return err
}

func (r *UserRefreshTokenRepo) UpdateRefreshTokenUser(refreshToken models.RefreshToken) error{
	sqlQuery := `UPDATE user_refresh_tokens SET token_hash = ?, expires_at = ? WHERE user_id = ?`
	formatted := refreshToken.ExpiresAt.Format("2006-01-02 15:04:05")
	_, err := r.DB.Exec(sqlQuery, refreshToken.RefreshToken, formatted, refreshToken.UserID)
	return err
}