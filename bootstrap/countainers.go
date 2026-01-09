package bootstrap

import (
	"firebase.google.com/go/auth"
	"github.com/muhammadfarrasfajri/koperasi-gerai/controllers"
	"github.com/muhammadfarrasfajri/koperasi-gerai/database"
	"github.com/muhammadfarrasfajri/koperasi-gerai/repository"
	"github.com/muhammadfarrasfajri/koperasi-gerai/services"
)

type Container struct {
	UserAuthController  *controllers.UserAuthController
}

func InitContainer(userAuth *auth.Client) *Container {
	userAuthRepo := repository.NewUserAuthRepo(database.DB)
	userRepo := repository.NewUserRepo(database.DB)


	userAuthService := services.NewUserAuthService(userAuthRepo, userRepo, userAuth)
	
	return &Container{
		UserAuthController: controllers.NewAuthController(userAuthService),
	}
}
