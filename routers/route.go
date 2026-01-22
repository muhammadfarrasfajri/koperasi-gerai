package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai/controllers"
)

func SetupRoutes(r *gin.Engine, userAuthController *controllers.UserAuthController) {

	// ===========================
	// AUTH ROUTES
	// ===========================
	auth := r.Group("/api/auth")
	{
		//auth User
		auth.POST("/user/register", userAuthController.RegisterUser)
		auth.POST("/user/login")
	}
}
