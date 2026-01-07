package services

import (
	"errors"

	"github.com/muhammadfarrasfajri/login-google/models"
	"github.com/muhammadfarrasfajri/login-google/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserSevice(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

// ------------------------- GET ALL USERS -----------------------------

func (s *UserService) GetAll() ([]models.BaseUser, error) {
	users, err := s.UserRepo.GetAll()
	if err != nil {
		return nil, ErrUserNotFound
	}
	return users, nil
}

// ------------------------- GET USER BY ID ----------------------------

func (s *UserService) GetByID(id string) (*models.BaseUser, error) {
	user, err := s.UserRepo.FindByID(id)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// --------------------------- UPDATE USER -----------------------------

func (s *UserService) Update(id, name, email, role, ProfilePicture string) (*models.BaseUser, error) {
	// cek apakah user ada
	existing, err := s.UserRepo.FindByID(id)
	if err != nil || existing == nil {
		return nil, ErrUserNotFound
	}

	// update field
	existing.Name = name
	existing.Email = email
	existing.Role = role
	existing.ProfilePicture = ProfilePicture

	err = s.UserRepo.Update(*existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// --------------------------- DELETE USER -----------------------------

func (s *UserService) Delete(id string) error {
	// cek user dulu
	user, err := s.UserRepo.FindByID(id)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	return s.UserRepo.Delete(id)
}
