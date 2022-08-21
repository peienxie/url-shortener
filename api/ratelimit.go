package api

import (
	"net/http"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func errorHandler(c *gin.Context, info ratelimit.Info) {
	errmsg := "too many requests. try again in " + time.Until(info.ResetTime).String()
	c.String(http.StatusTooManyRequests, errmsg)
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

// NewRateLimiter creates a rate limiter gin handler by given rate and limit
// parameters with default error handler function.
func NewRateLimiter(rate time.Duration, limit uint) gin.HandlerFunc {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  rate,
		Limit: limit,
	})

	limiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	return limiter
}
