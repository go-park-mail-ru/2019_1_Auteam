package main

import (
    "context"
    "net/http"
    "os"
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
            w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_IP"))
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusOK)
                return
            }
            next.ServeHTTP(w, r)
        })
}