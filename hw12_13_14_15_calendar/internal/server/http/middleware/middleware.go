package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	logger  Logger
	handler http.Handler
}

func NewLoggingMiddleware(logger Logger, handlerToWrap http.Handler) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger, handler: handlerToWrap}
}

func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	m.handler.ServeHTTP(w, r)

	params := []any{
		"proto", r.Proto,
		"method", r.Method,
		"path", r.URL.Path,
		"user-agent", r.UserAgent(),
		"ip", r.RemoteAddr,
	}
	m.logger.Info(fmt.Sprintf("Processed request in %v", time.Since(start)), params...)
}
