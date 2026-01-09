package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
	"github.com/muhammadfarrasfajri/koperasi-gerai/services"
)

type UserAuthController struct {
	AuthService *services.UserAuthService
}

func NewAuthController(userAuthService *services.UserAuthService) *UserAuthController {
	return &UserAuthController{
		AuthService: userAuthService,
	}
}

func (c *UserAuthController) RegisterUser(ctx *gin.Context) {

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
        user.KtpPicture = path // Masukkan path ke struct user yang sudah terisi tadi
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