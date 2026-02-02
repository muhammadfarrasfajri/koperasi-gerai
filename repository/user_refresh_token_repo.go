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

func (r *UserRefreshTokenRepo) FindRefreshTokenUser(userID int) (*models.RefreshToken, error) {
	sqlQuery := `
		SELECT id, token_hash, expires_at
		FROM user_refresh_tokens
		WHERE user_id = ?
		LIMIT 1
	`

	row := r.DB.QueryRow(sqlQuery, userID)

	token := models.RefreshToken{}
	err := row.Scan(&token.ID, &token.TokenHash, &token.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("refresh token not found")
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *UserRefreshTokenRepo) UpsertRefreshToken(rt models.RefreshToken) error {
	sqlQuery := `
	INSERT INTO user_refresh_tokens (user_id, token_hash, expires_at) 
	VALUES (?, ?, ?) 
	ON DUPLICATE KEY UPDATE
		token_hash = VALUES(token_hash),
		expires_at = VALUES(expires_at),
		revoked_at = NULL
	`

	_, err := r.DB.Exec(
		sqlQuery,
		rt.UserID,
		rt.TokenHash,
		rt.ExpiresAt,
	)
	return err
}

func (r *UserRefreshTokenRepo) DeleteRefreshToken(tokenHash string) error {
    // Perhatikan: parameternya adalah tokenHash (yang sudah di-hash di Service)
    query := "DELETE FROM user_refresh_tokens WHERE token_hash = ?"
    
    result, err := r.DB.Exec(query, tokenHash)
    if err != nil {
        return err
    }

    // (Opsional) Cek apakah ada baris yang terhapus?
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return errors.New("token tidak ditemukan (mungkin sudah logout duluan)")
    }

    return nil
}