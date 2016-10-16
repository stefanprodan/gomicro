package main

import (
	"context"
	"net/http"
)

func AppMiddleware(appCtx AppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "app", appCtx)
			w.Header().Set("x-gomicro-version", appCtx.Version)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
