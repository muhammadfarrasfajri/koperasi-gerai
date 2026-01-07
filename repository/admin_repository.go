package repository

import (
	"database/sql"
	"errors"

	"github.com/muhammadfarrasfajri/login-google/models"
)

// --------------------------- GET ALL ADMINS -----------------------------------

func (r *AdminRepository) GetAll() ([]models.BaseUser, error) {
	sqlQuery := `SELECT id, google_uid, name, email, google_picture, role, profile_picture FROM admins`
	rows, err := r.DB.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	admins := []models.BaseUser{}

	for rows.Next() {
		u := models.BaseUser{}
		err := rows.Scan(&u.ID, &u.GoogleUID, &u.Name, &u.Email, &u.GooglePicture, &u.Role, &u.ProfilePicture)
		if err != nil {
			return nil, err
		}
		admins = append(admins, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return admins, nil
}

// --------------------------- FIND BY ID --------------------------------------

// Get User Use Id
func (r *AdminRepository) FindByID(id string) (*models.BaseUser, error) {
	sqlQuery := `SELECT id, google_uid, name, email, google_picture, role FROM admins WHERE id = ?`
	row := r.DB.QueryRow(sqlQuery, id)
	admin := models.BaseUser{}
	err := row.Scan(&admin.ID, &admin.GoogleUID, &admin.Name, &admin.Email, &admin.GooglePicture, &admin.Role)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &admin, err
}

// --------------------------- UPDATE USER -------------------------------------

func (r *AdminRepository) Update(admin models.BaseUser) error {
	sqlQuery := `UPDATE admins SET name = ?, email = ?, role = ?, profile_picture = ? WHERE id = ?`
	_, err := r.DB.Exec(sqlQuery, admin.Name, admin.Email, admin.Role, admin.ProfilePicture, admin.ID)

	return err
}

// --------------------------- DELETE USER -------------------------------------

func (r *AdminRepository) Delete(id string) error {
	sqlQuery := `DELETE FROM admins WHERE id = ?`
	_, err := r.DB.Exec(sqlQuery, id)
	return err
}

// --------------------------- UPDATE PHOTO URL -------------------------------------
func (r *AdminRepository) UpdatePhotoURL(adminID int, url string) error {
	sqlQuery := `UPDATE admins SET profile_picture = ? WHERE id = ?`
	_, err := r.DB.Exec(sqlQuery, url, adminID)
	return err

}
