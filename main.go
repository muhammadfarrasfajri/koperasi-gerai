package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai/bootstrap"
	"github.com/muhammadfarrasfajri/koperasi-gerai/middleware"
	routes "github.com/muhammadfarrasfajri/koperasi-gerai/routers"
)

func main() {


	// Database
	bootstrap.InitDatabase()

	// Firebase
	userAuth := bootstrap.InitFirebase()

	// Build container (repositories, services, controllers)
	container := bootstrap.InitContainer(userAuth)

	// GIN
	r := gin.Default()
	
	r.Static("/public", "./public")

	r.MaxMultipartMemory = 50 << 20 // 50 MB

	// CORS Middleware
	middleware.AttachCORS(r)

	// ROUTES
	routes.SetupRoutes(r, container.UserAuthController, container.JWTManager)

	r.Run(":8080")
}
