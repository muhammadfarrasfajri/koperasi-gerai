package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/repository"
	"github.com/muhammadfarrasfajri/login-google/services"
)

type UserController struct {
	UserService *services.UserService
	Repo        *repository.UserRepository
}

func NewUserController(userService *services.UserService, repo *repository.UserRepository) *UserController{
	return &UserController{
		UserService: userService,
		Repo: repo,
	}
}

// GET /users
func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.UserService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get users",
		"users":   users,
	})
}

// GET /users/:id
func (c *UserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.UserService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get user",
		"user":    user,
	})
}

// PUT /users/:id
// PUT /users/:id
func (c *UserController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	// Ambil form-data
	name := ctx.PostForm("Name")
	email := ctx.PostForm("Email")
	role := ctx.PostForm("Role")
	fmt.Println(id, name, email, role)

	// Ambil file kalau ada
	file, _ := ctx.FormFile("Profile_picture")
	var filename string
	var publicPath string
	if file != nil {
		filename = fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		savePath := fmt.Sprintf("./public/uploads/images/%s", filename)
		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}
		publicPath = fmt.Sprintf("/public/uploads/images/%s", filename) // ini yang dikirim ke DB
	}

	// Kirim ke service/repo
	user, err := c.UserService.Update(id, name, email, role, publicPath)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format URL fotoqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq
	if user.ProfilePicture != "" && !strings.HasPrefix(user.ProfilePicture, "http") {
		user.ProfilePicture = fmt.Sprintf(user.ProfilePicture)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update success",
		"user":    user,
	})
}

// DELETE /users/:id
func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.UserService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete success",
	})
}

func (uc *UserController) UploadPhoto(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Buat nama file unik
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	savePath := fmt.Sprintf("public/uploads/images/%s", filename)

	// Simpan file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// URL publik
	publicURL := fmt.Sprintf("/public/uploads/images/%s", filename)

	// Simpan ke database via repository (contoh)
	userID, err := strconv.Atoi(c.PostForm("user_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "user_id must be an integer"})
		return
	}

	if err := uc.Repo.UpdatePhotoURL(userID, publicURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL in DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Photo uploaded successfully",
		"url":     publicURL,
	})
}
