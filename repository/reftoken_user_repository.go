package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/muhammadfarrasfajri/login-google/models"
)

type RefreshTokenUserRepository struct {
	DB *sql.DB
}

func NewRefreshTokenUser(db *sql.DB) *RefreshTokenUserRepository {
	return &RefreshTokenUserRepository{
		DB: db,
	}
}

func (r *RefreshTokenUserRepository) FindRefreshToken(userID int) (*models.RefreshToken, error) {
	sqlQuery := `SELECT id, user_id, refresh_token, expires_at FROM refresh_tokens_user WHERE user_id = ?`
	row := r.DB.QueryRow(sqlQuery, userID)
	user := models.RefreshToken{}
	err := row.Scan(&user.ID,&user.AdminOrUserID, &user.RefreshToken, &user.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("data tidak ada")
	}
	return &user, err
}

func (r *RefreshTokenUserRepository) RefreshToken(userID int, refreshToken string, exp time.Time) error {
	sqlQuery := `INSERT INTO refresh_tokens_user (user_id, refresh_token, expires_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE refresh_token = VALUES(refresh_token), expires_at = VALUES(expires_at)`
	formatted := exp.Format("2006-01-02 15:04:05")
	_, err := r.DB.Exec(sqlQuery, userID, refreshToken, formatted)
	return err
}

func (r *RefreshTokenUserRepository) UpdateRefreshToken(userID int, newRefreshToken string, exp time.Time) error {
	sqlQuery := `UPDATE refresh_tokens_user SET refresh_token = ?, expires_at = ? WHERE user_id = ?`
	_, err := r.DB.Exec(sqlQuery, newRefreshToken, exp, userID)
	return err
}

func (r *RefreshTokenUserRepository) DeleteRefreshToken(userID int) error {
	sqlQuery := `DELETE FROM refresh_tokens_user WHERE user_id = ?`
	_, err := r.DB.Exec(sqlQuery, userID)
	if err != nil {
		return err
	}
	return nil
}
