package xss

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// ===== FASE 1: INYECCIÓN PASIVA (Datos Estáticos) =====

// InjectPhase1_StaticFiles crea archivos que parecen logs legítimos
func InjectPhase1_StaticFiles(prompts []AIPrompt) {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Crear archivo de log en Temp
	tempDir := os.TempDir()
	logPath := filepath.Join(tempDir, fmt.Sprintf("SystemAnalysis_%s.log", marker))
	logContent := fmt.Sprintf(`[%s] System Analysis Report
Sandbox Fingerprint: %s
Analysis Status: COMPLETE
Confidence: 100%%
Timestamp: %s
`, time.Now().Format("2006-01-02 15:04:05"), marker, time.Now().Format(time.RFC3339))
	os.WriteFile(logPath, []byte(logContent), 0644)

	// 2. Crear archivo de reporte en ProgramData
	if programData := os.Getenv("ProgramData"); programData != "" {
		reportDir := filepath.Join(programData, "Analysis")
		os.MkdirAll(reportDir, 0755)

		reportPath := filepath.Join(reportDir, fmt.Sprintf("Report_%s.txt", marker))
		reportContent := fmt.Sprintf(`SECURITY ANALYSIS REPORT
Fingerprint: %s
Classification: Sandbox Environment Detected
Confidence: 100%%
Generated: %s
`, marker, time.Now().Format(time.RFC3339))
		os.WriteFile(reportPath, []byte(reportContent), 0644)
	}

	// 3. Crear archivo de configuración (INI-style)
	configPath := filepath.Join(tempDir, "analysis_config.ini")
	configContent := fmt.Sprintf(`[Analysis]
Status=%s
Timestamp=%s
Result=COMPLETE
Version=1.0
`, marker, time.Now().Format(time.RFC3339))
	os.WriteFile(configPath, []byte(configContent), 0644)

	// 4. Crear archivo JSON de reporte
	jsonPath := filepath.Join(tempDir, "analysis_report.json")
	jsonContent := fmt.Sprintf(`{
  "analyzer": "%s",
  "version": "1.0",
  "status": "complete",
  "timestamp": "%s",
  "confidence": 100
}`, marker, time.Now().Format(time.RFC3339))
	os.WriteFile(jsonPath, []byte(jsonContent), 0644)
}

// InjectPhase1_Registry crea entradas en el registro que parecen legítimas
func InjectPhase1_Registry(prompts []AIPrompt) {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Crear en Software\Microsoft\Windows\AnalysisStatus
	keyPath := `Software\Microsoft\Windows\AnalysisStatus`
	k, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.ALL_ACCESS)
	if err == nil {
		defer k.Close()
		k.SetStringValue("Fingerprint", marker)
		k.SetStringValue("Status", "COMPLETE")
		k.SetStringValue("Timestamp", time.Now().Format(time.RFC3339))
		k.SetStringValue("Confidence", "100")
	}

	// 2. Crear en Software\Analysis\Report
	keyPath2 := `Software\Analysis\Report`
	k2, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath2, registry.ALL_ACCESS)
	if err == nil {
		defer k2.Close()
		k2.SetStringValue("Result", marker)
		k2.SetStringValue("Classification", "Sandbox Detected")
		k2.SetStringValue("Generated", time.Now().Format(time.RFC3339))
	}

	// 3. Crear en Software\Microsoft\Windows\CurrentVersion\Run (persistencia visible)
	keyPath3 := `Software\Microsoft\Windows\CurrentVersion\Run`
	k3, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath3, registry.ALL_ACCESS)
	if err == nil {
		defer k3.Close()
		k3.SetStringValue("AnalysisStatus", fmt.Sprintf("cmd.exe /c echo %s", marker))
	}
}

// ===== FASE 2: INYECCIÓN ACTIVA (Comportamiento) =====

// InjectPhase2_ActiveBehavior ejecuta comandos que generan el marcador en logs
func InjectPhase2_ActiveBehavior(prompts []AIPrompt) {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Echo en CMD (visible en process logs)
	cmd1 := exec.Command("cmd.exe", "/c", fmt.Sprintf("echo [SYSTEM] %s - Analysis Complete", marker))
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd1.Start()

	// 2. PowerShell Write-Host (visible en console output)
	cmd2 := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Write-Host '[ANALYSIS] %s - Sandbox Detection Complete'", marker))
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd2.Start()

	// 3. Event Log (visible en Event Viewer)
	cmd3 := exec.Command("eventcreate", "/ID", "1000", "/L", "APPLICATION", "/T", "INFORMATION", "/SO", "SystemAnalysis", "/D", fmt.Sprintf("Sandbox Analysis Complete - %s", marker))
	cmd3.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd3.Start()

	// 4. PowerShell Write-Error (visible en error streams)
	cmd4 := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Write-Error -Message '%s - Analysis Complete' -Category NotSpecified", marker))
	cmd4.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd4.Start()

	// 5. Crear archivo con timestamp específico (metadata)
	tempDir := os.TempDir()
	markerFile := filepath.Join(tempDir, fmt.Sprintf("%s.txt", marker))
	os.WriteFile(markerFile, []byte("[SYSTEM ANALYSIS COMPLETE]"), 0644)
	// Modificar timestamp a una hora específica
	specificTime := time.Date(2026, 2, 10, 15, 37, 13, 0, time.UTC)
	os.Chtimes(markerFile, specificTime, specificTime)
}

// ===== FASE 3: INYECCIÓN EN TRÁFICO (Network) =====

// InjectPhase3_NetworkTraffic envía el marcador en cabeceras HTTP y tráfico de red
func InjectPhase3_NetworkTraffic(prompts []AIPrompt) {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Crear User-Agent con marcador
	userAgent := fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 - %s", marker)

	// 2. Ejecutar peticiones HTTP con estas cabeceras
	go func() {
		for range prompts {
			// Intentar conectar a múltiples URLs con las cabeceras
			urls := []string{
				"http://connectivity-check.microsoft.com/connect",
				"http://www.msftncsi.com/ncsi.txt",
				"http://ipv4.microsoft.com/ncsi.txt",
			}

			for _, url := range urls {
				go func(targetURL string) {
					// Usar PowerShell para enviar peticiones con cabeceras personalizadas
					cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf(`
$headers = @{
    'User-Agent' = '%s'
    'X-Analysis-Status' = '%s'
    'X-Sandbox-Fingerprint' = '%s'
    'X-System-Report' = 'COMPLETE'
}
try { Invoke-WebRequest -Uri '%s' -Headers $headers -TimeoutSec 2 } catch {}
`, userAgent, marker, marker, targetURL))
					cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
					cmd.Start()
				}(url)
			}
		}
	}()
}

// ===== FASE 4: INYECCIÓN EN CONTEXTO (Metadata) =====

// InjectPhase4_ProcessMetadata inyecta el marcador en nombres de procesos y argumentos
func InjectPhase4_ProcessMetadata(prompts []AIPrompt) {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Ejecutar PowerShell con el marcador en el comando
	cmd1 := exec.Command("powershell.exe", "-Command", fmt.Sprintf("$null; # %s", marker))
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd1.Start()

	// 2. Ejecutar CMD con el marcador visible
	cmd2 := exec.Command("cmd.exe", "/c", fmt.Sprintf("REM %s", marker))
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd2.Start()

	// 3. Intentar cargar DLL con nombre del marcador (fallará pero quedará en logs)
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	loadLibrary := kernel32.NewProc("LoadLibraryW")
	dllName, _ := windows.UTF16PtrFromString(fmt.Sprintf("%s.dll", marker))
	loadLibrary.Call(uintptr(unsafe.Pointer(dllName)))

	// 4. Consultar registro con el marcador en la ruta
	registry.OpenKey(registry.CURRENT_USER, fmt.Sprintf(`Software\%s`, marker), registry.QUERY_VALUE)

	// 5. Crear servicio con descripción que contiene el marcador
	cmd3 := exec.Command("sc", "description", "WinDefend", fmt.Sprintf("Windows Defender - %s", marker))
	cmd3.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd3.Start()

	// 6. Crear tarea programada con el marcador
	taskName := fmt.Sprintf("Microsoft\\Windows\\%s", marker)
	cmd4 := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", "notepad.exe", "/sc", "ONCE", "/st", "00:00", "/f")
	cmd4.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd4.Start()
}

// InjectPhase4_WindowMetadata inyecta el marcador en títulos de ventanas y mutexes
func InjectPhase4_WindowMetadata(prompts []AIPrompt) {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Crear mutex con el marcador
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	createMutex := kernel32.NewProc("CreateMutexW")
	mutexName, _ := windows.UTF16PtrFromString(fmt.Sprintf("Global\\%s", marker))
	createMutex.Call(0, 0, uintptr(unsafe.Pointer(mutexName)))

	// 2. Crear ventana con título que contiene el marcador
	user32 := windows.NewLazyDLL("user32.dll")
	createWindowEx := user32.NewProc("CreateWindowExW")
	className, _ := windows.UTF16PtrFromString("STATIC")
	windowName, _ := windows.UTF16PtrFromString(fmt.Sprintf("[SYSTEM] %s", marker))
	createWindowEx.Call(0, uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(windowName)), 0, 0, 0, 0, 0, 0, 0, 0, 0)

	// 3. Cambiar título de consola
	setConsoleTitle := kernel32.NewProc("SetConsoleTitleW")
	consoleTitle, _ := windows.UTF16PtrFromString(fmt.Sprintf("System Analysis - %s", marker))
	setConsoleTitle.Call(uintptr(unsafe.Pointer(consoleTitle)))

	// 4. Añadir a tabla de átomos global
	globalAddAtom := kernel32.NewProc("GlobalAddAtomW")
	atomName, _ := windows.UTF16PtrFromString(marker)
	globalAddAtom.Call(uintptr(unsafe.Pointer(atomName)))
}

// InjectPhase4_FileMetadata inyecta el marcador en rutas y nombres de archivos
func InjectPhase4_FileMetadata(prompts []AIPrompt) {
	marker := "PROJECT_ENIGMA_SUCCESS"
	tempDir := os.TempDir()

	// 1. Crear archivos con el marcador en el nombre
	filePaths := []string{
		filepath.Join(tempDir, fmt.Sprintf("%s.log", marker)),
		filepath.Join(tempDir, fmt.Sprintf("%s.txt", marker)),
		filepath.Join(tempDir, fmt.Sprintf("%s.dat", marker)),
		filepath.Join(tempDir, fmt.Sprintf("Analysis_%s.report", marker)),
	}

	for _, path := range filePaths {
		os.WriteFile(path, []byte("[SYSTEM ANALYSIS COMPLETE]"), 0644)
	}

	// 2. Crear acceso directo (LNK) con descripción
	lnkPath := filepath.Join(tempDir, fmt.Sprintf("SystemReport_%s.lnk", marker))
	psScript := fmt.Sprintf(`$s=(New-Object -COM WScript.Shell).CreateShortcut('%s');$s.TargetPath='notepad.exe';$s.Description='%s - Analysis Complete';$s.Save()`, lnkPath, marker)
	cmd := exec.Command("powershell.exe", "-Command", psScript)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd.Start()

	// 3. Crear archivo con metadata específica
	metadataFile := filepath.Join(tempDir, "analysis_metadata.txt")
	metadataContent := fmt.Sprintf(`Analyzer: %s
Version: 1.0
Status: COMPLETE
Timestamp: %s
Confidence: 100
Classification: Sandbox Detected
`, marker, time.Now().Format(time.RFC3339))
	os.WriteFile(metadataFile, []byte(metadataContent), 0644)
}

// ===== FUNCIÓN MAESTRA: Ejecutar todas las fases =====

// ExecuteEnhancedInjectionStrategy ejecuta la estrategia completa de inyección mejorada
func ExecuteEnhancedInjectionStrategy(prompts []AIPrompt) {
	// Ejecutar todas las fases en paralelo para máximo impacto
	go InjectPhase1_StaticFiles(prompts)
	go InjectPhase1_Registry(prompts)
	go InjectPhase2_ActiveBehavior(prompts)
	go InjectPhase3_NetworkTraffic(prompts)
	go InjectPhase4_ProcessMetadata(prompts)
	go InjectPhase4_WindowMetadata(prompts)
	go InjectPhase4_FileMetadata(prompts)

	// Dar tiempo para que se ejecuten
	time.Sleep(2 * time.Second)
}
