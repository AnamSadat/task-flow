package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// ANSI colors
const (
	reset  = "\033[0m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	cyan   = "\033[36m"
	gray   = "\033[90m"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return cyan
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		fmt.Printf("%s%s%s %s %s%d%s %s%.2fms%s\n",
			cyan, r.Method, reset,
			r.URL.Path,
			statusColor(rw.status), rw.status, reset,
			gray, float64(time.Since(start).Microseconds())/1000, reset,
		)
	})
}
