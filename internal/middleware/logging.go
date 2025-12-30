package middleware

import(
	"log"
	"net/http"
	"time"
)

type responseLogger struct {
	http.ResponseWriter
	statusCode int
}


func (rl *responseLogger) WriteHeader(statusCode int) {
	rl.statusCode = statusCode
	rl.ResponseWriter.WriteHeader(statusCode)
}

func Logging(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &responseLogger{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lw, r)
		duration := time.Since(start)
		log.Printf(
			"%s %s %d %v",
			r.Method,
			r.URL.Path,
			lw.statusCode,
			duration,
		)
	})
}