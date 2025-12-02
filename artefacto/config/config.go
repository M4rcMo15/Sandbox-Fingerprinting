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
	XSSAudit       bool
	CallbackServer string
}

func Load() *Config {
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://54.37.226.179/api/collect"
	}

	// Leer timeout del .env o usar 120 segundos por defecto
	timeout := 120 * time.Second
	if timeoutStr := os.Getenv("TIMEOUT"); timeoutStr != "" {
		if parsedTimeout, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = parsedTimeout
		}
	}

	// Callback server para XSS (por defecto el mismo que ServerURL pero sin /api/collect)
	callbackServer := os.Getenv("CALLBACK_SERVER")
	if callbackServer == "" {
		callbackServer = "http://54.37.226.179"
	}

	return &Config{
		ServerURL:      serverURL,
		Timeout:        timeout,
		EnableDebug:    os.Getenv("DEBUG") == "1",
		CollectorCount: 5, // Número de colectores paralelos
		XSSAudit:       os.Getenv("XSS_AUDIT") == "true",
		CallbackServer: callbackServer,
	}
}
