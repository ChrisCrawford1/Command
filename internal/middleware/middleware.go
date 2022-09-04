package middleware

import (
	"context"
	"encoding/json"
	"github.com/ChrisCrawford1/Command/internal/auth"
	"github.com/ChrisCrawford1/Command/internal/responses"
	"net/http"
	"strings"
)

// ContentTypeMiddleware - Sets the response type on all requests to be application/json
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// ValidateJwtToken - Check that the user exists and set their id into the context
func ValidateJwtToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			splitToken := strings.Split(r.Header["Authorization"][0], "Bearer ")[1]

			isValid, claims, err := auth.ValidateAccessToken(splitToken)

			if err != nil {
				//http.Error(w, "Invalid credentials", 401)
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(responses.Error{Message: "Invalid credentials"})
				return
			}

			if isValid {
				ctx := context.WithValue(r.Context(), "userId", claims["userId"])

				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Not authorized", 401)
			}
		} else {
			http.Error(w, "Not authorized", 401)
		}
	})
}
