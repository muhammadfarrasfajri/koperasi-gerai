package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	firebase "firebase.google.com/go/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/muhammadfarrasfajri/koperasi-gerai/middleware"
	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
	"github.com/muhammadfarrasfajri/koperasi-gerai/repository"
)

var (
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrUserNotRegistered = errors.New("user not registered, please register first")
	ErrExistingNIK = errors.New("NIK already usage")
	ErrInvalidNIK =  errors.New("Invalid NIK")
	ErrInvalidNoHp = errors.New("Mobile number must consist of 11-13 digits")
)

type UserAuthService struct {
	AuthRepo repository.AuthRepository
	UserRepo repository.UserRepository
	RefRepo repository.UserRefreshTokenRepository
	FirebaseAuth *firebase.Client
	JWTSecret    *middleware.JWTManager
}

func NewUserAuthService(userAuthRepository repository.AuthRepository, userRepository repository.UserRepository, refreshRepository repository.UserRefreshTokenRepository, firebaseAuth *firebase.Client, jwtsecret *middleware.JWTManager) *UserAuthService{
	return &UserAuthService{
		AuthRepo: userAuthRepository,
		UserRepo: userRepository,
		RefRepo: refreshRepository,
		FirebaseAuth: firebaseAuth,
		JWTSecret: jwtsecret,
	}
}

// --------------------------- REGISTER -----------------------------------

func (s *UserAuthService) Register(idToken string, user models.BaseUser) (*models.BaseUser, error) {
	ctx := context.Background()

	// 1. Verifikasi Firebase ID Token
	token, err := s.FirebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user.GoogleUID = token.UID
	
	if email, ok := token.Claims["email"].(string); ok {
        user.Email = email
    }
	
    if pic, ok := token.Claims["picture"].(string); ok {
        user.GooglePicture = pic
    }

	// Generate member ID
	newMemberID, err := s.AuthRepo.GenerateMemberID("KOP")
	if err != nil {
		return nil, errors.New("Generate ID member failed")
	}
	
	user.IDMember = newMemberID

	// 2. Cek apakah user sudah ada
	existing, err := s.AuthRepo.IsNIKExists(user.NIK)
	if err != nil {
		return nil, err
	}

	if existing {
		return nil, ErrExistingNIK
	}

	// 3. Tentukan nama
	if user.Name == "" {
        if n, ok := token.Claims["name"].(string); ok {
            user.Name = n
        }
    }

	// 4. Cek panjang no HP
	length := len(user.PhoneNumber)
	if length < 11 || length > 13 {
		return nil, ErrInvalidNoHp
	}

	// 5.Cek format NIK
	nik := len(user.NIK)
	if nik != 16 {
		return nil, ErrInvalidNIK
	}

	err = s.AuthRepo.CreateRegisterUser(user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserAuthService) Login(loginHistory models.BaseLoginHistory) (map[string]interface{}, error) {
	ctx := context.Background()

	tokenModel := models.RefreshToken{}
	// 1. Verifikasi Firebase Token
	token, err := s.FirebaseAuth.VerifyIDToken(ctx, loginHistory.IdToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	googleUID := token.UID

	// 2. Cek user di DB
	user, err := s.UserRepo.FindByGoogleUID(googleUID)
	if err != nil {
		return nil, err
	}
	
	if user.Is_verified == 0 {
		return nil, errors.New("User is not verified")
	}

	// 3. Cek user login
	u, err := s.RefRepo.FindRefreshTokenUser(user.ID)

	if err != nil || u == nil{
		// 4. Simpan Aktivitas Login
	    err = s.AuthRepo.HistoryLoginUser(loginHistory) 
		if err != nil {
			return nil, err
		}
	
		//5. Generate Access token
		accessToken, err := s.JWTSecret.GenerateAccessToken(user.ID, user.Email)
		if err != nil {
			return nil, err
		}

		//6. Generate Referesh Token
		refreshToken, err := s.JWTSecret.GenerateRefreshToken(user.ID)
		if err != nil {
			return nil, err
		}

		tokenModel.RefreshToken  = refreshToken

		//7. time exp refresh token
		expiresAt := time.Now().Add(7 * 24 * time.Hour) 

		tokenModel.ExpiresAt = expiresAt

		//8. Send refresh token to database
		err = s.RefRepo.CreateRefreshTokenUser(tokenModel)
		if err != nil {
			return nil, err
		}
	
		return map[string]interface{}{
			"message": "login success",
			"access_token":   accessToken,
			"refresh_token": refreshToken,
		}, nil
	}

	// 5. Simpan aktivitas login
	err = s.AuthRepo.HistoryLoginUser(loginHistory)
	if err != nil {
		return nil, err
	}

	//6. Generate Access token
	accessToken, err := s.JWTSecret.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	
	//7. Generate Referesh Token
	refreshToken, err := s.JWTSecret.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	tokenModel.RefreshToken = refreshToken
	
	//10. time exp refresh token
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	tokenModel.ExpiresAt = expiresAt

	//11. Send refresh token to database
	err = s.RefRepo.UpdateRefreshTokenUser( tokenModel)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "login success",
		"access_token":   accessToken,
		"refresh_token": refreshToken,
		}, nil
}

func (s *UserAuthService) RefreshToken(refreshToken string) (map[string]interface{}, error) {

	tokenModel := models.RefreshToken{}
		// parsing token
		token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
			return s.JWTSecret.RefreshSecret, nil
		})
			
		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))
	
	
		tokenCheck, err := s.RefRepo.FindRefreshTokenUser(userID)
	
	
		if err != nil || tokenCheck == nil {
			return nil, errors.New("refresh token not found")
		}
	
		if refreshToken != tokenCheck.RefreshToken {
			return nil, errors.New("refresh token not match")
		} 
	
		// ambil user dari db
		user, err := s.UserRepo.FindById(strconv.Itoa(userID))
		if err != nil || user == nil {
			return nil, err
		}
	
		// generate token baru (access + refresh)
		accessToken, err := s.JWTSecret.GenerateAccessToken(user.ID, user.Email)
		if err != nil {
			return nil, err
		}
	
		newRefreshToken, err := s.JWTSecret.GenerateRefreshToken(user.ID)
		if err != nil {
			return nil, err
		}

		tokenModel.RefreshToken = refreshToken
		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		tokenModel.ExpiresAt = expiresAt
	
		err = s.RefRepo.UpdateRefreshTokenUser(tokenModel)
		if err != nil {
			return nil, err
		}
	
		return map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": newRefreshToken,
		}, nil
}