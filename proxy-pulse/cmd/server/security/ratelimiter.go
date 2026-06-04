package security

import (
	"sync"
	"time"
)

type tokenBucket struct {
	limit      float64
	refill     time.Duration
	tokens     float64
	lastRefill time.Time
}

func newTokenBucket(limit float64, window time.Duration) *tokenBucket {
	return &tokenBucket{limit: limit, refill: window, tokens: limit, lastRefill: time.Now()}
}

func (tb *tokenBucket) Allow() bool {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += elapsed * tb.limit / tb.refill.Seconds()
	if tb.tokens > tb.limit {
		tb.tokens = tb.limit
	}
	tb.lastRefill = now
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}
	return false
}

type RateLimiter struct {
	buckets sync.Map
	limit   float64
	window  time.Duration
}

func NewRateLimiter(qps float64) *RateLimiter {
	return &RateLimiter{limit: qps, window: time.Second}
}

func (rl *RateLimiter) Allow(ip string) bool {
	b, _ := rl.buckets.LoadOrStore(ip, newTokenBucket(rl.limit, rl.window))
	return b.(*tokenBucket).Allow()
}
