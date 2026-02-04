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

    if err := ctx.ShouldBind(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, models.APIResponse{
            Error: true,
            Message: "Invalid data format",
            Type: "ValidationError",
        })
        return
    }

    idToken := ctx.PostForm("id_token")
    if idToken == "" {
        ctx.JSON(http.StatusBadRequest, models.APIResponse{
            Error: true,
            Message: "Failed to send id token",
            Type: "TokenError",
        })      
        return
    }

    fileKTP, err := ctx.FormFile("ktp_picture")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, models.APIResponse{
            Error: true,
            Message: "Failed to send KTP photo",
            Type: "KTPError",
        })
        return
    }
    pathKtp := "public/uploads/ktp/" + fmt.Sprintf("%d_%s", time.Now().Unix(), fileKTP.Filename)
    if err := ctx.SaveUploadedFile(fileKTP, pathKtp); err != nil {
         ctx.JSON(http.StatusInternalServerError, models.APIResponse{
            Error: true,
            Message: "Failed to save KTP photo",
            Type: "ServerError",
         })
         return
    }
    user.KtpPicture = pathKtp 

    // -- PROFILE PICTURE --
    fileProfile, err := ctx.FormFile("profile_picture")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, models.APIResponse{
            Error: true,
            Message: "Failed to send profile photo",
            Type: "ProfileError",
        })
        return
    }
    
    pathProfile := "public/uploads/profile/" + fmt.Sprintf("%d_%s", time.Now().Unix(), fileProfile.Filename)
    if err := ctx.SaveUploadedFile(fileProfile, pathProfile); err != nil {
         ctx.JSON(http.StatusInternalServerError, models.APIResponse{
            Error: true,
            Message: "Failed to save profile photo",
            Type: "ServerError",
         })
         return
    }
    user.ProfilePicture = pathProfile

    // IP Address
    user.RegisterIP = ctx.ClientIP()
    
    resUser, err := c.AuthService.Register(idToken, user)
    
    if err != nil {
        ctx.JSON(http.StatusBadRequest, models.APIResponse{
            Error: true,
            Message: "Register Failed",
            Type: "Register User",
        })
        return
    }

    ctx.JSON(http.StatusOK, models.APIResponse{
        Error: false,
        Message: "Register Success",
        Data: resUser,
    })
}

func (c *UserAuthController) LoginUser(ctx *gin.Context){
    loginRequest := models.LoginRequest{}
    
    err := ctx.BindJSON(&loginRequest)

    if err != nil {
        ctx.JSON(http.StatusBadRequest, models.APIResponse{
            Error: true,
            Message: "Invalid data format",
            Type: "ValidationError",
        })
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
		ctx.JSON((http.StatusBadRequest), models.APIResponse{
            Error: true,
            Message: err.Error(),
            Type: "AuthenticationError",
        })
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
        Error:   false,
        Message: "Login success",
        Data:    result,
    })
}

func (c *UserAuthController) RefreshToken(ctx *gin.Context) {

    type RefreshTokenReq struct {
        RefreshToken string `json:"token_hash" binding:"required"`
    }

    var req RefreshTokenReq

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, models.APIResponse{
            Error: true,
            Message: "Invalid data format",
            Type: "ValidationError",
        })
        return
    }

    result, err := c.AuthService.RefreshToken(req.RefreshToken)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, models.APIResponse{
            Error: true,
            Message: err.Error(),
            Type: "RefreshTokenError",
        })
        return
    }
    
    ctx.JSON(http.StatusOK, models.APIResponse{
        Error: false,
        Message: "Generate refresh token success",
        Data: result,
    })
}

func (c *UserAuthController) LogoutUser(ctx *gin.Context) {
    // Ambil refresh token dari body request
    type LogoutReq struct {
        RefreshToken string `json:"token_hash" binding:"required"`
    }

    var req LogoutReq

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, models.APIResponse{
            Error: true,
            Message: "Invalid data format",
            Type: "ValidationError",
        })
        return
    }

    err := c.AuthService.Logout(req.RefreshToken)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, models.APIResponse{
            Error: true,
            Message: "Logout Failed",
            Type: "LogoutError",
        })
        return
    }

    ctx.JSON(http.StatusOK, models.APIResponse{
        Error: false,
        Message: "Logout Success",
    })
}