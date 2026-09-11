package service

import (
	"strings"
	"testing"
	"time"

	"github.com/kansaok/go-boilerplate/internal/config"
)

func newMemLimiter() *AccountLimiter {
	return newAccountLimiter(&config.RedisConfig{Enabled: false})
}

func TestLoginLockoutAfterMaxAttempts(t *testing.T) {
	l := newMemLimiter()
	email := "alice@example.com"

	for i := 0; i < loginMaxAttempts; i++ {
		if i < loginMaxAttempts-1 && l.CheckLoginLocked(email) {
			t.Fatalf("akun terkunci sebelum max attempts (percobaan ke-%d)", i+1)
		}
		l.RecordLoginFailure(email)
	}

	if !l.CheckLoginLocked(email) {
		t.Fatal("akun harus terkunci setelah max attempts tercapai")
	}

	l.RecordLoginFailure(email)
	if !l.CheckLoginLocked(email) {
		t.Fatal("akun tetap terkunci selama window")
	}
}

func TestLoginLockedResetOnSuccess(t *testing.T) {
	l := newMemLimiter()
	email := "bob@example.com"
	for i := 0; i < loginMaxAttempts; i++ {
		l.RecordLoginFailure(email)
	}
	if !l.CheckLoginLocked(email) {
		t.Fatal("akun harus terkunci sebelum reset")
	}
	l.ResetLoginFailures(email)
	if l.CheckLoginLocked(email) {
		t.Fatal("akun tidak boleh terkunci setelah reset")
	}
}

func TestLoginLockExpiresAfterWindow(t *testing.T) {
	l := newMemLimiter()
	email := "carol@example.com"
	for i := 0; i < loginMaxAttempts; i++ {
		l.RecordLoginFailure(email)
	}

	key := accountKey("login:attempts", email)
	l.mu.Lock()
	l.mem[key].resetAt = time.Now().Add(-time.Second)
	l.mu.Unlock()

	if l.CheckLoginLocked(email) {
		t.Fatal("kunci harus kedaluwarsa setelah window berlalu")
	}
}

func TestRegisterThrottle(t *testing.T) {
	l := newMemLimiter()
	email := "dave@example.com"

	for i := 0; i < registerMaxAttempts; i++ {
		l.RecordRegisterAttempt(email)
	}
	if !l.CheckRegisterThrottled(email) {
		t.Fatal("registrasi harus ter-throttle setelah melewati batas")
	}
}

func TestDifferentAccountsAreIsolated(t *testing.T) {
	l := newMemLimiter()
	l.RecordLoginFailure("one@example.com")
	for i := 0; i < loginMaxAttempts; i++ {
		l.RecordLoginFailure("two@example.com")
	}
	if l.CheckLoginLocked("one@example.com") {
		t.Fatal("akun pertama tidak boleh ikut terkunci")
	}
	if !l.CheckLoginLocked("two@example.com") {
		t.Fatal("akun kedua harus terkunci")
	}
}

func TestRedisDisabledFallback(t *testing.T) {
	l := newMemLimiter()
	if l.IsRedisEnabled() {
		t.Fatal("limiter tanpa Redis harus dalam mode in-memory")
	}
	l.RecordLoginFailure("x@example.com")
	_ = l.CheckLoginLocked("x@example.com")
	l.ResetLoginFailures("x@example.com")
	_ = l.CheckRegisterThrottled("x@example.com")
	l.RecordRegisterAttempt("x@example.com")
}

func TestAccountKeyNormalizesCaseAndWhitespace(t *testing.T) {
	normalized := accountKey("login:attempts", "  Alice@Example.COM ")
	raw := accountKey("login:attempts", "alice@example.com")
	if normalized != raw {
		t.Fatalf("key harus dinormalisasi: %s != %s", normalized, raw)
	}
	if len(normalized) > 64 {
		t.Fatal("key terlalu panjang")
	}
	if strings.Contains(normalized, "alice") {
		t.Fatal("key tidak boleh berisi email mentah (privasi)")
	}
}