package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/redis/go-redis/v9"
)

const (
	loginMaxAttempts    = 5
	loginWindow         = 15 * time.Minute
	registerMaxAttempts = 5
	registerWindow      = 1 * time.Hour
	memMaxEntries       = 100000
	redisOpTimeout      = 1 * time.Second
)

// AccountLimiter membatasi percobaan login/registrasi per-akun (email).
// Saat Redis dikonfigurasi, counter dibagikan antar-instance (wajib untuk
// deployment multi-instance). Tanpa Redis, digunakan penyimpanan in-memory
// yang masih efektif untuk instalasi single-instance.
type AccountLimiter struct {
	client  *redis.Client
	enabled bool

	mu  sync.Mutex
	mem map[string]*memAttempt
}

type memAttempt struct {
	count   int
	resetAt time.Time
}

var (
	limiterOnce sync.Once
	limiter     *AccountLimiter
)

// GetAccountLimiter mengembalikan singleton limiter (diinisialisasi sekali).
func GetAccountLimiter() *AccountLimiter {
	limiterOnce.Do(func() {
		limiter = newAccountLimiter(config.LoadConfig().RedisConfig)
	})
	return limiter
}

// SetAccountLimiter menggantikan singleton (digunakan oleh unit test).
func SetAccountLimiter(l *AccountLimiter) {
	limiterOnce.Do(func() {})
	limiter = l
}

func newAccountLimiter(cfg *config.RedisConfig) *AccountLimiter {
	l := &AccountLimiter{mem: make(map[string]*memAttempt, 1024)}
	if cfg == nil || !cfg.Enabled {
		log.Println("Account limiter: REDIS_HOST belum di-set, pakai in-memory fallback")
		return l
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("WARNING: Redis tidak tersedia (%v); fallback ke in-memory limiter", err)
		_ = client.Close()
		return l
	}

	l.client = client
	l.enabled = true
	log.Println("Account limiter menggunakan Redis")
	return l
}

func (l *AccountLimiter) IsRedisEnabled() bool {
	return l.enabled
}

// CheckLoginLocked true jika akun sedang dikunci karena terlalu banyak gagal.
func (l *AccountLimiter) CheckLoginLocked(account string) bool {
	key := accountKey("login:attempts", account)
	if l.enabled {
		ctx, cancel := context.WithTimeout(context.Background(), redisOpTimeout)
		defer cancel()
		count, err := l.client.Get(ctx, key).Int()
		if err != nil {
			return false
		}
		return count >= loginMaxAttempts
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	return l.checkMemLocked(key, loginMaxAttempts, time.Now())
}

// RecordLoginFailure mencatat satu kegagalan login untuk sebuah akun.
func (l *AccountLimiter) RecordLoginFailure(account string) {
	key := accountKey("login:attempts", account)
	if l.enabled {
		ctx, cancel := context.WithTimeout(context.Background(), redisOpTimeout)
		defer cancel()
		count, err := l.client.Incr(ctx, key).Result()
		if err != nil {
			return
		}
		if count == 1 {
			l.client.Expire(ctx, key, loginWindow)
		}
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.recordMem(key, loginWindow, time.Now())
}

// ResetLoginFailures menghapus counter kegagalan (dipanggil saat login sukses).
func (l *AccountLimiter) ResetLoginFailures(account string) {
	key := accountKey("login:attempts", account)
	if l.enabled {
		ctx, cancel := context.WithTimeout(context.Background(), redisOpTimeout)
		defer cancel()
		l.client.Del(ctx, key)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.mem, key)
}

// CheckRegisterThrottled true jika email sudah melewati batas percobaan
// registrasi dalam satu window (anti CPU-flood bcrypt via /register).
func (l *AccountLimiter) CheckRegisterThrottled(email string) bool {
	key := accountKey("register:attempts", email)
	if l.enabled {
		ctx, cancel := context.WithTimeout(context.Background(), redisOpTimeout)
		defer cancel()
		count, err := l.client.Get(ctx, key).Int()
		if err != nil {
			return false
		}
		return count >= registerMaxAttempts
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	return l.checkMemLocked(key, registerMaxAttempts, time.Now())
}

// RecordRegisterAttempt mencatat satu percobaan registrasi.
func (l *AccountLimiter) RecordRegisterAttempt(email string) {
	key := accountKey("register:attempts", email)
	if l.enabled {
		ctx, cancel := context.WithTimeout(context.Background(), redisOpTimeout)
		defer cancel()
		count, err := l.client.Incr(ctx, key).Result()
		if err != nil {
			return
		}
		if count == 1 {
			l.client.Expire(ctx, key, registerWindow)
		}
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.recordMem(key, registerWindow, time.Now())
}

func (l *AccountLimiter) checkMemLocked(key string, maxAttempts int, now time.Time) bool {
	attempt, ok := l.mem[key]
	if !ok {
		return false
	}
	if now.After(attempt.resetAt) {
		delete(l.mem, key)
		return false
	}
	return attempt.count >= maxAttempts
}

func (l *AccountLimiter) recordMem(key string, window time.Duration, now time.Time) {
	if len(l.mem) >= memMaxEntries {
		l.purgeExpiredLocked(now)
	}
	attempt, ok := l.mem[key]
	if !ok {
		l.mem[key] = &memAttempt{count: 1, resetAt: now.Add(window)}
		return
	}
	attempt.count++
}

func (l *AccountLimiter) purgeExpiredLocked(now time.Time) {
	for key, attempt := range l.mem {
		if now.After(attempt.resetAt) {
			delete(l.mem, key)
		}
	}
}

func accountKey(kind, account string) string {
	normalized := strings.ToLower(strings.TrimSpace(account))
	sum := sha256.Sum256([]byte(normalized))
	return fmt.Sprintf("auth:%s:%x", kind, sum[:12])
}