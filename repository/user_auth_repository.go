package repository

import "github.com/muhammadfarrasfajri/koperasi-gerai/models"

	type AuthRepository interface {
		CreateRegisterUser(user models.BaseUser) error
		HistoryLoginUser(user models.BaseLoginHistory) error
		IsNIKExists(nik string) (bool, error)
		GetMemberId(prefix string) (string, error)
		IsNoHPExists(noHp string) (bool, error)
		FindByEmail(email string) (*models.BaseUser, error)
		LinkGoogleAccount(email string, googleUID string, googlePic string) error
	} 