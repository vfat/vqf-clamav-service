package ratelimit

import (
	"testing"
	"time"
)

func TestTokenBucketLimiter_Allow(t *testing.T) {
	limiter := NewLimiter()

	key := "api_key_test_123"
	limitRPM := 5

	// First 5 requests must be allowed
	for i := 0; i < limitRPM; i++ {
		allowed, remaining, _ := limiter.Allow(key, limitRPM)
		if !allowed {
			t.Fatalf("request %d should have been allowed", i+1)
		}
		expectedRemaining := limitRPM - (i + 1)
		if remaining != expectedRemaining {
			t.Errorf("expected remaining %d, got %d", expectedRemaining, remaining)
		}
	}

	// 6th request must be rejected
	allowed, remaining, resetTime := limiter.Allow(key, limitRPM)
	if allowed {
		t.Fatalf("6th request should be blocked (rate limit exceeded)")
	}
	if remaining != 0 {
		t.Errorf("expected remaining 0, got %d", remaining)
	}
	if resetTime <= time.Now().Unix() {
		t.Errorf("expected resetTime in future, got %d", resetTime)
	}
}

func TestTokenBucketLimiter_GlobalRPS(t *testing.T) {
	limiter := NewLimiter()
	globalRPS := 3

	for i := 0; i < globalRPS; i++ {
		if !limiter.AllowGlobal(globalRPS) {
			t.Fatalf("global request %d should be allowed", i+1)
		}
	}

	// Exceeded
	if limiter.AllowGlobal(globalRPS) {
		t.Fatalf("global request %d should be rejected", globalRPS+1)
	}
}
