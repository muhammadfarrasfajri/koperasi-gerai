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

func (r *UserRepo) FindByGoogleUID(uid string) (*models.BaseUser, error) {
	sqlQuery := `SELECT id, google_uid, name, email, is_verified FROM users WHERE google_uid = ? LIMIT 1`
	row := r.DB.QueryRow(sqlQuery, uid)
	user := models.BaseUser{}
	err := row.Scan(&user.ID, &user.GoogleUID, &user.Name, &user.Email, &user.Is_verified)
	if err != nil {
		if err == sql.ErrNoRows {
		return nil, err
	}
	return nil, err
}
	return &user, err
}

func (r *UserRepo) FindById(id string) (*models.BaseUser, error) {
	sqlQuery := `SELECT id, google_uid, name, email FROM users WHERE id = ? LIMIT 1`
	row := r.DB.QueryRow(sqlQuery, id)
	user := models.BaseUser{}
	err := row.Scan(&user.ID, &user.GoogleUID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}
