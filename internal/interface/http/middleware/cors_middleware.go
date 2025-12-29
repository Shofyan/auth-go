package middleware

import (
	"net/http"
)

// CORSMiddleware handles CORS
type CORSMiddleware struct{}

// NewCORSMiddleware creates a new CORS middleware
func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}

// Handle handles CORS headers
func (m *CORSMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
