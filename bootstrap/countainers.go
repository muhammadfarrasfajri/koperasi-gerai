package bootstrap

import (
	"firebase.google.com/go/auth"
	"github.com/muhammadfarrasfajri/koperasi-gerai/controllers"
	"github.com/muhammadfarrasfajri/koperasi-gerai/database"
	"github.com/muhammadfarrasfajri/koperasi-gerai/middleware"
	"github.com/muhammadfarrasfajri/koperasi-gerai/repository"
	"github.com/muhammadfarrasfajri/koperasi-gerai/services"
)

type Container struct {
	UserAuthController  *controllers.UserAuthController
	JWTManager *middleware.JWTManager
}

func InitContainer(userAuth *auth.Client) *Container {
	userAuthRepo := repository.NewUserAuthRepo(database.DB)
	userRepo := repository.NewUserRepo(database.DB)
	userRefRepo := repository.NewUserRefreshTokenRepo(database.DB)

	jwtManager := middleware.NewJWTManager("79329e633bbbd5652893feea5c27f60faa0ee69688e65e29bf03419889be965adcab16420e07fa88c62ab8d1f7c82804aee66e30d237f6381b002e1ae1109187","aaa46ab1983939bfaa571d4b6581e2012d0cec9a67ee1cec975f64af716f6850f080d389b12b70b6918e94bf417c7448133cbd129c8c4bf567bdb7f82bbfa3a1")

	userAuthService := services.NewUserAuthService(userAuthRepo, userRepo, userRefRepo, userAuth, jwtManager)
	
	return &Container{
		UserAuthController: controllers.NewAuthController(userAuthService),
		JWTManager: jwtManager,
	}
}
