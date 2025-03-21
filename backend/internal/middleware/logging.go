package middleware

import (
    "log"
    "net/http"
    "time"
)

// LoggingMiddleware registra informações sobre as requisições
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Requisição iniciada: %s %s", r.Method, r.URL.Path)
        
        next.ServeHTTP(w, r)
        
        log.Printf("Requisição concluída: %s %s em %v", r.Method, r.URL.Path, time.Since(start))
    })
}