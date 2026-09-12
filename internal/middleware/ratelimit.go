package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	rateLimit = 10
	burstSize = 20
)

func RateLimiter() gin.HandlerFunc {
	var (
		mu      sync.Mutex
		clients = make(map[string]*rate.Limiter)
	)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()

		limiter, exists := clients[ip]
		if !exists {
			limiter = rate.NewLimiter(rateLimit, burstSize)
			clients[ip] = limiter
		}

		mu.Unlock()

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
