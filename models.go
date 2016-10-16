package main

type AppContext struct {
	Role    string `json:"role"`
	Version string `json:"version"`
	Env     string `json:"env"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	WorkDir string `json:"workdir"`
}
