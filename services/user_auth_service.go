package services

import (
	"context"
	"errors"

	firebase "firebase.google.com/go/auth"
	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
	"github.com/muhammadfarrasfajri/koperasi-gerai/repository"
)

var (
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrUserNotRegistered = errors.New("user not registered, please register first")
	ExistingNIK = errors.New("NIK already usage")
)

type UserAuthService struct {
	AuthRepo 	repository.AuthRepository
	UserRepo 	repository.UserRepository
	FirebaseAuth *firebase.Client
}

func NewUserAuthService(userAuthRepository repository.AuthRepository, userRepository repository.UserRepository, firebaseAuth *firebase.Client) *UserAuthService{
	return &UserAuthService{
		AuthRepo: userAuthRepository,
		UserRepo: userRepository,
		FirebaseAuth: firebaseAuth,
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

	// 2. Cek apakah user sudah ada
	existing, err := s.UserRepo.IsNIKExists(user.NIK)
	if err != nil {
		return nil, err
	}

	if existing {
		return nil, ExistingNIK
	}

	// 3. Tentukan nama
	if user.Name == "" {
        if n, ok := token.Claims["name"].(string); ok {
            user.Name = n
        }
    }

	length := len(user.PhoneNumber)
	
	if length < 11 || length > 13 {
		return nil, errors.New("nomor HP harus 11–13 digit")
	}

	err = s.AuthRepo.CreateRegisterUser(user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}