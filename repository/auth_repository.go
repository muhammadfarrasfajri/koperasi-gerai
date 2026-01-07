package repository

import (
	"github.com/muhammadfarrasfajri/login-google/models"
)

type AuthRepository interface {
	// Register and Login
	Create(user models.BaseUser) error
	SaveLoginHistory(userID int, deviceInfo, ip string) error
	UpdateLoginStatus(id int, status int) error

	//CRUD
	FindByGoogleUID(uid string) (*models.BaseUser, error)
	FindByID(id string) (*models.BaseUser, error)
	GetAll() ([]models.BaseUser, error)
	Update(user models.BaseUser) error
	Delete(id string) error
	UpdatePhotoURL(userID int, url string) error
}
