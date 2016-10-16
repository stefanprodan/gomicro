package main

import (
	"net/http"

	//	"github.com/unrolled/render"
)

func homeIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	app, ok := ctx.Value("app").(AppContext)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	rnd.HTML(w, http.StatusOK, "home", app)
}

func homeIndexJson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	app, ok := ctx.Value("app").(AppContext)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	rnd.JSON(w, http.StatusOK, app)
}

func homePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func homePanic(w http.ResponseWriter, r *http.Request) {
	panic("ERROR!!!")
}
