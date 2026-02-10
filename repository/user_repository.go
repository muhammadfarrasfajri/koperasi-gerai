package repository

import "github.com/muhammadfarrasfajri/koperasi-gerai/models"

type UserRepository interface {
	FindByNIK(nik string) (*models.BaseUser, error)
	FindByGoogleUID(uid string) (*models.BaseUser, error)
	FindById(id string) (*models.BaseUser, error)
	
}