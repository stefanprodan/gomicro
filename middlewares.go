package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AppMiddleware(app AppSettings) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "app", app)
			w.Header().Set("x-gomicro-version", app.Version)
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
			app, _ := r.Context().Value("app").(AppSettings)
			http_requests_total.WithLabelValues(app.Role, r.Method, r.URL.Path, status).Inc()
			//http_requests_latency.WithLabelValues(r.Method, r.URL.Path, status).Observe(took.Seconds())
		}
	})
}

func LogHttpErrorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		uri := r.RequestURI

		i := &interceptor{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(i, r)

		if !(100 <= i.statusCode && i.statusCode < 400) {
			log.Printf("%s %s (%d) %s", r.Method, uri, i.statusCode, time.Since(begin))
		}
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%s %s panic: %+v", r.Method, r.RequestURI, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func PromMetrics(next http.Handler) http.Handler {
	return promhttp.Handler()
}
