package middleware

import (
	"context"
	"net/http"
	"strings"

	"task-flow/internal/httpx"
	"task-flow/internal/pkg/jwt"
)

type ctxKey string

type AuthMiddleware struct {
	JWT *jwt.JWT
}

func NewAuthMiddleware(j *jwt.JWT) *AuthMiddleware {
	return &AuthMiddleware{
		JWT: j,
	}
}

const userIDKey ctxKey = "user_id"

func RequireAccessJWT(authSvc *AuthMiddleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			parts := strings.SplitN(h, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				httpx.Error(w, http.StatusUnauthorized, "missing/invalid authorization")
				return
			}

			userID, err := authSvc.JWT.Verify(parts[1])
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserID(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey).(string)
	return v
}
