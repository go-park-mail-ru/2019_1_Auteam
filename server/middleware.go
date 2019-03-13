package server

import (
	"context"
	"net/http"
	"log"
)

func (s *Server) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			isAuth := false
			ctx := r.Context()
			cookie, err := r.Cookie("sessionID")
			if err != http.ErrNoCookie {
				userId, err := s.CheckSession(cookie.Value)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				isAuth = true
				ctx = context.WithValue(ctx, "userId", userId)
				ctx = context.WithValue(ctx, "sessionID", cookie.Value)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, "isAuth", isAuth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func SetCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "dev.mycodestory.ru")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Logs: %s: %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
}