package utils

import (
	"net/http"
	"time"
)

const rpsLimit = 10

func RateLimitedHandler(next http.Handler) http.Handler {
	ticker := time.NewTicker(time.Second / time.Duration(rpsLimit))
	sem := make(chan struct{}, rpsLimit)

	go func() {
		for range ticker.C {
			select {
			case sem <- struct{}{}:
			default:
			}
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-sem:
			next.ServeHTTP(w, r)
		default:
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
		}
	})
}
