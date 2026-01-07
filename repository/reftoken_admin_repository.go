package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/muhammadfarrasfajri/login-google/models"
)


type RefreshTokenAdminRepository struct {
	DB *sql.DB
}

func NewRefreshTokenAdmin(db *sql.DB) *RefreshTokenAdminRepository {
	return &RefreshTokenAdminRepository{
		DB: db,
	}
}

func (r *RefreshTokenAdminRepository) FindRefreshToken(adminID int) (*models.RefreshToken, error) {
	sqlQuery := `SELECT id, admin_id, refresh_token, expires_at FROM refresh_tokens_admin WHERE admin_id = ?`
	row := r.DB.QueryRow(sqlQuery, adminID)
	admin := models.RefreshToken{}
	err := row.Scan(&admin.ID, &admin.AdminOrUserID, &admin.RefreshToken, &admin.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("refresh token not found")
	}
	return &admin, err
}

func (r *RefreshTokenAdminRepository) RefreshToken(adminID int, refreshToken string, exp time.Time) error {
	sqlQuery := `INSERT INTO refresh_tokens_admin (admin_id, refresh_token, expires_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE refresh_token = VALUES(refresh_token), expires_at = VALUES(expires_at)`
	formatted := exp.Format("2006-01-02 15:04:05")
	_, err := r.DB.Exec(sqlQuery, adminID, refreshToken, formatted)
	return err
}

func (r *RefreshTokenAdminRepository) UpdateRefreshToken(adminID int, newRefreshToken string, exp time.Time) error {
	sqlQuery := `UPDATE refresh_tokens_admin SET refresh_token = ?, expires_at = ? WHERE admin_id = ?`
	_, err := r.DB.Exec(sqlQuery, newRefreshToken, exp, adminID)
	return err
}

func (r *RefreshTokenAdminRepository) DeleteRefreshToken(adminID int) error {
	sqlQuery := `DELETE FROM refresh_tokens_admin WHERE admin_id = ?`
	_, err := r.DB.Exec(sqlQuery, adminID)
	if err != nil {
		return err
	}
	return nil
}
