package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/auth"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Authhandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authheader := r.Header.Get("Authorization")
		if authheader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("missing token")
			return
		}
		parts := strings.SplitAfter(authheader, " ")
		if len(parts) != 2 || parts[0] == "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("invalid format")
			return
		}

		user_id, err := auth.ValidateToken(parts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("invalid or expired token")
		}
		ctx := context.WithValue(r.Context(), UserIDKey, user_id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
