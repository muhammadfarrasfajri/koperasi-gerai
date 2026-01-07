package repository

import (
	"database/sql"

	"github.com/muhammadfarrasfajri/login-google/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}


func (r *UserRepository) Create(user models.BaseUser) error {
    // Gunakan baris baru agar mudah menghitung jumlah kolom dan tanda tanya
    sqlQuery := `INSERT INTO users (
        google_uid, name, email, nik, npwp, 
        jenis_kelamin, agama, tempat_lahir, tanggal_lahir, 
        alamat_domisili, register_location, register_ip, pekerjaan, 
        status_perkawinan, warga_negara, no_hp, google_picture, profile_picture, ktp_image_path
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)` 
    _, err := r.DB.Exec(sqlQuery, 
        user.GoogleUID, user.Name, user.Email, user.Nik, user.Npwp, 
        user.JenisKelamin, user.Agama, user.TempatLahir, user.TanggalLahir, 
        user.AlamatDomisili, user.RegisterLocation, user.RegisterIP, user.Pekerjaan, 
        user.StatusPerkawinan, user.WargaNegara, user.NoHp, user.GooglePicture, user.ProfilePicture, user.KtpImagePath,
    )
    return err
}

func (r *UserRepository) FindByGoogleUID(uid string) (*models.BaseUser, error) {
	sqlQuery := `SELECT id, google_uid, name, email, google_picture, is_logged_in FROM users WHERE google_uid = ? LIMIT 1`
	row := r.DB.QueryRow(sqlQuery, uid)
	user := models.BaseUser{}
	err := row.Scan(&user.ID, &user.GoogleUID, &user.Name, &user.Email, &user.GooglePicture, &user.IsLoggedIn)
	if err != nil {
		if err == sql.ErrNoRows {
		return nil, err
	}
	return nil, err
}
	return &user, err
}

func (r *UserRepository) UpdateLoginStatus(id int, status int) error {
    query := `UPDATE users SET is_logged_in = ? WHERE id = ?`
    _, err := r.DB.Exec(query, status, id)
    return err
}

func (r *UserRepository) SaveLoginHistory(userID int, deviceInfo, ip string) error {
	sqlQuery := `INSERT INTO login_history_user (user_id, login_at, device_info, ip_address) VALUES (?, NOW(), ?, ?)`
	_, err := r.DB.Exec(sqlQuery, userID, deviceInfo, ip)
	return err
}
