package pkg

import (
	"log"
	"os"
	"strings"
)

var damiURL string = "http://localhost:8001"

const damiURLKey = "DAMI_URL"

type Env struct {
	DamiURL string
}

func (e *Env) ParseEnv() {
	if v := os.Getenv(damiURLKey); strings.TrimSpace(v) != "" {
		log.Println("updating env key")
		damiURL = v
	}
	e.DamiURL = damiURL
	log.Println("env is ", e.DamiURL)
}
