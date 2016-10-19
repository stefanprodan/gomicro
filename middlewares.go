package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func PromMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//begin := time.Now()

		interceptor := &interceptor{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(interceptor, r)

		var (
			status = strconv.Itoa(interceptor.statusCode)
			//took   = time.Since(begin)
		)

		// ignore websockets and the /metrics endpoint
		if !isWSRequest(r) && !isPromRequest(r) {
			http_requests_total.WithLabelValues(r.Method, r.URL.Path, status).Inc()
			//http_requests_latency.WithLabelValues(r.Method, r.URL.Path, status).Observe(took.Seconds())
		}
	})
}

func PromMetrics(next http.Handler) http.Handler {
	return promhttp.Handler()
}
