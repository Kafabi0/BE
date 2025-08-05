package middleware

import (
	"context"
	"net/http"
	"strings"

	"loginapp/handlers"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("supersecretkey")

type contextKey string

const ClaimsContextKey contextKey = "claims"

// JWTMiddleware validasi token JWT dan simpan claims ke context
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token tidak ditemukan", http.StatusUnauthorized)
			return
		}

		// Format header harus: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Format Authorization header salah", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		claims := &handlers.Claims{}

		// ekstra dan validasi
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token tidak valid", http.StatusUnauthorized)
			return
		}

		// Simpan claims ke context request
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx)) // lanjutkan ke handler berikutnya
	})
}

// RoleAuthorization middleware cek role yang diizinkan
func RoleAuthorization(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsContextKey).(*handlers.Claims) // ambil claims dari context
			if !ok {
				http.Error(w, "Claims tidak ditemukan", http.StatusUnauthorized)
				return
			}

			allowed := false
			// Cek apakah role user ada di daftar allowedRoles
			for _, role := range allowedRoles {
				if claims.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Akses ditolak: role tidak sesuai", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
