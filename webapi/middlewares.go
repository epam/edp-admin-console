package webapi

import (
	"net/http"

	"go.uber.org/zap"
)

func WithLoggerMw(logger *zap.Logger) func(next http.Handler) http.Handler {
	mw := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			r = r.WithContext(ContextWithLogger(ctx, logger))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return mw
}
