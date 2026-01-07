package repository

import (
	"database/sql"
	"fmt"

	"github.com/muhammadfarrasfajri/login-google/models"
)

// --------------------------- GET ALL USERS -----------------------------------

func (r *UserRepository) GetAll() ([]models.BaseUser, error) {
	sqlQuery := `SELECT id, google_uid, name, email, google_picture, role, profile_picture FROM users`
	rows, err := r.DB.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.BaseUser{}

	for rows.Next() {
		u := models.BaseUser{}
		err := rows.Scan(&u.ID, &u.GoogleUID, &u.Name, &u.Email, &u.GooglePicture, &u.Role, &u.ProfilePicture)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// --------------------------- FIND BY ID --------------------------------------

func (r *UserRepository) FindByID(id string) (*models.BaseUser, error) {
	row := r.DB.QueryRow(`
        SELECT id, google_uid, name, email, google_picture, role, profile_picture
        FROM users WHERE id = ?
    `, id)

	user := models.BaseUser{}
	err := row.Scan(&user.ID, &user.GoogleUID, &user.Name, &user.Email, &user.GooglePicture, &user.Role, &user.ProfilePicture)

	if err == sql.ErrNoRows {
		fmt.Println("No user found with ID:", err)
		return nil, nil
	}

	return &user, err
}

// --------------------------- UPDATE USER -------------------------------------

func (r *UserRepository) Update(user models.BaseUser) error {
	_, err := r.DB.Exec(`
        UPDATE users SET 
            name = ?, 
            email = ?, 
            role = ?,
			profile_picture = ?
        WHERE id = ?
    `, user.Name, user.Email, user.Role, user.ProfilePicture, user.ID)

	return err
}

// --------------------------- DELETE USER -------------------------------------

func (r *UserRepository) Delete(id string) error {
	// Baru hapus user
	_, err := r.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// --------------------------- UPDATE PHOTO URL -------------------------------------
func (r *UserRepository) UpdatePhotoURL(userID int, url string) error {
	// Update kolom avatar_url di tabel users
	_, err := r.DB.Exec("UPDATE users SET profile_picture = ? WHERE id = ?", url, userID)
	fmt.Println("errornya:", err)
	return err

}
