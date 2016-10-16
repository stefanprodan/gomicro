package main

import (
	"net/http"
)

func homeIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	app, ok := ctx.Value("app").(AppContext)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if r.Header.Get("Content-Type") == "application/json" {
		rnd.JSON(w, http.StatusOK, app)
	} else {
		rnd.HTML(w, http.StatusOK, "home", app)
	}
}

func homePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func homePanic(w http.ResponseWriter, r *http.Request) {
	panic("ERROR!!!")
}
func emptyHandler(w http.ResponseWriter, r *http.Request) {

}
