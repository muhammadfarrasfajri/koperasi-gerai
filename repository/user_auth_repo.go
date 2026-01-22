package repository

import (
	"database/sql"
	"fmt"
	"time"

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
	sqlQuery := `INSERT INTO users (id_member, google_uid, name, email, npwp, nik, place_of_birth, birth, gender, address, pos_code, religion, marital_status, job, citizenship, blood_type, phone_number, register_location, register_ip, ktp_picture, google_picture, last_education, active_as, mother_name) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec( sqlQuery, user.IDMember, user.GoogleUID, user.Name, user.Email, user.NPWP, user.NIK, user.PlaceOfBirth, user.Birth, user.Gender, user.Address, user.PosCode, user.Religion, user.MaritalStatus, user.Job, user.Citizenship, user.Blood_type,  user.PhoneNumber, user.RegisterLocation, user.RegisterIP, user.KtpPicture, user.GooglePicture, user.LastEducation, user.ActiveAs, user.Mother_name)

	return err
}

func (r *UserAuthRepo) HistoryLoginUser(user models.BaseLoginHistory) error {
	sqlQuery := `INSERT INTO history_login_user (user_id, login_at, status, user_agent, ip_address, device_info, location, ) VALUES (?, NOW(), ?, ?)`
	_, err := r.DB.Exec(sqlQuery, user.UserID, user.LoginAt, user.Status, user.UserAgent, user.IPAddress, user.DeviceInfo, user.Location)
	return err
}

func (r UserAuthRepo) IsGoogleUIDExists(googleUID string) (bool, error) {
	query := `SELECT 1 FROM users WHERE google_uid = ? LIMIT 1`

	var exists int
	err := r.DB.QueryRow(query, googleUID).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r UserAuthRepo) GenerateMemberID(prefix string) (string, error) {
	// 1. Ambil waktu sekarang untuk Year (YY) dan Month (MM)
	now := time.Now()
	yearStr := now.Format("06")  // 2 Digit Tahun (misal: 26)
	monthStr := now.Format("01") // 2 Digit Bulan (misal: 01)

	var lastNumber int

	// 2. Query Locking: SELECT ... FOR UPDATE
	// Query ini akan mencari counter sekaligus MENGUNCI baris tersebut.
	querySelect := `
		SELECT last_number 
		FROM member_counters 
		WHERE prefix = ? AND year = ? AND month = ? 
		FOR UPDATE
	`
	// Eksekusi query
	err := r.DB.QueryRow(querySelect, prefix, yearStr, monthStr).Scan(&lastNumber)

	if err != nil {
		if err == sql.ErrNoRows {
			// KASUS A: Data belum ada (Awal Bulan), kita INSERT baru
			lastNumber = 1
			queryInsert := `
				INSERT INTO member_counters (prefix, year, month, last_number) 
				VALUES (?, ?, ?, ?)
			`
			_, errInsert := r.DB.Exec(queryInsert, prefix, yearStr, monthStr, lastNumber)
			if errInsert != nil {
				return "", errInsert
			}
		} else {
			// Error lain (koneksi putus, dsb)
			return "", err
		}
	} else {
		// KASUS B: Data ditemukan, kita UPDATE
		lastNumber++ // Tambah 1
		queryUpdate := `
			UPDATE counters 
			SET last_number = ? 
			WHERE prefix = ? AND year = ? AND month = ?
		`
		_, errUpdate := r.DB.Exec(queryUpdate, lastNumber, prefix, yearStr, monthStr)
		if errUpdate != nil {
			return "", errUpdate
		}
	}

	// 3. Format hasil akhir (Contoh: MBR260100005)
	newID := fmt.Sprintf("%s%s%s%05d", prefix, yearStr, monthStr, lastNumber)
	
	return newID, nil
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

