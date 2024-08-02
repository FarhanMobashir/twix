package middlewares

import (
	"net/http"
)

func ContentType(allowedTypes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			requestContentType := r.Header.Get("Content-Type")

			for _, contentType := range allowedTypes {
				if requestContentType == contentType {
					w.Header().Set("Content-Type", contentType)
					break
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
