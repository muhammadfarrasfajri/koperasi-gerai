package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/bootstrap"
	"github.com/muhammadfarrasfajri/login-google/middleware"
	routes "github.com/muhammadfarrasfajri/login-google/routers"
)

func main() {

	// ENV
	bootstrap.InitEnv()

	// Encryption Key
	middleware.InitEncryptionKey()

	// Database
	bootstrap.InitDatabase()

	// Firebase
	adminAuth, userAuth := bootstrap.InitFirebase()

	// Build container (repositories, services, controllers)
	container := bootstrap.InitContainer(adminAuth, userAuth)

	// GIN
	r := gin.Default()
	r.Static("/public", "./public")

	// CORS Middleware
	middleware.AttachCORS(r)

	// ROUTES
	routes.SetupRoutes(
		r,
		container.AuthAdminController,
		container.AuthUserController,
		container.UserController,
		container.JWTManager,
	)

	r.Run(":8080")
}
