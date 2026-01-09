package repository

import "github.com/muhammadfarrasfajri/koperasi-gerai/models"

type UserRepository interface {
	FindByNIK(nik string) (*models.BaseUser, error)
	IsNIKExists(nik string) (bool, error)
}