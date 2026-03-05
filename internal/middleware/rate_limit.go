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
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	var mu sync.Mutex
	limiters := make(map[string]*ipLimiter)

	// Periodic cleanup of stale entries every 3 minutes.
	go func() {
		for {
			time.Sleep(3 * time.Minute)
			mu.Lock()
			for ip, l := range limiters {
				if time.Since(l.lastSeen) > 5*time.Minute {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
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
