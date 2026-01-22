package repository

import "github.com/muhammadfarrasfajri/koperasi-gerai/models"

	type AuthRepository interface {
		CreateRegisterUser(user models.BaseUser) error
		HistoryLoginUser(user models.BaseLoginHistory) error
		IsGoogleUIDExists(googleUID string) (bool, error)
		IsNIKExists(nik string) (bool, error)
		GenerateMemberID(prefix string) (string, error)
	} 