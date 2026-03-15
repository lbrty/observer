package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimit returns a Gin middleware that enforces per-IP rate limiting.
// rps is the sustained requests per second, burst is the maximum burst size.
// Stale IP entries are cleaned up lazily (at most once every 3 minutes) so
// no background goroutine is required.
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	var mu sync.Mutex
	limiters := make(map[string]*ipLimiter)
	lastCleanup := time.Now()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if time.Since(lastCleanup) > 3*time.Minute {
			for k, l := range limiters {
				if time.Since(l.lastSeen) > 5*time.Minute {
					delete(limiters, k)
				}
			}
			lastCleanup = time.Now()
		}
		l, exists := limiters[ip]
		if !exists {
			l = &ipLimiter{
				limiter: rate.NewLimiter(rate.Limit(rps), burst),
			}
			limiters[ip] = l
		}
		l.lastSeen = time.Now()
		mu.Unlock()

		if !l.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
				"code":  "errors.rateLimit.exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
