package middleware

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/util"
)

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*client
	interval time.Duration
	maxReqs  int
}

type client struct {
	tokens    int
	lastSeen  time.Time
}

var globalLimiter = &RateLimiter{
	clients:  make(map[string]*client),
	interval: time.Minute,
	maxReqs:  100,
}

func NewRateLimiter(interval time.Duration, maxReqs int) *RateLimiter {
	return &RateLimiter{
		clients:  make(map[string]*client),
		interval: interval,
		maxReqs:  maxReqs,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[key]
	if !exists {
		rl.clients[key] = &client{
			tokens:   rl.maxReqs,
			lastSeen: time.Now(),
		}
		return true
	}

	// Reset tokens if interval has passed
	if time.Since(c.lastSeen) > rl.interval {
		c.tokens = rl.maxReqs
		c.lastSeen = time.Now()
	}

	if c.tokens > 0 {
		c.tokens--
		return true
	}

	return false
}

func (rl *RateLimiter) GetRemaining(ctx context.Context, key string) int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[key]
	if !exists {
		return rl.maxReqs
	}

	if time.Since(c.lastSeen) > rl.interval {
		return rl.maxReqs
	}

	return c.tokens
}

// RateLimitMiddleware returns a middleware that rate limits requests per IP
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.Allow(c.Request.Context(), ip) {
			util.RespondWithError(c, util.CodeBadRequest, "Too many requests. Please try again later.", nil)
			c.Abort()
			return
		}

		c.Set("rate_limit_remaining", rl.GetRemaining(c.Request.Context(), ip))
		c.Set("rate_limit_reset", time.Now().Add(rl.interval).Unix())
		c.Next()
	}
}

// RateLimitAuth returns a more restrictive middleware for auth endpoints
func RateLimitAuth() gin.HandlerFunc {
	authLimiter := &RateLimiter{
		clients:  make(map[string]*client),
		interval: time.Minute * 15,
		maxReqs:  10,
	}
	return RateLimitMiddleware(authLimiter)
}

// GlobalLimiter is the default global rate limiter
var GlobalLimiter = globalLimiter
