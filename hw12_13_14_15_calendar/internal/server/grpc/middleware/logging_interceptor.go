package middleware

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func InterceptorLogger(l Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]any, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			f = append(f, key, value)
		}

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg, f...)
		case logging.LevelInfo:
			l.Info(msg, f...)
		case logging.LevelWarn:
			l.Warn(msg, f...)
		case logging.LevelError:
			l.Error(msg, f...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
