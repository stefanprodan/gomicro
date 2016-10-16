package main

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/unrolled/render"
)

var rnd *render.Render

func StartServer(appCtx AppContext) {

	rnd = render.New(render.Options{
		IndentJSON: true,
		Layout:     "layout",
	})

	r := chi.NewRouter()

	r.Use(middleware.CloseNotify)

	//prometheus metrics
	promRegister()
	r.Use(PromMiddleware)

	// chi middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// global middleware
	r.Use(AppMiddleware(appCtx))

	// routes
	r.Mount("/", homeRouter())
	r.Mount("/metrics", promRouter())

	//file server
	filesDir := filepath.Join(appCtx.WorkDir, "assets")
	r.FileServer("/assets", http.Dir(filesDir))

	log.Println("Starting HTTP server on port " + appCtx.Port + " work dir " + appCtx.WorkDir + ".")

	http.ListenAndServe(":"+appCtx.Port, r)
}
