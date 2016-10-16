package main

import (
	"net/http"

	"github.com/pressly/chi"
)

func homeRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", homeIndex)
	r.Get("/ping", homePing)
	r.Get("/panic", homePanic)
	return r
}

func promRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(Prom)
	r.Get("/", emptyHandler)
	return r
}
