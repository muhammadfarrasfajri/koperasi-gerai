package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
	"github.com/muhammadfarrasfajri/koperasi-gerai/services"
	"github.com/muhammadfarrasfajri/koperasi-gerai/utils"
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

    // 1. AUTO BINDING
    if err := ctx.ShouldBind(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Ambil ID Token
    idToken := ctx.PostForm("id_token")
    if idToken == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "id_token wajib diisi"})
        return
    }
    // -- KTP --
    fileKTP, err := ctx.FormFile("ktp_picture")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wajib upload KTP (key: ktp_picture)"})
        return
    }
    pathKtp := "public/uploads/ktp/" + fmt.Sprintf("%d_%s", time.Now().Unix(), fileKTP.Filename)
    if err := ctx.SaveUploadedFile(fileKTP, pathKtp); err != nil {
         ctx.JSON(500, gin.H{"error": "Gagal save KTP"})
         return
    }
    user.KtpPicture = pathKtp 

    // -- PROFILE PICTURE --
    fileProfile, err := ctx.FormFile("profile_picture")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wajib upload Profile Picture"})
        return
    }
    pathProfile := "public/uploads/profile/" + fmt.Sprintf("%d_%s", time.Now().Unix(), fileProfile.Filename)
    if err := ctx.SaveUploadedFile(fileProfile, pathProfile); err != nil {
         ctx.JSON(500, gin.H{"error": "Gagal save Profile Picture"})
         return
    }
    user.ProfilePicture = pathProfile

    // IP Address
    user.RegisterIP = ctx.ClientIP()
    
    resUser, err := c.AuthService.Register(idToken, user)
    
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
            "message": "Registrasi Failed",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Register Success",
        "data": gin.H{
            "id_member": resUser.IDMember,
            "name":      resUser.Name,
            "email":     resUser.Email,
            "role":      resUser.ActiveAs,
        },
    })
}

func (c *UserAuthController) LoginUser(ctx *gin.Context){
    loginRequest := models.LoginRequest{}
    
    err := ctx.BindJSON(&loginRequest)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
    }
    
    uaRaw := ctx.Request.UserAgent()
    if uaRaw == ""{
        uaRaw = "unknown"
    }
    
    deviceInfo := utils.ParseDeviceInfo(uaRaw)

    loginHistory := models.BaseLoginHistory{
        IPAddress: ctx.ClientIP(),
        UserAgent: uaRaw,
        DeviceInfo: deviceInfo,
        Location: loginRequest.Location,
    }

	result, err := c.AuthService.Login(loginRequest.IdToken, loginHistory)
	if err != nil {
		ctx.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}
    
	ctx.JSON(http.StatusOK, result)
}

func (c *UserAuthController) RefreshToken(ctx *gin.Context) {

    type RefreshTokenReq struct {
        RefreshToken string `json:"token_hash" binding:"required"`
    }

    var req RefreshTokenReq

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token wajib dikirim dalam body JSON"})
        return
    }

    result, err := c.AuthService.RefreshToken(req.RefreshToken)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message":       "success",
        "access_token":  result["access_token"],  // Token buat akses API
        "refresh_token": result["refresh_token"], // Token buat diputar (Rotation)
    })
}

func (c *UserAuthController) LogoutUser(ctx *gin.Context) {
    // Ambil refresh token dari body request
    type LogoutReq struct {
        RefreshToken string `json:"token_hash" binding:"required"`
    }

    var req LogoutReq

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"error": "Refresh token wajib dikirim"})
        return
    }

    // Panggil Service
    err := c.AuthService.Logout(req.RefreshToken)
    if err != nil {
        ctx.JSON(500, gin.H{"error": "Gagal logout"})
        return
    }

    ctx.JSON(200, gin.H{"message": "Berhasil logout"})
}