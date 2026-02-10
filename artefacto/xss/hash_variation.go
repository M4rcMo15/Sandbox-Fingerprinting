package xss

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

// HashVariationConfig contiene datos para variar el hash del binario
type HashVariationConfig struct {
	CampaignID    string
	VariationID   string
	Timestamp     string
	RandomPadding string
}

// GenerateHashVariation crea una variación única para cada compilación
func GenerateHashVariation() *HashVariationConfig {
	// Generar ID de campaña basado en timestamp
	campaignID := fmt.Sprintf("ENIGMA_%d", time.Now().Unix())

	// Generar ID de variación aleatorio
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	variationID := hex.EncodeToString(randomBytes)

	// Timestamp actual
	timestamp := time.Now().Format(time.RFC3339)

	// Generar padding aleatorio (256-512 bytes)
	paddingSize := 256 + int(randomBytes[0])%256
	paddingBytes := make([]byte, paddingSize)
	rand.Read(paddingBytes)
	randomPadding := hex.EncodeToString(paddingBytes)

	return &HashVariationConfig{
		CampaignID:    campaignID,
		VariationID:   variationID,
		Timestamp:     timestamp,
		RandomPadding: randomPadding,
	}
}

// GetHashVariationFromEnv obtiene la variación del hash desde variables de entorno
func GetHashVariationFromEnv() *HashVariationConfig {
	config := &HashVariationConfig{
		CampaignID:    os.Getenv("CAMPAIGN_ID"),
		VariationID:   os.Getenv("VARIATION_ID"),
		Timestamp:     os.Getenv("VARIATION_TIMESTAMP"),
		RandomPadding: os.Getenv("RANDOM_PADDING"),
	}

	// Si no hay variables de entorno, generar nuevas
	if config.CampaignID == "" {
		return GenerateHashVariation()
	}

	return config
}

// EmbedHashVariation embebe la variación en el binario (como strings globales)
// Esto asegura que cada compilación tenga un hash diferente
var (
	// Estos strings se embeben en el binario y varían por compilación
	EmbeddedCampaignID    = "ENIGMA_DEFAULT"
	EmbeddedVariationID   = "default_variation"
	EmbeddedTimestamp     = "2026-02-10T00:00:00Z"
	EmbeddedRandomPadding = "default_padding"
)

// InitializeHashVariation inicializa la variación del hash
func InitializeHashVariation() {
	config := GetHashVariationFromEnv()

	EmbeddedCampaignID = config.CampaignID
	EmbeddedVariationID = config.VariationID
	EmbeddedTimestamp = config.Timestamp
	EmbeddedRandomPadding = config.RandomPadding
}

// GetCampaignID retorna el ID de campaña
func GetCampaignID() string {
	return EmbeddedCampaignID
}

// GetVariationID retorna el ID de variación
func GetVariationID() string {
	return EmbeddedVariationID
}

// GetVariationTimestamp retorna el timestamp de variación
func GetVariationTimestamp() string {
	return EmbeddedTimestamp
}
