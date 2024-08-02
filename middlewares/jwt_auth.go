package middlewares

import (
	"context"
	"net/http"

	"github.com/farhanmobashir/twix"
	"github.com/golang-jwt/jwt/v5"
)

type TokenSource string

const (
	Header TokenSource = "header"
	Cookie TokenSource = "cookie"
)

type JWTConfig struct {
	SecretKey   []byte
	TokenSource TokenSource
	CookieName  string
}

func JWTAuth(config JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string

			switch config.TokenSource {
			case Header:
				tokenString = r.Header.Get("Authorization")
				if tokenString == "" {
					http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
					return
				}

				if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
					tokenString = tokenString[7:]
				}
			case Cookie:
				cookie, err := r.Cookie(config.CookieName)
				if err != nil {
					if err == http.ErrNoCookie {
						http.Error(w, "Missing authentication cookie", http.StatusUnauthorized)
					} else {
						http.Error(w, "Error retrieving authentication cookie", http.StatusInternalServerError)
					}
					return
				}
				tokenString = cookie.Value
			default:
				http.Error(w, "Invalid Token Source", http.StatusInternalServerError)
				return
			}

			// Token parsing
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return config.SecretKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Retrieve existing context or create a new one
			ctx, ok := r.Context().Value(twix.TwixContextKey).(*twix.Context)
			if !ok {
				ctx = &twix.Context{
					ResponseWriter: w,
					Request:        r,
					Params:         make(map[string]string),
				}
			}

			// Store the token claims in the context
			ctx.TokenClaims = token.Claims

			// Pass the context to the next handler
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), twix.TwixContextKey, ctx)))
		})
	}
}
