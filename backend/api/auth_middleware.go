package api

import (
	"context"
	"net/http"
	"strings"

	"Log45/budget/backend/services"
)

type authenticatedUserIDKey struct{}

// AuthMiddleware requires a valid bearer token and adds its user ID to the
// request context for downstream handlers and services.
func AuthMiddleware(auth services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			scheme, token, found := strings.Cut(header, " ")
			if !found || !strings.EqualFold(scheme, "Bearer") || token == "" || strings.ContainsAny(token, " \t") {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateToken(token)
			if err != nil {
				http.Error(w, "invalid authentication token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), authenticatedUserIDKey{}, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthenticatedUserID returns the authenticated user's ID from a request context.
func AuthenticatedUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(authenticatedUserIDKey{}).(int64)
	return userID, ok && userID > 0
}
