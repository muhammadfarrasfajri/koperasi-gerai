package repository

import (
	"context"
	"database/sql"
)

type AdminDeviceTokenRepository interface {
	GetAllTokens(ctx context.Context) ([]string, error)
}

type adminDeviceTokenRepo struct {
	db *sql.DB
}

func NewAdminDeviceTokenRepository(db *sql.DB) AdminDeviceTokenRepository {
	return &adminDeviceTokenRepo{db: db}
}

func (r *adminDeviceTokenRepo) GetAllTokens(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT fcm_token
		FROM admin_device_tokens
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
