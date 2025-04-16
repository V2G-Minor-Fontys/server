package middleware

import (
	"context"
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/V2G-Minor-Fontys/server/pkg/jwt"
	"net/http"
	"strings"
)

const IdentityIDKey string = "identityID"

func AuthVerifier(cfg *config.Jwt) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				httpx.ProblemResponseWithJSON(w, httpx.Unauthorized("Missing or invalid Authorization header", r.RequestURI))
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwt.VerifyAccessToken(tokenStr, cfg)
			if err != nil {
				httpx.ProblemResponseWithJSON(w, httpx.Unauthorized("Invalid or expired token", r.RequestURI))
				return
			}

			sub, _ := claims.GetSubject()
			ctx := context.WithValue(r.Context(), IdentityIDKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
