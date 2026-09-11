package middleware

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter_AllowsUntilLimit(t *testing.T) {
	rl := NewRateLimiter(time.Hour, 3)
	ctx := context.Background()

	allowed := []bool{
		rl.Allow(ctx, "1.2.3.4"),
		rl.Allow(ctx, "1.2.3.4"),
		rl.Allow(ctx, "1.2.3.4"),
	}
	for i, a := range allowed {
		if !a {
			t.Fatalf("request %d should be allowed", i)
		}
	}

	if rl.Allow(ctx, "1.2.3.4") {
		t.Fatal("4th request should be blocked (limit 3)")
	}
}

func TestRateLimiter_IndependentClients(t *testing.T) {
	rl := NewRateLimiter(time.Hour, 2)
	ctx := context.Background()

	// client-a exhausts its full quota: 2 allowed, 3rd blocked
	if !rl.Allow(ctx, "client-a") {
		t.Fatal("client-a request 1 should pass")
	}
	if !rl.Allow(ctx, "client-a") {
		t.Fatal("client-a request 2 should pass")
	}
	if rl.Allow(ctx, "client-a") {
		t.Fatal("client-a request 3 should be blocked")
	}

	// client-b must be unaffected by client-a's exhaustion
	if !rl.Allow(ctx, "client-b") {
		t.Fatal("client-b request 1 should pass independently")
	}
}

func TestRateLimiter_ResetsAfterInterval(t *testing.T) {
	rl := NewRateLimiter(50*time.Millisecond, 1)
	ctx := context.Background()

	if !rl.Allow(ctx, "x") {
		t.Fatal("first request should pass")
	}
	if rl.Allow(ctx, "x") {
		t.Fatal("second request should be blocked within interval")
	}

	time.Sleep(80 * time.Millisecond)
	if !rl.Allow(ctx, "x") {
		t.Fatal("request after interval should be allowed")
	}
}

func TestRateLimiter_CleanupEvictsStaleEntries(t *testing.T) {
	rl := NewRateLimiter(50*time.Millisecond, 1)
	ctx := context.Background()

	rl.Allow(ctx, "stale-client")
	if _, exists := rl.clients["stale-client"]; !exists {
		t.Fatal("client should be tracked initially")
	}

	time.Sleep(150 * time.Millisecond)
	rl.cleanup()

	if _, exists := rl.clients["stale-client"]; exists {
		t.Fatal("stale client entry should be evicted after 2x interval")
	}
}

func TestRateLimiter_GetRemaining(t *testing.T) {
	rl := NewRateLimiter(time.Hour, 5)
	ctx := context.Background()

	if got := rl.GetRemaining(ctx, "new-client"); got != 5 {
		t.Fatalf("expected 5 remaining for unseen client, got %d", got)
	}

	rl.Allow(ctx, "new-client")
	if got := rl.GetRemaining(ctx, "new-client"); got != 4 {
		t.Fatalf("expected 4 remaining after one request, got %d", got)
	}
}