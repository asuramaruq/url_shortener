package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(slog.String("component", "middleware/logger"))

		log.Info("logger middleware enabled")

		// actual handler code
		fn := func(w http.ResponseWriter, r *http.Request) {
			// initial request data
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			// abstraction layer on top of the default responsewriter to get info on the response
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// tracking time it took for request processing
			t1 := time.Now()
			// using defer to log the completion time, as it will run after the request processing is complete
			defer func() {
				entry.Info(
					"request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()
			// transferring ownership to next handler in the middleware chain
			next.ServeHTTP(ww, r)
		}
		// return written handler wrapping it into handlerfunc
		return http.HandlerFunc(fn)
	}
}
