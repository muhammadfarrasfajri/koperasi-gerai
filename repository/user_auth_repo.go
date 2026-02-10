package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
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
    sqlQuery := `INSERT INTO users (
        id_member, google_uid, name, email, npwp, nik,
        place_of_birth, birth, gender, address, pos_code, religion,
        marital_status, job, citizenship, blood_type, phone_number,
        register_location, register_ip, ktp_picture, profile_picture, google_picture,
        last_education, active_as, mother_name
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

    _, err := r.DB.Exec(sqlQuery,
        user.IDMember, user.GoogleUID, user.Name, user.Email, user.NPWP, user.NIK,
        user.PlaceOfBirth, user.Birth, user.Gender, user.Address, user.PosCode, user.Religion,
        user.MaritalStatus, user.Job, user.Citizenship, user.Blood_type, user.PhoneNumber,
        user.RegisterLocation, user.RegisterIP, user.KtpPicture, user.ProfilePicture, user.GooglePicture,
        user.LastEducation, user.ActiveAs, user.Mother_name,
    )

    if err != nil {

        if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
            
            errorMessage := mysqlErr.Message 

            if strings.Contains(errorMessage, "nik") {
                 return errors.New("NIK sudah terdaftar")
            }

            if strings.Contains(errorMessage, "phone") || strings.Contains(errorMessage, "no_hp") {
                 return errors.New("Nomor HP sudah terdaftar")
            }

            if strings.Contains(errorMessage, "email") {
                 return errors.New("Email sudah terdaftar")
            }
            
            return errors.New("Data akun sudah ada (Duplicate User)")
        }
        
        return err
    }

    return nil
}

func (r *UserAuthRepo) FindByEmail(email string) (*models.BaseUser, error) {

    query := `SELECT id_member, name, email, google_uid FROM users WHERE email = ? LIMIT 1`

    var user models.BaseUser
    err := r.DB.QueryRow(query, email).Scan(&user.IDMember, &user.Name, &user.Email, &user.GoogleUID)

    if err != nil {
        return nil, err 
    }

    return &user, nil
}

func (r *UserAuthRepo) LinkGoogleAccount(email string, googleUID string, googlePic string) error {
    query := `UPDATE users 
              SET google_uid = ?, google_picture = ? 
              WHERE email = ?`

    _, err := r.DB.Exec(query, googleUID, googlePic, email)
    
    return err
}


func (r *UserAuthRepo) HistoryLoginUser(user models.BaseLoginHistory) error {
    sqlQuery := `
        INSERT INTO user_login_histories
        (user_id, login_at, status, user_agent, ip_address, device_info, location, error_message)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    // Debug 1: Print Query & Data yang mau dimasukkan
    fmt.Println("\n--- [REPO DEBUG START] ---")
    fmt.Printf("Mencoba Insert ID: %d, Waktu: %v, Status: %s\n", user.UserID, user.LoginAt, user.Status)

    res, err := r.DB.Exec(
        sqlQuery,
        user.UserID,
        user.LoginAt,
        user.Status,
        user.UserAgent,
        user.IPAddress,
        user.DeviceInfo,
        user.Location,
        user.ErrorMessage,
    )

    if err != nil {
        fmt.Printf("!!! REPO ERROR: %v\n", err)
        return err
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        fmt.Println("!!! BAHAYA: Tidak ada Error, TAPI RowsAffected = 0. Data DITOLAK Database !!!")
    } else {
        fmt.Printf(">>> SUKSES INSERT: %d baris berhasil masuk ke tabel.\n", rows)
    }

    // 4. CEK TRANSAKSI (GORM / SQLX)
    // Jika r.DB adalah *sql.Tx, kita harus memastikan dia di-commit di level service/handler
    // Jika r.DB adalah *sql.DB (koneksi biasa), ini otomatis commit.
    
    fmt.Println("--- [REPO DEBUG END] ---")
    return nil
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

func (r *UserAuthRepo) IsNoHPExists(noHp string) (bool, error) {
	query := `SELECT 1 FROM users WHERE phone_number = ? LIMIT 1`

	var exists int
	
	err := r.DB.QueryRow(query, noHp).Scan(&exists)
	
	if err == sql.ErrNoRows {
		return false, nil
	}
	
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserAuthRepo) GetMemberId(prefix string) (string, error){

	var lastID string

	query := `SELECT id_member FROM users WHERE id_member LIKE ? ORDER BY id_member DESC LIMIT 1`
	
	searchPattern := prefix + "%"

	err := r.DB.QueryRow(query, searchPattern).Scan(&lastID)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return lastID, nil
}