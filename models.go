package main

import (
	"time"
)

type AppContext struct {
	Role      string    `json:"role"`
	Version   string    `json:"version"`
	Env       string    `json:"env"`
	Host      string    `json:"host"`
	Port      string    `json:"port"`
	WorkDir   string    `json:"workdir"`
	StartTime time.Time `json:"start"`
	Endpoints string    `json:"endpoints"`
}

type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Payload struct {
	Data string `json:"data"`
}
