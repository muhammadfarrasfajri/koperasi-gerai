	package repository

	import "github.com/muhammadfarrasfajri/koperasi-gerai/models"

	type AuthRepository interface {
		CreateRegisterUser(models.BaseUser) error
	}