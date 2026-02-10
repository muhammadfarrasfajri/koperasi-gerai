package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	firebase "firebase.google.com/go/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/muhammadfarrasfajri/koperasi-gerai/middleware"
	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
	"github.com/muhammadfarrasfajri/koperasi-gerai/repository"
	"github.com/muhammadfarrasfajri/koperasi-gerai/utils"
)

var (
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrUserNotRegistered = errors.New("user not registered, please register first")
	ErrExistingNIK = errors.New("NIK already usage")
	ErrInvalidNIK =  errors.New("Invalid NIK")
	ErrInvalidNoHp = errors.New("Mobile number must consist of 11-13 digits")
	ErrDuplicateUser = errors.New("Duplicate user")
	ErrExistingPhoneNo = errors.New("Phone number already usage")
    ErrGenMemberId = errors.New("Generate member id failed")
    ErrInvalidPoscode = errors.New("Invalid pos code")
    ErrRegisterFailed = errors.New("Register Failed")
    ErrInvalidUser = errors.New("invalid user data")
    ErrUserNotVerified = errors.New("User is not verified")
    ErrAccessToken = errors.New("Failed generate JWT access")
    ErrRefreshToken = errors.New("Failed generate JWT refresh")
    ErrloginSession= errors.New("Failed to save login session")
    ErrParseToken = errors.New("unexpected signing method")
    ErrInvalidRefreshToken = errors.New("refresh token invalid or expired")
    ErrRefreshTokenNotFound = errors.New("Refresh token not found in db")
    ErrRefreshTokenNotMatch = errors.New("refresh token reuse detected / not match")
    ErrUserNotFound = errors.New("User not found")
    ErrSaveNewSession = errors.New("failed to save new session")
    ErrLogout = errors.New("Logout Failed / token not found")
)

type UserAuthService struct {
	AuthRepo repository.AuthRepository
	UserRepo repository.UserRepository
	RefRepo repository.UserRefreshTokenRepository
	FirebaseAuth *firebase.Client				
	JWTSecret    *middleware.JWTManager
}

func NewUserAuthService(userAuthRepository repository.AuthRepository, userRepository repository.UserRepository, refreshRepository repository.UserRefreshTokenRepository, firebaseAuth *firebase.Client, jwtsecret *middleware.JWTManager) *UserAuthService{
	return &UserAuthService{
		AuthRepo: userAuthRepository,
		UserRepo: userRepository,
		RefRepo: refreshRepository,
		FirebaseAuth: firebaseAuth,
		JWTSecret: jwtsecret,
	}
}

// --------------------------- REGISTER -----------------------------------

func (s *UserAuthService) Register(idToken string, user models.BaseUser) (err error) {
    ctx := context.Background()

    logError := func(err error, context string) {
        prettyJSON, _ := json.MarshalIndent(user, "", "  ")
       
        shortToken := ""
        if len(idToken) > 10 { shortToken = idToken[:10] + "..." }

        log.Printf("[ERROR] %s | Token: %s | Err: %v\n", 
        context, shortToken, err)
        fmt.Printf("\n[DEBUG DATA USER]:\n%s\n\n", string(prettyJSON))
    }
    // checking cell phone numbers
    if len(user.PhoneNumber) < 11 || len(user.PhoneNumber) > 15 {
        logError(ErrInvalidNoHp, "Phone number validation")
        return  ErrInvalidNoHp
    }
    // checking NIK
    if len(user.NIK) != 16 {
        logError(ErrExistingNIK, "NIK validation")
        return  ErrInvalidNIK
    }
    //checking pos code
    if len(user.PosCode) != 5 {
        logError(ErrInvalidPoscode, "Pos code validation")
        return  ErrInvalidPoscode
    }
    //checking NPWP
    isValid, message := utils.ValidateNPWP(user.NPWP)
    if isValid == false || message != "" {
        logError(errors.New(message), "NPWP validation")
        return  errors.New(message)
    }

    // Verification Token Firebase
    token, err := s.FirebaseAuth.VerifyIDToken(ctx, idToken)
    if err != nil {
        logError(ErrInvalidToken, "Firebase token verification")
        return  ErrInvalidToken
    }

    // Get data from uid Google
    user.GoogleUID = token.UID

    if email, ok := token.Claims["email"].(string); ok {
        user.Email = email
    }

    if pic, ok := token.Claims["picture"].(string); ok {
        user.GooglePicture = pic
    }

    if user.Name == "" {
        if n, ok := token.Claims["name"].(string); ok {
            user.Name = n
        }
    }

    // Get user by email
    existingUser, err := s.AuthRepo.FindByEmail(user.Email)

    if err == nil && existingUser != nil {
        if existingUser.GoogleUID == "" || existingUser.GoogleUID != user.GoogleUID {
            errLink := s.AuthRepo.LinkGoogleAccount(user.Email, user.GoogleUID, user.GooglePicture)
            if errLink != nil {
                return  fmt.Errorf("gagal menghubungkan akun: %v", errLink)
            }
            user.IDMember = existingUser.IDMember
        }
        return nil
    }
    // checking NIk availability
    nikExists, _ := s.AuthRepo.IsNIKExists(user.NIK)
    if nikExists {
        logError(ErrExistingNIK, "Checking exists NIK")
        return ErrExistingNIK
    }

    // checking phone number availability
    hpExists, _ := s.AuthRepo.IsNoHPExists(user.PhoneNumber)
    if hpExists {
        logError(ErrExistingPhoneNo, "Checking exists phone number")
        return ErrExistingPhoneNo
    }

    currentYear := time.Now().Format("06")

    prefix := fmt.Sprintf("KOP-KF-%s", currentYear)
   
    lastID, err := s.AuthRepo.GetMemberId(prefix)
    if err != nil {
        logError(err, "Get last id member id")
        return err 
    }
    
    newNumber := 1
    if lastID != "" {
        parts := strings.Split(lastID, "-")
        if len(parts) > 0 {
            lastNumberStr := parts[len(parts)-1]
            if num, errAtoi := strconv.Atoi(lastNumberStr); errAtoi == nil {
                newNumber = num + 1
            }
        }
    }

    user.IDMember = fmt.Sprintf("%s-%010d", prefix, newNumber)

    err = s.AuthRepo.CreateRegisterUser(user)
    if err != nil {
        logError(ErrRegisterFailed, "Create User register")
        return ErrRegisterFailed
    }

    fmt.Println("Register Success")
    prettyJSON, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println("Data user : ", string(prettyJSON))

    return nil
}

func (s *UserAuthService) Login(idToken string, loginHistory models.BaseLoginHistory) (map[string]interface{}, error) {

    ctx := context.Background()
    
    var savedUser *models.BaseUser 
    finalStatus := "FAILED"
    emailInfo := ""

    defer func() {
        if savedUser != nil {
            log.Printf("[DEFER CHECK] User ID: %d | Status: %s\n", savedUser.ID, finalStatus)
            
            // Masukkan data ke history
            loginHistory.UserID = savedUser.ID
            loginHistory.Status = finalStatus

            if loginHistory.LoginAt.IsZero() {
                loginHistory.LoginAt = time.Now().UTC()
            }
            if errHist := s.AuthRepo.HistoryLoginUser(loginHistory); errHist != nil {
                log.Printf("[ERROR DB] Gagal simpan history: %v\n", errHist)
            } else {
                log.Println("[SUCCESS DB] History login berhasil disimpan.")
            }
        } else {
            log.Println("[DEFER SKIP] History tidak disimpan karena savedUser masih NIL.")
        }
    }()

    logError := func(err error, context string) {

        prettyJSON, _ := json.MarshalIndent(loginHistory, "", "  ")
        emailFinal := emailInfo
        shortToken := ""
        if len(idToken) > 10 { shortToken = idToken[:10] + "..." }

        log.Printf("[ERROR] %s | Token: %s | Err: %v\n", 
        context, shortToken, err)
        log.Println(emailFinal)
        fmt.Printf("\n[DEBUG DATA LOGIN]:\n%s\n\n", string(prettyJSON))
    }

    token, err := s.FirebaseAuth.VerifyIDToken(ctx, idToken)
    if err != nil {
        logError(ErrInvalidToken, "verifikasi firebase")
        return nil, ErrInvalidToken
    }

    user, err := s.UserRepo.FindByGoogleUID(token.UID)
    if err != nil { 
        email, _ := token.Claims["email"].(string)
        existingUser, errEmail := s.AuthRepo.FindByEmail(email)
        if errEmail == nil && existingUser != nil {
            pic, _ := token.Claims["picture"].(string)
            _ = s.AuthRepo.LinkGoogleAccount(email, token.UID, pic)
            user, err = s.UserRepo.FindByGoogleUID(token.UID)
            if err != nil {
                emailInfo = email
                logError(errors.New("gagal mengambil data user setelah linking"), "Find By Google UID")
                return nil, errors.New("gagal mengambil data user setelah linking")
            }
        } else {
            emailInfo = email
            logError(errors.New("User not found. Please Register first."), "User not found")
            return nil, errors.New("User not found. Please Register first.")
        }
    }
    emailInfo = user.Email
    savedUser = user

    if user == nil {
        logError(ErrInvalidUser, "invalid user")
        return nil, ErrInvalidUser
    }

    finalStatus = "FAILED" 

    // 3. Cek Status Verifikasi
    if user.Is_verified == 0 {
        loginHistory.ErrorMessage = "User is not verified"
        logError(ErrUserNotVerified, "Checking status verifikasi")
        return nil, ErrUserNotVerified
    }

    // 4. Generate Access Token
    accessToken, err := s.JWTSecret.GenerateAccessToken(user.ID, user.Email)
    if err != nil {
        logError(ErrAccessToken, "Generate access token")
        return nil, ErrAccessToken
    }

    // 5. Generate Refresh Token
    refreshToken, err := s.JWTSecret.GenerateRefreshToken(user.ID)
    if err != nil {
         logError(ErrRefreshToken, "Generate refresh token")
        return nil, ErrRefreshToken
    }

    // 6. Hash & Siapkan Model
    refreshTokenHash := middleware.HashToken(refreshToken)
    expiresAt := time.Now().Add(7 * 24 * time.Hour) 

    tokenModel := models.RefreshToken{
        UserID:    user.ID,
        TokenHash: refreshTokenHash,
        ExpiresAt: expiresAt,
    }

    // 7. Simpan ke Database (Upsert)
    err = s.RefRepo.UpsertRefreshToken(tokenModel)
    if err != nil {
        logError(ErrloginSession, "upsert history login")
        return nil, ErrloginSession
    }

    loginHistory.LoginAt = time.Now().UTC()

    finalStatus = "SUCCESS"
    loginHistory.Status = finalStatus
    loginHistory.UserID = user.ID   
    loginHistory.Status = finalStatus

    //LOG User
    prettyHitory, _ := json.MarshalIndent(loginHistory, "", "  ")
    prettyToken, _ := json.MarshalIndent(tokenModel, "", " ")
    fmt.Println(string(prettyToken))
    fmt.Println(string(prettyHitory))
    log.Println(emailInfo)

    return map[string]interface{}{
        "access_token":  accessToken,
        "token_hash": refreshToken,
        "user": map[string]interface{}{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
            "role":  user.ActiveAs,
        },
    }, nil
}

func (s *UserAuthService) RefreshToken(rawRefreshToken string) (map[string]interface{}, error) {

     logError := func(err error, context string) {
        log.Printf("[ERROR] %s | RefreshToken: %s | Err: %v\n", 
        context, rawRefreshToken, err)
    }

    token, err := jwt.Parse(rawRefreshToken, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            logError(ErrParseToken, "parse Token")
            return nil, ErrParseToken
        }
        return s.JWTSecret.RefreshSecret, nil
    })

    if err != nil {
        logError(ErrInvalidRefreshToken, "Checking refresh token")
        return nil, ErrInvalidRefreshToken
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        logError(errors.New("Invalid token claims"), "get claims")
        return nil, errors.New("invalid token claims")
    }

    userIDFloat, ok := claims["user_id"].(float64)
    if !ok {
        logError(errors.New("Invalid user id in token"), "get user id")
        return nil, errors.New("invalid user id in token")
    }

    userID := int(userIDFloat)

    // 2. Cek Token di Database
    tokenCheck, err := s.RefRepo.FindRefreshTokenUser(userID)
    if err != nil || tokenCheck == nil {
        logError(ErrRefreshTokenNotFound, "Checking refresh token in database")
        return nil, ErrInvalidRefreshToken
    }

    // PERBAIKAN 1: Hash dulu token dari user, baru bandingkan
    incomingTokenHash := middleware.HashToken(rawRefreshToken)
    
    if incomingTokenHash != tokenCheck.TokenHash {
        logError(ErrRefreshTokenNotMatch, "compare refresh token")
        return nil,ErrRefreshTokenNotMatch
    }

    // 3. Ambil User untuk Data Token Baru
    user, err := s.UserRepo.FindById(strconv.Itoa(userID))
    if err != nil || user == nil {
        logError(ErrUserNotFound, "get user with id")
        return nil, ErrUserNotFound
    }

    // 4. Generate Token Baru
    newAccessToken, err := s.JWTSecret.GenerateAccessToken(user.ID, user.Email)
    if err != nil {
        logError(ErrAccessToken, "Generate access token")
        return nil, err
    }

    newRefreshToken, err := s.JWTSecret.GenerateRefreshToken(user.ID)
    if err != nil {
        logError(ErrRefreshToken, "Generate refresh token")
        return nil, err
    }

    // 5. Simpan Token Baru ke Database
    newRefreshTokenHash := middleware.HashToken(newRefreshToken)

    expiresAt := time.Now().Add(7 * 24 * time.Hour)

    tokenModel := models.RefreshToken{
        UserID:    user.ID,    
        TokenHash: newRefreshTokenHash,
        ExpiresAt: expiresAt,
    }

    err = s.RefRepo.UpsertRefreshToken(tokenModel)
    if err != nil {
        logError(ErrSaveNewSession, "Update refresh token")
        return nil, ErrSaveNewSession
    }

    prettyToken, _ := json.MarshalIndent(tokenModel, "", " ")
    fmt.Println(string(prettyToken))

    return map[string]interface{}{
        "access_token":  newAccessToken,
        "token_hash": newRefreshToken,
    }, nil
}

func (s *UserAuthService) Logout(rawRefreshToken string) error {
    
    logError := func(err error, context string) {

        log.Printf("[ERROR] %s | Token: %s | Err: %v\n", 
        context, rawRefreshToken, err)
    }
    tokenHash := middleware.HashToken(rawRefreshToken)

    err := s.RefRepo.DeleteRefreshToken(tokenHash)

    if err != nil {
        logError(ErrLogout, "Delete Refresh Token")
        return ErrLogout
    }

    return nil
}