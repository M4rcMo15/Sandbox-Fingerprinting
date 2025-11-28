package config

import (
	"os"
	"time"
)

type Config struct {
	ServerURL      string
	Timeout        time.Duration
	EnableDebug    bool
	CollectorCount int
}

func Load() *Config {
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://192.168.1.143:8080/content"
	}

	return &Config{
		ServerURL:      serverURL,
		Timeout:        30 * time.Second,
		EnableDebug:    os.Getenv("DEBUG") == "1",
		CollectorCount: 5, // Número de colectores paralelos
	}
}
