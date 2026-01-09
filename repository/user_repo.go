package repository

import (
	"database/sql"

	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo{
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) FindByNIK(nik string) (*models.BaseUser, error) {
	sqlQuery := `SELECT nik FROM users WHERE nik = ? LIMIT 1`
	row := r.DB.QueryRow(sqlQuery, nik)
	user := models.BaseUser{}
	err := row.Scan(user.NIK)	
	if err != nil {
		if err == sql.ErrNoRows {
		return nil, err
	}
	return nil, err
}
	return &user, nil
}

func (r *UserRepo) IsNIKExists(nik string) (bool, error) {
	query := `SELECT 1 FROM users WHERE nik = ? LIMIT 1`

	var exists int
	err := r.DB.QueryRow(query, nik).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

