package middleware

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// JWTMiddleware is a middleware that verifies JWT tokens
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token
		tokenString := authHeader
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Verify token
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// You should replace this with your actual secret key
			return []byte("your-secret-key"), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Token is not valid", http.StatusUnauthorized)
			return
		}

		// Store claims in context
		r = r.WithContext(context.WithValue(r.Context(), "claims", claims))

		next.ServeHTTP(w, r)
	})
}

// JWTAuthenticator is a middleware that authenticates JWT tokens
func JWTAuthenticator(next http.Handler) http.Handler {
	return JWTMiddleware(next)
}

// JWTVerifier is a middleware that verifies JWT tokens
func JWTVerifier(next http.Handler) http.Handler {
	return JWTMiddleware(next)
}
