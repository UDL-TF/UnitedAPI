package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/UDL-TF/UnitedAPI/internal/response"
	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.RWMutex
)

// RateLimit middleware limits requests per IP
// This is a simple in-memory implementation
// For production, consider using Redis-based rate limiting
func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	// Cleanup old visitors every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{
				lastSeen: time.Now(),
				count:    1,
			}
			mu.Unlock()
			c.Next()
			return
		}

		if time.Since(v.lastSeen) > time.Minute {
			v.lastSeen = time.Now()
			v.count = 1
			mu.Unlock()
			c.Next()
			return
		}

		if v.count >= requestsPerMinute {
			mu.Unlock()
			response.Error(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		v.count++
		mu.Unlock()
		c.Next()
	}
}
