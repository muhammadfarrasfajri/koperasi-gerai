package services

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	firebase "firebase.google.com/go/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/muhammadfarrasfajri/login-google/middleware"
	"github.com/muhammadfarrasfajri/login-google/models"
	"github.com/muhammadfarrasfajri/login-google/repository"
)

var (
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrUserNotRegistered = errors.New("user not registered, please register first")
)

type AuthService struct {
	Repo     repository.AuthRepository
	Ref      repository.RefreshTokenRepository
	FirebaseAuth *firebase.Client
	JWTSecret    *middleware.JWTManager
}

func NewAuthService(repositoryUser repository.AuthRepository, refRefo repository.RefreshTokenRepository, firebaseAuth *firebase.Client, jwtsecret *middleware.JWTManager) *AuthService{
	return &AuthService{
		Repo: repositoryUser,
		Ref: refRefo,
		FirebaseAuth: firebaseAuth,
		JWTSecret: jwtsecret,
	}
}

// --------------------------- REGISTER -----------------------------------

func (s *AuthService) Register(idToken string, user models.BaseUser) (*models.BaseUser, error) {
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

	// 2. Cek apakah user sudah ada
	existing, _ := s.Repo.FindByGoogleUID(user.GoogleUID)
	if existing != nil {
		return nil, errors.New("user already registered, please login")
	}

	// 3. Tentukan nama
	if user.Name == "" {
        if n, ok := token.Claims["name"].(string); ok {
            user.Name = n
        }
    }
	
	length := len(user.NoHp)
	
	if length < 11 || length > 13 {
		return nil, errors.New("nomor HP harus 11–13 digit")
	}

	err = s.Repo.Create(user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// -------------------------- LOGIN ----------------------------------------

func (s *AuthService) Login(idToken string, deviceInfo string, ip string) (map[string]interface{}, error) {
	ctx := context.Background()

	// 1. Verifikasi Firebase Token
	token, err := s.FirebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	googleUID := token.UID

	// 2. Cek user di DB
	user, err := s.Repo.FindByGoogleUID(googleUID)
	if err != nil || user == nil {
		return nil, ErrUserNotRegistered
	}

	// 3. Cek user login
	u, err := s.Ref.FindRefreshToken(user.ID)
	if err != nil || u == nil{
		// 4. Update status login
		if err := s.Repo.UpdateLoginStatus(user.ID, 1); err != nil {
			return nil, err
		}
	
		// 5. Simpan aktivitas login
		err = s.Repo.SaveLoginHistory(user.ID, deviceInfo, ip)
		if err != nil {
			return nil, err
		}
	
		//6. Generate Access token
		accessToken, err := s.JWTSecret.GenerateAccessToken(user.ID, user.Email, user.Role)
		if err != nil {
			return nil, err
		}
	
		//7. Generate Referesh Token
		refreshToken, err := s.JWTSecret.GenerateRefreshToken(user.ID)
		if err != nil {
			return nil, err
		}
	
		//8. encryp Refresh Token
		encryptedRefresh, err := middleware.Encrypt(refreshToken)
		if err != nil {
			return nil, err
		}
	
		//9. encode base64 
		encodedToken := base64.URLEncoding.EncodeToString([]byte(encryptedRefresh))
	
		//10. time exp refresh token
		expiresAt := time.Now().Add(7 * 24 * time.Hour)
	
		//11. Send refresh token to database
		if err := s.Ref.RefreshToken(user.ID, encodedToken, expiresAt); err != nil {
			return nil, err
		}
	
		return map[string]interface{}{
			"message": "login success",
			"access_token":   accessToken,
			"refresh_token": encodedToken,
		}, nil
	}
	// 5. Simpan aktivitas login
	err = s.Repo.SaveLoginHistory(user.ID, deviceInfo, ip)
	if err != nil {
		return nil, err
	}
	//6. Generate Access token
	accessToken, err := s.JWTSecret.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	
	//7. Generate Referesh Token
	refreshToken, err := s.JWTSecret.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}
	
	//8. encryp Refresh Token
	encryptedRefresh, err := middleware.Encrypt(refreshToken)
	if err != nil {
		return nil, err
	}
	
	//9. encode base64 
	encodedToken := base64.URLEncoding.EncodeToString([]byte(encryptedRefresh))
	
	//10. time exp refresh token
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	
	//11. Send refresh token to database
	if err := s.Ref.UpdateRefreshToken(user.ID, encodedToken, expiresAt); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"message": "login success",
		"access_token":   accessToken,
		"refresh_token": encodedToken,
		}, nil
}
	
	// -------------------------- REFRESH TOKEN ------------------------
func (s *AuthService) RefreshToken(encryptedToken string) (map[string]interface{}, error) {
	
	
		decodedBytes, err := base64.URLEncoding.DecodeString(encryptedToken)
		
		if err != nil {
			return nil, ErrInvalidToken
		}
	
		refreshToken, err := middleware.Decrypt(string(decodedBytes))
		if err != nil {
			return nil, ErrInvalidToken
		}
	
		// parsing token
		token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
			return s.JWTSecret.RefreshSecret, nil
		})
			
		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))
	
	
		tokenCheck, err := s.Ref.FindRefreshToken(userID)
	
	
		if err != nil || tokenCheck == nil {
			return nil, errors.New("refresh token not found")
		}
	
		if encryptedToken != tokenCheck.RefreshToken {
			return nil, errors.New("refresh token not match")
		} 
	
		// ambil user dari db
		user, err := s.Repo.FindByID(strconv.Itoa(userID))
		if err != nil || user == nil {
			return nil, err
		}
	
		// generate token baru (access + refresh)
		accessToken, err := s.JWTSecret.GenerateAccessToken(user.ID, user.Email, user.Role)
		if err != nil {
			return nil, err
		}
	
		newRefreshToken, err := s.JWTSecret.GenerateRefreshToken(user.ID)
		if err != nil {
			return nil, err
		}
		encryptedRefresh, err := middleware.Encrypt(newRefreshToken)
		if err != nil {
			return nil, err
		}
		encodedToken := base64.URLEncoding.EncodeToString([]byte(encryptedRefresh))
	
		expiresAt := time.Now().Add(7 * 24 * time.Hour)
	
		err = s.Ref.UpdateRefreshToken(user.ID, encodedToken, expiresAt)
		if err != nil {
			return nil, err
		}
	
		return map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": encodedToken,
		}, nil
	}

func (s *AuthService) Logout(userID int) error {
	if err := s.Repo.UpdateLoginStatus(userID, 0); err != nil {
		return err
	}
	if err := s.Ref.DeleteRefreshToken(userID); err != nil {
		return err
	}
	return nil
}