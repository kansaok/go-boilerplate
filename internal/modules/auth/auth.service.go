package auth

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kansaok/go-boilerplate/internal/config"
	usr "github.com/kansaok/go-boilerplate/internal/modules/user"
	"github.com/kansaok/go-boilerplate/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser mendaftarkan user baru dengan email dan password
func RegisterUser(ctx context.Context, req RegisterRequest) (map[string]interface{}, error) {
    if req.Password != req.ConfirmPassword {
        return nil, errors.New("password dan konfirmasi password tidak cocok")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    bod, _ := util.ParseDate(req.Bod)

    user := usr.User{
        Email:             req.Email,
        Password:          string(hashedPassword),
        Title:             req.Title,
        FirstName:         req.FirstName,
        MiddleName:        &req.MiddleName,
        LastName:          &req.LastName,
        Gender:            req.Gender,
        Bod:               bod,
        Pob:               req.Pob,
        PhoneNumber:       req.PhoneNumber,
    }

    createdUser, err := CreateUser(ctx, user)
    if err != nil {
        return nil, err
    }

    return createdUser, nil
}

func AuthenticateUser(email, password string, jwtConfig *config.JWTConfig) (string, error) {
    user, err := GetUserByEmail(email)
    if err != nil {
        return "", err
    }

    if user == nil {
        return "", errors.New("email tidak ditemukan")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("password salah")
    }

    token, err := GenerateToken(user.Email, jwtConfig)
    if err != nil {
        return "", err
    }

    return token, nil
}

func GenerateToken(email string, jwtConfig *config.JWTConfig) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(jwtConfig.AccessTokenLifetime).Unix(),
    })
    return token.SignedString([]byte(jwtConfig.SecretKey))
}
