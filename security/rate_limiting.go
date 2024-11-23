package security

import (
	"sync"
	"time"
)

type RateLimiter struct {
	Component struct{}
	limits    map[string]*userRateLimit
	mu        sync.Mutex
}

type userRateLimit struct {
	tokens     int
	lastAccess time.Time
}

const maxTokens = 50
const refillInterval = time.Minute

func (rl *RateLimiter) PostConstruct() {
	rl.limits = make(map[string]*userRateLimit)
}

func (rl *RateLimiter) AllowRequest(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	ul, exists := rl.limits[ip]

	if !exists {
		ul = &userRateLimit{tokens: maxTokens, lastAccess: now}
		rl.limits[ip] = ul
	}

	// Refill tokens based on time passed
	if now.Sub(ul.lastAccess) > refillInterval {
		ul.tokens = maxTokens
		ul.lastAccess = now
	}

	if ul.tokens > 0 {
		ul.tokens--
		return true
	}

	return false
}
