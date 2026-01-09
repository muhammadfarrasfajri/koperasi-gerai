package repository

import (
	"database/sql"

	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
)

type UserAuthRepo struct {
	DB *sql.DB
}

func NewUserAuthRepo(db *sql.DB) *UserAuthRepo{
	return &UserAuthRepo{
		DB: db,
	}
}

func (r *UserAuthRepo) CreateRegisterUser(user models.BaseUser) error {
	sqlQuery := `INSERT INTO users ( id_member, google_uid, name, email, nik, npwp, gender, religion, place_of_birth, birth, address, register_location, register_ip, job, marital_status, citizenship, phone_number, google_picture, profile_picture, ktp_picture) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `
	_, err := r.DB.Exec( sqlQuery, user.IDMember, user.GoogleUID, user.Name, user.Email, user.NIK, user.NPWP, user.Gender, user.Religion, user.PlaceOfBirth, user.Birth, user.Address, user.RegisterLocation, user.RegisterIP, user.Job, user.MaritalStatus, user.Citizenship, user.PhoneNumber, user.GooglePicture, user.ProfilePicture, user.KtpPicture)

	return err
}

func (r *UserAuthRepo) IsNIKExists(nik string) (bool, error) {
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

