package rate

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	limit    int
	interval time.Duration
	mtx      sync.Mutex
	times    []time.Time
}

// New creates a new rate limiter for the limit and interval.
func New(limit int, interval time.Duration) *RateLimiter {
	lim := &RateLimiter{
		limit:    limit,
		interval: interval,
	}
	lim.times = append(lim.times, time.Now())
	return lim
}

// Try returns true if under the rate limit, or false if over and the
// remaining time before the rate limit expires.
func (r *RateLimiter) Try() (ok bool, remaining time.Duration) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	now := time.Now()
	if l := len(r.times); l < r.limit {
		r.times = append(r.times, time.Now())
		return true, 0
	}
	frnt := r.times[len(r.times)-1]
	if diff := now.Sub(frnt); diff < r.interval {
		return false, r.interval - diff
	}
	fmt.Println(now.Sub(frnt))
	frnt = now
	r.times = append(r.times, frnt)

	return true, 0
}
