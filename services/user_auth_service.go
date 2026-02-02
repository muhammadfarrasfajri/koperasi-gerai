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

func (s *UserAuthService) Register(idToken string, user models.BaseUser) (res *models.BaseUser, err error) {
    ctx := context.Background()

    logError := func(err error, context string) {
        prettyJSON, _ := json.MarshalIndent(user, "", "  ")
       
        shortToken := ""
        if len(idToken) > 10 { shortToken = idToken[:10] + "..." }

        log.Printf("[ERROR] %s | Token: %s | Err: %v\n", 
        context, shortToken, err)
        fmt.Printf("\n[DEBUG DATA USER]:\n%s\n\n", string(prettyJSON))
    }
    
    if len(user.PhoneNumber) < 11 || len(user.PhoneNumber) > 15 {
        logError(ErrInvalidNoHp, "Phone number validation")
        return nil, ErrInvalidNoHp
    }

    if len(user.NIK) != 16 {
        logError(ErrExistingNIK, "NIK validation")
        return nil, ErrInvalidNIK
    }

    if len(user.PosCode) != 5 {
        return nil, ErrInvalidPoscode
    }

    // 2. Verifikasi Token Firebase
    token, err := s.FirebaseAuth.VerifyIDToken(ctx, idToken)
    if err != nil {
        logError(ErrInvalidToken, "Firebase token verification")
        return nil, ErrInvalidToken
    }

    // Ambil data dari token
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
    
    existingUser, err := s.AuthRepo.FindByEmail(user.Email)

    if err == nil && existingUser != nil {

        if existingUser.GoogleUID == "" || existingUser.GoogleUID != user.GoogleUID {
            errLink := s.AuthRepo.LinkGoogleAccount(user.Email, user.GoogleUID, user.GooglePicture)
            if errLink != nil {
                return nil, fmt.Errorf("gagal menghubungkan akun: %v", errLink)
            }

            user.IDMember = existingUser.IDMember
        }
        return &user, nil
    }
    
    // Cek NIK
    nikExists, _ := s.AuthRepo.IsNIKExists(user.NIK)
    if nikExists {
        logError(ErrExistingNIK, "Checking exists NIK")
        return nil, ErrExistingNIK
    }

    // Cek No HP
    hpExists, _ := s.AuthRepo.IsNoHPExists(user.PhoneNumber)
    if hpExists {
        logError(ErrExistingPhoneNo, "Checking exists phone number")
        return nil, ErrExistingPhoneNo
    }

    currentYear := time.Now().Format("06")

    prefix := fmt.Sprintf("KOP-KF-%s", currentYear)
   
    lastID, err := s.AuthRepo.GetMemberId(prefix)
    if err != nil {
        return nil, err 
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

        return nil, err
    }

    return &user, nil
}

func (s *UserAuthService) Login(idToken string, loginHistory models.BaseLoginHistory) (map[string]interface{}, error) {
    ctx := context.Background()

    // 1. Verifikasi Firebase Token
    token, err := s.FirebaseAuth.VerifyIDToken(ctx, idToken)
    if err != nil {
        login, _ := json.Marshal(loginHistory)
        log.Println("Id token = ", idToken)
        log.Println("Data User = ", string(login))
        return nil, ErrInvalidToken
    }

    // 2. Cari User via Google UID
    user, err := s.UserRepo.FindByGoogleUID(token.UID)
    
    // ------------------------------------------------------------------
    // LOGIC ACCOUNT LINKING
    // ------------------------------------------------------------------
    if err != nil { 

        email, _ := token.Claims["email"].(string)
        existingUser, errEmail := s.AuthRepo.FindByEmail(email)

        if errEmail == nil && existingUser != nil {

            pic, _ := token.Claims["picture"].(string)
            _ = s.AuthRepo.LinkGoogleAccount(email, token.UID, pic)
            
            user, err = s.UserRepo.FindByGoogleUID(token.UID)
            if err != nil {
                return nil, errors.New("gagal mengambil data user setelah linking")
            }
        } else {
            return nil, errors.New("User not found. Please Register first.")
        }
    }

    if user == nil {
        return nil, errors.New("user data is invalid")
    }
    finalStatus := "FAILED" 

    defer func() {
        loginHistory.UserID = user.ID
        loginHistory.Status = finalStatus

        if errHist := s.AuthRepo.HistoryLoginUser(loginHistory); errHist != nil {
            log.Println("Gagal simpan history:", errHist)
        }
    }()

    // 3. Cek Status Verifikasi
    if user.Is_verified == 0 {
        // Otomatis defer jalan -> Status FAILED
        return nil, errors.New("User is not verified")
    }

    // [DIHAPUS] HistoryLoginUser manual disini sudah dihapus agar tidak double

    // 4. Generate Access Token
    accessToken, err := s.JWTSecret.GenerateAccessToken(user.ID, user.Email)
    if err != nil {
        return nil, errors.New("error JWT access")
    }

    // 5. Generate Refresh Token
    refreshToken, err := s.JWTSecret.GenerateRefreshToken(user.ID)
    if err != nil {
        return nil, errors.New("error JWT refresh")
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
        return nil, errors.New("gagal menyimpan session login")
    }

    // 8. Sukses!
    finalStatus = "SUCCESS"

    return map[string]interface{}{
        "message":       "login success",
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "user": map[string]interface{}{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
            "role":  user.ActiveAs,
        },
    }, nil
}

func (s *UserAuthService) RefreshToken(rawRefreshToken string) (map[string]interface{}, error) {
    
    // 1. Parse Token dengan Safety Check
    token, err := jwt.Parse(rawRefreshToken, func(t *jwt.Token) (interface{}, error) {
        // Best Practice: Cek Signing Method
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return s.JWTSecret.RefreshSecret, nil
    })

    // PERBAIKAN 3: Cek Error DULU sebelum lanjut
    if err != nil {
        return nil, errors.New("refresh token invalid or expired")
    }
    
    // Ambil Claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token claims")
    }

    // Ambil User ID
    userIDFloat, ok := claims["user_id"].(float64)
    if !ok {
        return nil, errors.New("invalid user id in token")
    }
    userID := int(userIDFloat)

    // 2. Cek Token di Database
    tokenCheck, err := s.RefRepo.FindRefreshTokenUser(userID)
    if err != nil || tokenCheck == nil {
        return nil, errors.New("refresh token not found in db")
    }

    // PERBAIKAN 1: Hash dulu token dari user, baru bandingkan
    incomingTokenHash := middleware.HashToken(rawRefreshToken)
    
    if incomingTokenHash != tokenCheck.TokenHash {
        // Ini indikasi Token Reuse Attack (Bahaya!)
        // Opsional: Kamu bisa hapus token di DB biar user dipaksa login ulang demi keamanan
        return nil, errors.New("refresh token reuse detected / not match")
    }

    // 3. Ambil User untuk Data Token Baru
    user, err := s.UserRepo.FindById(strconv.Itoa(userID))
    if err != nil || user == nil {
        return nil, errors.New("user not found")
    }

    // 4. Generate Token Baru
    newAccessToken, err := s.JWTSecret.GenerateAccessToken(user.ID, user.Email)
    if err != nil {
        return nil, err
    }

    newRefreshToken, err := s.JWTSecret.GenerateRefreshToken(user.ID)
    if err != nil {
        return nil, err
    }

    // 5. Simpan Token Baru ke Database
    newRefreshTokenHash := middleware.HashToken(newRefreshToken)
    expiresAt := time.Now().Add(7 * 24 * time.Hour)

    tokenModel := models.RefreshToken{
        UserID:    user.ID,        // PERBAIKAN 2: Jangan lupa isi UserID!
        TokenHash: newRefreshTokenHash,
        ExpiresAt: expiresAt,
    }

    err = s.RefRepo.UpsertRefreshToken(tokenModel)
    if err != nil {
        return nil, errors.New("failed to save new session")
    }

    return map[string]interface{}{
        "access_token":  newAccessToken,
        "refresh_token": newRefreshToken,
    }, nil
}

func (s *UserAuthService) Logout(rawRefreshToken string) error {
    
    // 1. LAKUKAN HASHING ULANG DISINI
    // Gunakan fungsi yang SAMA PERSIS dengan yang kamu pakai di fungsi Login
    tokenHash := middleware.HashToken(rawRefreshToken)

    // 2. Sekarang 'tokenHash' isinya sudah cocok dengan yang di Database
    // Panggil Repo untuk hapus berdasarkan hash tersebut
    err := s.RefRepo.DeleteRefreshToken(tokenHash)
    
    if err != nil {
        return errors.New("gagal logout / token tidak ditemukan")
    }
    
    return nil
}