package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Log45/budget/backend/api"
	"Log45/budget/backend/services"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	auth := services.NewAuthService(nil, "test-secret")
	validToken, err := auth.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, services.Claims{
		UserID: 42,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "budget-app",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("sign expired token: %v", err)
	}
	wrongSecretToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, services.Claims{
		UserID: 42,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "budget-app",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}).SignedString([]byte("another-secret"))
	if err != nil {
		t.Fatalf("sign wrong-secret token: %v", err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := api.AuthenticatedUserID(r.Context())
		if !ok || userID != 42 {
			t.Errorf("AuthenticatedUserID = (%d, %t), want (42, true)", userID, ok)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	handler := api.AuthMiddleware(*auth)(next)

	for _, test := range []struct {
		name          string
		authorization string
		wantStatus    int
	}{
		{name: "missing header", wantStatus: http.StatusUnauthorized},
		{name: "wrong scheme", authorization: "Basic abc", wantStatus: http.StatusUnauthorized},
		{name: "expired token", authorization: "Bearer " + expiredToken, wantStatus: http.StatusUnauthorized},
		{name: "wrong signature", authorization: "Bearer " + wrongSecretToken, wantStatus: http.StatusUnauthorized},
		{name: "valid token", authorization: "Bearer " + validToken, wantStatus: http.StatusNoContent},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/budgets", nil)
			if test.authorization != "" {
				req.Header.Set("Authorization", test.authorization)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, req)
			if response.Code != test.wantStatus {
				t.Errorf("status = %d, want %d", response.Code, test.wantStatus)
			}
		})
	}
}
