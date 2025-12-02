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
		serverURL = "http://192.168.1.143:8080/api/collect"
	}

	// Leer timeout del .env o usar 120 segundos por defecto
	timeout := 120 * time.Second
	if timeoutStr := os.Getenv("TIMEOUT"); timeoutStr != "" {
		if parsedTimeout, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = parsedTimeout
		}
	}

	return &Config{
		ServerURL:      serverURL,
		Timeout:        timeout,
		EnableDebug:    os.Getenv("DEBUG") == "1",
		CollectorCount: 5, // Número de colectores paralelos
	}
}
