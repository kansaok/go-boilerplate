package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/modules/auth"
	"github.com/kansaok/go-boilerplate/internal/util"
)

type FailedLoginStore struct {
	mu      sync.Mutex
	counts  map[string]*loginAttempt
	entries []string
}

type loginAttempt struct {
	count       int
	lastFailed  time.Time
	lockedUntil time.Time
}

const maxEntries = 100000

var failedLogins = &FailedLoginStore{
	counts: make(map[string]*loginAttempt),
}

const (
	maxLoginAttempts = 5
	lockoutDuration  = time.Minute * 15
)

func (s *FailedLoginStore) increment(ip string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.entries) > maxEntries {
		s.evictStale()
	}

	attempt, exists := s.counts[ip]
	if !exists {
		s.counts[ip] = &loginAttempt{count: 1, lastFailed: time.Now()}
		s.entries = append(s.entries, ip)
		return 1
	}

	if time.Now().After(attempt.lockedUntil) && attempt.count >= maxLoginAttempts {
		attempt.count = 1
		attempt.lastFailed = time.Now()
		return 1
	}

	attempt.count++
	attempt.lastFailed = time.Now()

	if attempt.count >= maxLoginAttempts {
		attempt.lockedUntil = time.Now().Add(lockoutDuration)
	}

	return attempt.count
}

func (s *FailedLoginStore) reset(ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.counts, ip)
}

func (s *FailedLoginStore) isLocked(ip string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	attempt, exists := s.counts[ip]
	if !exists {
		return false
	}

	if time.Now().After(attempt.lockedUntil) {
		delete(s.counts, ip)
		return false
	}

	return true
}

func (s *FailedLoginStore) evictStale() {
	cutoff := time.Now().Add(-lockoutDuration * 2)
	var remaining []string
	for _, ip := range s.entries {
		attempt, exists := s.counts[ip]
		if !exists || attempt.lastFailed.Before(cutoff) {
			delete(s.counts, ip)
		} else {
			remaining = append(remaining, ip)
		}
	}
	s.entries = remaining
}

func AuthMiddleware(jwtConfig *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if c.FullPath() == "/metrics" || c.FullPath() == "/health" {
			c.Next()
			return
		}

		if failedLogins.isLocked(ip) {
			util.RespondWithError(c, util.CodeForbidden, "Account temporarily locked due to too many failed attempts", nil)
			c.Abort()
			return
		}

		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			util.RespondWithError(c, util.CodeUnauthorized, util.MESSAGES["TOKEN_NOTFOUND"], nil)
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenStr, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtConfig.SecretKey), nil
		})

		if err != nil || !token.Valid {
			failedLogins.increment(ip)
			util.RespondWithError(c, util.CodeUnauthorized, util.MESSAGES["INVALID_TOKEN"], nil)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*auth.Claims)
		if !ok || claims.Email == "" {
			util.RespondWithError(c, util.CodeUnauthorized, util.MESSAGES["INVALID_TOKEN"], nil)
			c.Abort()
			return
		}

		c.Set("userEmail", claims.Email)
		c.Set("userSubject", claims.Subject)

		failedLogins.reset(ip)
		c.Next()
	}
}

func CSRFValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// JWT-based requests authenticate via the Authorization header, not cookies,
		// so CSRF protection does not apply to them.
		if c.GetHeader("Authorization") != "" {
			c.Next()
			return
		}

		if _, err := c.Cookie("session_id"); err != nil {
			c.Next()
			return
		}

		csrfCookie, err := c.Cookie("csrf_token")
		if err != nil || csrfCookie == "" {
			util.RespondWithError(c, util.CodeForbidden, "CSRF token missing", nil)
			c.Abort()
			return
		}

		csrfHeader := c.GetHeader("X-CSRF-Token")
		if csrfHeader == "" || csrfHeader != csrfCookie {
			util.RespondWithError(c, util.CodeForbidden, "CSRF token mismatch", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
