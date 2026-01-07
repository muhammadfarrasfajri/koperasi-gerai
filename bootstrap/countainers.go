package bootstrap

import (
	"os"

	"firebase.google.com/go/auth"
	"github.com/muhammadfarrasfajri/login-google/controllers"
	"github.com/muhammadfarrasfajri/login-google/database"
	"github.com/muhammadfarrasfajri/login-google/middleware"
	"github.com/muhammadfarrasfajri/login-google/repository"
	"github.com/muhammadfarrasfajri/login-google/services"
)

type Container struct {
	AuthAdminController *controllers.AuthController
	AuthUserController  *controllers.AuthController
	UserController      *controllers.UserController
	JWTManager          *middleware.JWTManager
}

func InitContainer(adminAuth, userAuth *auth.Client) *Container {
	adminRepo := repository.NewAdminRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	adminRef := repository.NewRefreshTokenAdmin(database.DB)
	userRef := repository.NewRefreshTokenUser(database.DB)

	jwtManager := middleware.NewJWTManager(
		os.Getenv("JWT_SECRET"),
		os.Getenv("REFRESH_SECRET"),
	)

	authAdminService := services.NewAuthService(adminRepo, adminRef, adminAuth, jwtManager)
	authUserService := services.NewAuthService(userRepo, userRef, userAuth, jwtManager)
	userService := services.NewUserSevice(userRepo)

	return &Container{
		AuthAdminController: controllers.NewAuthController(authAdminService),
		AuthUserController:  controllers.NewAuthController(authUserService),
		UserController:      controllers.NewUserController(userService, userRepo),
		JWTManager:          jwtManager,
	}
}
