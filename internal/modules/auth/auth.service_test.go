package auth

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kansaok/go-boilerplate/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func testJWTConfig() *config.JWTConfig {
	return &config.JWTConfig{
		SecretKey:            "test-secret-key-that-is-long-enough-32-plus",
		AccessTokenLifetime:  time.Minute * 15,
		RefreshTokenLifetime: time.Hour * 24,
	}
}

func TestGenerateToken_ProducesValidHS256Token(t *testing.T) {
	cfg := testJWTConfig()
	email := "user@example.com"

	tokenStr, err := GenerateToken(email, cfg)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.SecretKey), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("token should be valid, got err=%v valid=%v", err, token.Valid)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Fatal("expected Claims type")
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
	if claims.Subject != email {
		t.Errorf("expected subject %s, got %s", email, claims.Subject)
	}
}

func TestGenerateToken_RejectsWrongSecret(t *testing.T) {
	cfg := testJWTConfig()
	tokenStr, err := GenerateToken("a@b.c", cfg)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	_, err = jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("different-secret-key-not-the-right-one!!!"), nil
	})
	if err == nil {
		t.Fatal("expected signature validation error with wrong secret")
	}
}

func TestGenerateToken_RejectsNonHS256Algorithm(t *testing.T) {
	cfg := testJWTConfig()

	// Craft a raw JWT with header "alg":"RS256" (algorithm confusion / key-substitution attack).
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"email":"attacker@example.com","exp":9999999999}`))
	sig := base64.RawURLEncoding.EncodeToString([]byte("fake-rsa-signature"))
	rsaToken := header + "." + payload + "." + sig

	_, err := jwt.ParseWithClaims(rsaToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.SecretKey), nil
	})
	if err == nil {
		t.Fatal("algorithm confusion token must be rejected")
	}
}

func TestRegisterUser_PasswordMismatch(t *testing.T) {
	req := RegisterRequest{
		Password:        "Ab1!defgh",
		ConfirmPassword: "different",
	}

	_, err := RegisterUser(context.Background(), req)
	if err == nil {
		t.Fatal("expected password mismatch error")
	}
	if !strings.Contains(err.Error(), "tidak cocok") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestInit_DummyBcryptHashPrepared(t *testing.T) {
	if len(dummyBcryptHash) == 0 {
		t.Fatal("dummy bcrypt hash should be prepared at init")
	}
	// Timing-equalization compare must succeed against the dummy value
	if err := bcrypt.CompareHashAndPassword(dummyBcryptHash, []byte("timing-equalization-dummy")); err != nil {
		t.Fatalf("dummy hash should match its source password: %v", err)
	}
}