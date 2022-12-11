package middleware

import (
	"net/http"

	"github.com/cheebz/go-auth-helpers"
)

func NewAuthMiddleware(endpoint string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := auth.Authenticate(w, r, endpoint)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
