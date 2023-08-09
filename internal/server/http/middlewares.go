package http

import (
	"github.com/romandnk/advertisement/internal/logger"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func loggingMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lrw := negroni.NewResponseWriter(w)

			start := time.Now()

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			log.Info("Request info HTTP",
				zap.String("client ip", r.RemoteAddr),
				zap.String("method", r.Method),
				zap.String("method path", r.URL.Path),
				zap.String("HTTP version", r.Proto),
				zap.String("status code", strconv.Itoa(lrw.Status())),
				zap.String("processing time", duration.String()),
			)
		})
	}
}
