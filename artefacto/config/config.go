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
		serverURL = "https://releases.life/api/collect"
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
		callbackServer = "https://releases.life"
	}

	// XSS Audit: Activado por defecto si no se especifica lo contrario
	xssAudit := true // Por defecto activado
	if xssAuditEnv := os.Getenv("XSS_AUDIT"); xssAuditEnv != "" {
		xssAudit = xssAuditEnv == "true"
	}

	return &Config{
		ServerURL:      serverURL,
		Timeout:        timeout,
		EnableDebug:    os.Getenv("DEBUG") == "1",
		CollectorCount: 5, // Número de colectores paralelos
		XSSAudit:       xssAudit,
		CallbackServer: callbackServer,
	}
}
