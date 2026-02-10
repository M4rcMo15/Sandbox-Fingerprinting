package config

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type Config struct {
	ServerURL      string
	Timeout        time.Duration
	EnableDebug    bool
	CollectorCount int
	XSSAudit       bool
	CallbackServer string
	TargetSandbox  string
}

func Load(targetSandbox string) *Config {
	// 1. Cargar variables desde .env si existe
	loadEnvFile()

	// 2. Establecer valores por defecto si no existen en variables de entorno
	setDefaults()

	serverURL := os.Getenv("SERVER_URL")

	// Leer timeout del .env o usar 120 segundos por defecto
	timeout := 120 * time.Second
	if timeoutStr := os.Getenv("TIMEOUT"); timeoutStr != "" {
		if parsedTimeout, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = parsedTimeout
		}
	}

	// Callback server para XSS (por defecto el mismo que ServerURL pero sin /api/collect)
	callbackServer := os.Getenv("CALLBACK_SERVER")

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
		TargetSandbox:  targetSandbox,
	}
}

// loadEnvFile carga las variables de entorno desde el archivo .env
func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Ignorar líneas vacías y comentarios
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Separar clave=valor
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
}

func setDefaults() {
	defaults := map[string]string{
		"SERVER_URL":      "https://releases.life/api/collect",
		"CALLBACK_SERVER": "https://xss.releases.life",
		"XSS_AUDIT":       "true",
		"TIMEOUT":         "120s",
	}
	for k, v := range defaults {
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}
