package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/models"
	"github.com/muhammadfarrasfajri/login-google/services"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authservice *services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authservice,
	}
}

func (c *AuthController) RegisterUser(ctx *gin.Context) {

    var user models.BaseUser

    // 1. AUTO BINDING (Pengganti manual satu-satu)
    if err := ctx.ShouldBind(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Ambil ID Token manual karena dia tidak masuk struct BaseUser
    idToken := ctx.PostForm("id_token")
    if idToken == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "id_token wajib diisi"})
        return
    }

    // 2. PROSES UPLOAD FILE (Tetap manual karena perlu Logic Simpan & Rename)
    fileKTP, err := ctx.FormFile("ktp_image")
    if err == nil {
        path := "public/uploads/ktp/" + fmt.Sprintf("%d_%s", time.Now().Unix(), fileKTP.Filename)
        ctx.SaveUploadedFile(fileKTP, path)
        user.KtpImagePath = path // Masukkan path ke struct user yang sudah terisi tadi
    }

    fileProfile, err := ctx.FormFile("profile_image")
    if err == nil {
        path := "public/uploads/profile/" + fmt.Sprintf("%d_%s", time.Now().Unix(), fileProfile.Filename)
        ctx.SaveUploadedFile(fileProfile, path)
        user.ProfilePicture = path
    }

	//id register
	ip := ctx.ClientIP()
	user.RegisterIP = ip
	
    // 3. Panggil Service
    resUser, err := c.AuthService.Register(idToken, user)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Register success", "user": resUser})
}

func (c *AuthController) LoginUser(ctx *gin.Context) {
	var req struct {
		IDToken    string `json:"id_token"`
		DeviceInfo string `json:"device_info"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ip := ctx.ClientIP()

	result, err := c.AuthService.Login(req.IDToken, req.DeviceInfo, ip)
	if err != nil {
		ctx.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *AuthController) RefreshTokenUser(ctx *gin.Context) {

	refreshToken, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	result, err := c.AuthService.RefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *AuthController) LogoutUser(ctx *gin.Context) {

	userID := ctx.GetInt("user_id")

	err := c.AuthService.Logout(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logout success"})
}