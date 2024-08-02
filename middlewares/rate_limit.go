package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimitConfig struct {
	RequestLimit int
	WindowSize   time.Duration
}

type RateLimitData struct {
	Count     int64
	Timestamp time.Time
}

var rateLimitStore = make(map[string]RateLimitData)
var mu sync.Mutex

func RateLimit(config RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the IP address
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			// Check and update the rate limit data
			if data, exists := rateLimitStore[ip]; exists {
				if time.Since(data.Timestamp) < config.WindowSize {
					if data.Count >= int64(config.RequestLimit) {
						w.WriteHeader(http.StatusTooManyRequests)
						w.Write([]byte("429 - Too Many Requests"))
						return
					}
					data.Count++
				} else {
					data.Count = 1
					data.Timestamp = time.Now()
				}
				rateLimitStore[ip] = data
			} else {
				rateLimitStore[ip] = RateLimitData{Count: 1, Timestamp: time.Now()}
			}

			next.ServeHTTP(w, r)
		})
	}
}
