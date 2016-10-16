package main

import (
	"log"
	"os"
	"time"
)

func main() {
	var (
		role       = os.Getenv("ROLE") // worker, proxy, ping
		env        = os.Getenv("ENV")  // DEBUG, DEV, STG, PROD
		port       = os.Getenv("PORT")
		endpoints  = os.Getenv("ENDPOINTS")
		host, _    = os.Hostname()
		workDir, _ = os.Getwd()
	)

	// defaults
	if env == "" {
		role = "proxy"
		env = "DEBUG"
		port = "3001"
		endpoints = "http://localhost:3001"
	}

	// reading version from file
	version, err := ParseVersionFile("VERSION")
	if err != nil {
		log.Fatal(err)
	}

	appCtx := AppContext{
		Role:      role,
		Version:   version,
		Env:       env,
		Host:      host,
		Port:      port,
		WorkDir:   workDir,
		StartTime: time.Now(),
		Endpoints: endpoints,
	}

	log.Println("Starting gomicro v" + appCtx.Version + " on " + appCtx.Host + ":" + appCtx.Port + " in " + appCtx.Env + " mode.")

	// start HTTP server
	StartServer(appCtx)
}
