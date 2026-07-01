package middleware

import (
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/util"
)

type FailedLoginStore struct {
	mu     sync.Mutex
	counts map[string]*loginAttempt
}

type loginAttempt struct {
	count       int
	lastFailed  time.Time
	lockedUntil time.Time
}

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

	attempt, exists := s.counts[ip]
	if !exists {
		s.counts[ip] = &loginAttempt{count: 1, lastFailed: time.Now()}
		return 1
	}

	// Reset if lockout expired
	if time.Now().After(attempt.lockedUntil) {
		attempt.count = 1
		attempt.lastFailed = time.Now()
		return 1
	}

	attempt.count++
	attempt.lastFailed = time.Now()

	// Lock after max attempts
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

func AuthMiddleware(jwtConfig *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Skip rate limiting for health checks and public routes
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

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.SecretKey), nil
		})

		if err != nil || !token.Valid {
			failedLogins.increment(ip)
			util.RespondWithError(c, util.CodeUnauthorized, util.MESSAGES["INVALID_TOKEN"], nil)
			c.Abort()
			return
		}

		// Reset counter on successful auth
		failedLogins.reset(ip)
		c.Next()
	}
}
