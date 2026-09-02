package ratelimit

import (
	"sync"
	"time"
)

type bucket struct {
	tokens     int
	lastRefill time.Time
}

// Limiter manages in-memory token bucket rate limiting.
type Limiter struct {
	mu           sync.Mutex
	buckets      map[string]*bucket
	globalBucket *bucket
}

// NewLimiter initializes a rate limiter.
func NewLimiter() *Limiter {
	return &Limiter{
		buckets:      make(map[string]*bucket),
		globalBucket: nil,
	}
}

// Allow checks if a request with given key is allowed under limitRPM (requests per minute).
// Returns (allowed, remainingTokens, resetEpochSeconds).
func (l *Limiter) Allow(key string, limitRPM int) (bool, int, int64) {
	if limitRPM <= 0 {
		return true, 100, time.Now().Add(time.Minute).Unix()
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, exists := l.buckets[key]
	if !exists {
		b = &bucket{
			tokens:     limitRPM,
			lastRefill: now,
		}
		l.buckets[key] = b
	}

	// Refill if 1 minute window has elapsed
	if now.Sub(b.lastRefill) >= time.Minute {
		b.tokens = limitRPM
		b.lastRefill = now
	}

	resetTime := b.lastRefill.Add(time.Minute).Unix()

	if b.tokens > 0 {
		b.tokens--
		return true, b.tokens, resetTime
	}

	return false, 0, resetTime
}

// AllowGlobal checks if the request is allowed under global RPS (requests per second).
func (l *Limiter) AllowGlobal(limitRPS int) bool {
	if limitRPS <= 0 {
		return true
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if l.globalBucket == nil || l.globalBucket.lastRefill.IsZero() || now.Sub(l.globalBucket.lastRefill) >= time.Second {
		l.globalBucket = &bucket{
			tokens:     limitRPS,
			lastRefill: now,
		}
	}

	if l.globalBucket.tokens > 0 {
		l.globalBucket.tokens--
		return true
	}

	return false
}
