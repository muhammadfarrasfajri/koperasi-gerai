package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai/controllers"
	"github.com/muhammadfarrasfajri/koperasi-gerai/middleware"
)

func SetupRoutes(r *gin.Engine, userAuthController *controllers.UserAuthController, jwtManager *middleware.JWTManager) {

	// ===========================
	// AUTH ROUTES
	// ===========================
	auth := r.Group("/api/auth")
	{
		//auth User
		auth.POST("/user/register", userAuthController.RegisterUser)
		auth.POST("/user/login", userAuthController.LoginUser)
		auth.POST("/user/refresh", userAuthController.RefreshToken)
		auth.POST("/user/logout", jwtManager.AuthMiddleware(), userAuthController.LogoutUser)
	}
}
