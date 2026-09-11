package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kansaok/go-boilerplate/internal/config"
	usr "github.com/kansaok/go-boilerplate/internal/modules/user"
	"github.com/kansaok/go-boilerplate/internal/service"
	"github.com/kansaok/go-boilerplate/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var dummyBcryptHash []byte

func init() {
	// Pre-computed dummy hash to equalize login timing for unknown emails
	dummyBcryptHash, _ = bcrypt.GenerateFromPassword([]byte("timing-equalization-dummy"), 12)
}

// RegisterUser mendaftarkan user baru dengan email dan password
func RegisterUser(ctx context.Context, req RegisterRequest) (map[string]interface{}, error) {
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("password dan konfirmasi password tidak cocok")
	}

	limiter := service.GetAccountLimiter()
	if limiter.CheckRegisterThrottled(req.Email) {
		return nil, errors.New("terlalu banyak percobaan registrasi untuk email ini")
	}
	limiter.RecordRegisterAttempt(req.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	bod, _ := util.ParseDate(req.Bod)

	user := usr.User{
		Email:       req.Email,
		Password:    string(hashedPassword),
		Title:       req.Title,
		FirstName:   req.FirstName,
		MiddleName:  &req.MiddleName,
		LastName:    &req.LastName,
		Gender:      req.Gender,
		Bod:         bod,
		Pob:         req.Pob,
		PhoneNumber: req.PhoneNumber,
		CreatedBy:   "self",
	}

	createdUser, err := CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func AuthenticateUser(email, password string, jwtConfig *config.JWTConfig) (string, error) {
	limiter := service.GetAccountLimiter()
	if limiter.CheckLoginLocked(email) {
		return "", errors.New("akun dikunci sementara karena terlalu banyak percobaan gagal")
	}

	user, err := GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if user == nil {
		// Equalize timing with the bcrypt compare to prevent user enumeration
		bcrypt.CompareHashAndPassword(dummyBcryptHash, []byte(password))
		limiter.RecordLoginFailure(email)
		return "", errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		limiter.RecordLoginFailure(email)
		return "", errors.New("email atau password salah")
	}

	limiter.ResetLoginFailures(email)

	token, err := GenerateToken(user.Email, jwtConfig)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateToken(email string, jwtConfig *config.JWTConfig) (string, error) {
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.AccessTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.SecretKey))
}
