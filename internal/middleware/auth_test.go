package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/config"
)

func setupTestRoute(jwtConfig *config.JWTConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware(jwtConfig))
	router.GET("/protected", func(c *gin.Context) {
		email, _ := c.Get("userEmail")
		c.JSON(http.StatusOK, gin.H{"email": email})
	})
	return router
}

func getSecretKey() string {
	return "test-secret-key-that-is-long-enough-32-plus"
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	router := setupTestRoute(&config.JWTConfig{SecretKey: getSecretKey(), AccessTokenLifetime: time.Minute})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := setupTestRoute(&config.JWTConfig{SecretKey: getSecretKey(), AccessTokenLifetime: time.Minute})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_NoneAlgorithmRejected(t *testing.T) {
	router := setupTestRoute(&config.JWTConfig{SecretKey: getSecretKey(), AccessTokenLifetime: time.Minute})

	claims := jwt.MapClaims{
		"email": "attacker@example.com",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("failed crafting none-alg token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("alg:none token must be rejected, got %d", w.Code)
	}
}

func TestAuthMiddleware_HS384DowngradeRejected(t *testing.T) {
	router := setupTestRoute(&config.JWTConfig{SecretKey: getSecretKey(), AccessTokenLifetime: time.Minute})

	claims := jwt.MapClaims{
		"email": "attacker@example.com",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenStr, err := token.SignedString([]byte(getSecretKey()))
	if err != nil {
		t.Fatalf("failed crafting HS384 token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("HS384 downgrade must be rejected, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidTokenPassesAndSetsIdentity(t *testing.T) {
	jwtCfg := &config.JWTConfig{SecretKey: getSecretKey(), AccessTokenLifetime: time.Minute}
	router := setupTestRoute(jwtCfg)

	claims := struct {
		Email string `json:"email"`
		jwt.RegisteredClaims
	}{
		Email: "user@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(getSecretKey()))
	if err != nil {
		t.Fatalf("failed crafting valid token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		body := w.Body.String()
		t.Fatalf("expected 200, got %d body=%s", w.Code, body)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp["email"] != "user@example.com" {
		t.Errorf("expected user email in context, got %v", resp["email"])
	}
}