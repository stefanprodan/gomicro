package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	var (
		role       = os.Getenv("ROLE") // worker, proxy, monitor
		env        = os.Getenv("ENV")  // DEBUG, DEV, STG, PROD
		port       = os.Getenv("PORT")
		endpoints  = os.Getenv("ENDPOINTS")
		host, _    = os.Hostname()
		workDir, _ = os.Getwd()
	)

	// defaults
	if env == "" {
		role = "monitor"
		env = "DEBUG"
		port = "3001"
		endpoints = "http://localhost:3001"
	}

	// reading version from file
	version, err := ParseVersionFile("VERSION")
	if err != nil {
		log.Fatal(err)
	}

	app := AppSettings{
		Role:         role,
		Version:      version,
		Env:          env,
		Host:         host,
		Port:         port,
		WorkDir:      workDir,
		StartTime:    time.Now(),
		Endpoints:    endpoints,
		PingInterval: 1000,
	}

	log.Println("Starting gomicro v" + app.Version + " on " + app.Host + ":" + app.Port + " in " + app.Env + " mode.")

	PromRegister()

	// start services
	go StartServer(app)

	if role == "monitor" {
		go StartHealthCheck(app)
	}

	// block
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt, os.Kill)
	osSignal := <-osChan

	log.Printf("Exiting! OS signal: %v", osSignal)
}
