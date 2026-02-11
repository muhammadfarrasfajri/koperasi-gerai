package services

import "github.com/muhammadfarrasfajri/koperasi-gerai/repository"

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService{
	return &UserService{
		UserRepo: userRepository,
	}
}