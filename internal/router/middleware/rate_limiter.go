package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	RedisClient *redis.Client
	Limit       int           // max requests
	Window      time.Duration // e.g., 1 minute
}

func NewRateLimiter(redisClient *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		RedisClient: redisClient,
		Limit:       limit,
		Window:      window,
	}
}

func (rl *RateLimiter) LimitByIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		ip := r.RemoteAddr

		key := fmt.Sprintf("rate:%s", ip)

		count, err := rl.RedisClient.Incr(ctx, key).Result()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if count == 1 {
			rl.RedisClient.Expire(ctx, key, rl.Window)
		}

		if count > int64(rl.Limit) {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(rl.Window.Seconds())))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
