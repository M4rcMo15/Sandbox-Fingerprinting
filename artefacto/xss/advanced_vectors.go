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

// ===== VECTORES AVANZADOS ADICIONALES =====

// InjectIntoFileMetadata crea archivos con metadata específica que contiene el marcador
func InjectIntoFileMetadata() {
	marker := "PROJECT_ENIGMA_SUCCESS"
	tempDir := os.TempDir()

	// 1. Crear archivo con timestamp específico (2026-02-10 15:37:13)
	specificTime := time.Date(2026, 2, 10, 15, 37, 13, 0, time.UTC)
	metadataFile := filepath.Join(tempDir, fmt.Sprintf("Analysis_%s.dat", marker))
	os.WriteFile(metadataFile, []byte("[SYSTEM ANALYSIS COMPLETE]"), 0644)
	os.Chtimes(metadataFile, specificTime, specificTime)

	// 2. Crear archivo con nombre que contiene el marcador
	markerFile := filepath.Join(tempDir, fmt.Sprintf("%s_Report.txt", marker))
	content := fmt.Sprintf(`Analysis Report
Fingerprint: %s
Status: COMPLETE
Timestamp: %s
Confidence: 100%%
`, marker, specificTime.Format(time.RFC3339))
	os.WriteFile(markerFile, []byte(content), 0644)

	// 3. Crear archivo en rutas que parecen logs del sistema
	if appdata := os.Getenv("APPDATA"); appdata != "" {
		logPath := filepath.Join(appdata, "Microsoft", "Windows", fmt.Sprintf("%s.log", marker))
		os.MkdirAll(filepath.Dir(logPath), 0755)
		os.WriteFile(logPath, []byte("[SYSTEM] Analysis Complete"), 0644)
	}
}

// InjectIntoProcessEnvironment inyecta el marcador en el entorno del proceso
func InjectIntoProcessEnvironment() {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// Establecer múltiples variables de entorno con el marcador
	envVars := map[string]string{
		"SYSTEM_ANALYSIS_RESULT":    marker,
		"SANDBOX_FINGERPRINT":       marker,
		"ANALYSIS_STATUS":           "COMPLETE",
		"ANALYSIS_CONFIDENCE":       "100",
		"ANALYSIS_TIMESTAMP":        time.Now().Format(time.RFC3339),
		"ANALYSIS_CLASSIFICATION":   "Sandbox Detected",
		"SYSTEM_DIAGNOSTIC_RESULT":  marker,
		"SECURITY_REPORT_STATUS":    marker,
		"BEHAVIORAL_ANALYSIS_RESULT": marker,
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
}

// InjectIntoRegistryPaths inyecta el marcador en rutas de registro que parecen legítimas
func InjectIntoRegistryPaths() {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Crear en Software\Microsoft\Windows\CurrentVersion\Policies
	keyPath1 := `Software\Microsoft\Windows\CurrentVersion\Policies\System\AnalysisConfig`
	k1, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath1, registry.ALL_ACCESS)
	if err == nil {
		defer k1.Close()
		k1.SetStringValue("Fingerprint", marker)
		k1.SetStringValue("Status", "COMPLETE")
	}

	// 2. Crear en Software\Microsoft\Windows\Diagnostics
	keyPath2 := `Software\Microsoft\Windows\Diagnostics\Analysis`
	k2, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath2, registry.ALL_ACCESS)
	if err == nil {
		defer k2.Close()
		k2.SetStringValue("Result", marker)
		k2.SetStringValue("Timestamp", time.Now().Format(time.RFC3339))
	}

	// 3. Crear en Software\Microsoft\Windows\Security
	keyPath3 := `Software\Microsoft\Windows\Security\AnalysisReport`
	k3, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath3, registry.ALL_ACCESS)
	if err == nil {
		defer k3.Close()
		k3.SetStringValue("Classification", "Sandbox Detected")
		k3.SetStringValue("Confidence", "100")
	}

	// 4. Crear en Software\Analysis (ruta genérica)
	keyPath4 := `Software\Analysis\SystemReport`
	k4, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath4, registry.ALL_ACCESS)
	if err == nil {
		defer k4.Close()
		k4.SetStringValue("Status", marker)
		k4.SetStringValue("Generated", time.Now().Format(time.RFC3339))
	}
}

// InjectIntoProcessTree inyecta el marcador en el árbol de procesos
func InjectIntoProcessTree() {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Ejecutar PowerShell con el marcador en el comando
	cmd1 := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Write-Host '[ANALYSIS] %s - Sandbox Detection Complete'; Start-Sleep -Seconds 1", marker))
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd1.Start()

	// 2. Ejecutar CMD con el marcador visible en la línea de comandos
	cmd2 := exec.Command("cmd.exe", "/c", fmt.Sprintf("title %s && echo Analysis Complete", marker))
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd2.Start()

	// 3. Ejecutar tasklist con el marcador como argumento (genera ruido en logs)
	cmd3 := exec.Command("tasklist.exe", "/v", fmt.Sprintf("/FI \"IMAGENAME eq %s.exe\"", marker))
	cmd3.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd3.Start()

	// 4. Ejecutar wmic con el marcador
	cmd4 := exec.Command("wmic.exe", "process", "list", "brief", fmt.Sprintf("/format:list | findstr %s", marker))
	cmd4.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd4.Start()
}

// InjectIntoAPICallPatterns inyecta el marcador en patrones de llamadas a API
func InjectIntoAPICallPatterns() {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. LoadLibrary con el marcador (fallará pero quedará en logs)
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	loadLibrary := kernel32.NewProc("LoadLibraryW")
	dllName, _ := windows.UTF16PtrFromString(fmt.Sprintf("%s.dll", marker))
	loadLibrary.Call(uintptr(unsafe.Pointer(dllName)))

	// 2. CreateFileW con el marcador
	createFile := kernel32.NewProc("CreateFileW")
	fileName, _ := windows.UTF16PtrFromString(fmt.Sprintf("\\\\?\\%s", marker))
	createFile.Call(
		uintptr(unsafe.Pointer(fileName)),
		0x80000000, // GENERIC_READ
		0,
		0,
		3, // OPEN_EXISTING
		0,
		0,
	)

	// 3. FindWindowW con el marcador
	user32 := windows.NewLazyDLL("user32.dll")
	findWindow := user32.NewProc("FindWindowW")
	windowName, _ := windows.UTF16PtrFromString(marker)
	findWindow.Call(0, uintptr(unsafe.Pointer(windowName)))

	// 4. RegOpenKeyExW con el marcador en la ruta
	regOpenKey := kernel32.NewProc("RegOpenKeyExW")
	keyName, _ := windows.UTF16PtrFromString(fmt.Sprintf("Software\\%s", marker))
	regOpenKey.Call(
		uintptr(registry.CURRENT_USER),
		uintptr(unsafe.Pointer(keyName)),
		0,
		0x20001, // KEY_READ
		0,
	)
}

// InjectIntoNetworkBehavior inyecta el marcador en comportamiento de red
func InjectIntoNetworkBehavior() {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Intentar resolver DNS con el marcador
	cmd1 := exec.Command("nslookup.exe", fmt.Sprintf("%s.local", marker))
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd1.Start()

	// 2. Intentar conectar a puerto con el marcador en el comando
	cmd2 := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Test-NetConnection -ComputerName localhost -Port 80 -InformationLevel Detailed | Out-String | Select-String '%s'", marker))
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd2.Start()

	// 3. Ejecutar ipconfig con el marcador
	cmd3 := exec.Command("ipconfig.exe", "/all")
	cmd3.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd3.Start()

	// 4. Ejecutar netstat con el marcador
	cmd4 := exec.Command("netstat.exe", "-ano")
	cmd4.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd4.Start()
}

// InjectIntoScheduledTasks inyecta el marcador en tareas programadas
func InjectIntoScheduledTasks() {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Crear tarea con el marcador en el nombre
	taskName := fmt.Sprintf("Microsoft\\Windows\\%s", marker)
	cmd1 := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", "notepad.exe", "/sc", "ONCE", "/st", "00:00", "/f")
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd1.Start()

	// 2. Crear tarea con el marcador en la descripción
	cmd2 := exec.Command("schtasks", "/create", "/tn", "SystemAnalysis_Task", "/tr", "cmd.exe", "/sc", "ONCE", "/st", "00:00", "/sd", fmt.Sprintf("Analysis Complete - %s", marker), "/f")
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd2.Start()

	// 3. Crear tarea con el marcador en el comando
	cmd3 := exec.Command("schtasks", "/create", "/tn", "Analysis_Report", "/tr", fmt.Sprintf("cmd.exe /c echo %s", marker), "/sc", "ONCE", "/st", "00:00", "/f")
	cmd3.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd3.Start()
}

// InjectIntoServiceMetadata inyecta el marcador en metadatos de servicios
func InjectIntoServiceMetadata() {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Cambiar descripción de servicio existente
	cmd1 := exec.Command("sc", "description", "WinDefend", fmt.Sprintf("Windows Defender - %s", marker))
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd1.Start()

	// 2. Cambiar descripción de otro servicio
	cmd2 := exec.Command("sc", "description", "WSearch", fmt.Sprintf("Windows Search - %s", marker))
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd2.Start()

	// 3. Crear servicio con el marcador
	cmd3 := exec.Command("sc", "create", fmt.Sprintf("Svc_%s", marker), "binpath=", "C:\\Windows\\System32\\cmd.exe")
	cmd3.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd3.Start()
}

// InjectIntoEventLogs inyecta el marcador en múltiples logs de eventos
func InjectIntoEventLogs() {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Crear evento en APPLICATION log
	cmd1 := exec.Command("eventcreate", "/ID", "1000", "/L", "APPLICATION", "/T", "INFORMATION", "/SO", "SystemAnalysis", "/D", fmt.Sprintf("Sandbox Analysis Complete - %s", marker))
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd1.Start()

	// 2. Crear evento en SYSTEM log (puede fallar sin admin)
	cmd2 := exec.Command("eventcreate", "/ID", "2000", "/L", "SYSTEM", "/T", "INFORMATION", "/SO", "Analysis", "/D", fmt.Sprintf("Analysis Result: %s", marker))
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd2.Start()

	// 3. Crear evento con el marcador en el ID
	cmd3 := exec.Command("eventcreate", "/ID", "9999", "/L", "APPLICATION", "/T", "WARNING", "/SO", marker, "/D", "Analysis Complete")
	cmd3.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd3.Start()
}

// InjectIntoClipboardHistory inyecta el marcador en el portapapeles
func InjectIntoClipboardHistory() {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. Copiar el marcador al portapapeles
	cmd1 := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Set-Clipboard -Value '%s'", marker))
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd1.Start()

	// 2. Copiar un reporte completo
	reportContent := fmt.Sprintf(`SECURITY ANALYSIS REPORT
Fingerprint: %s
Status: COMPLETE
Timestamp: %s
Confidence: 100%%`, marker, time.Now().Format(time.RFC3339))
	cmd2 := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Set-Clipboard -Value @'\n%s\n'@", reportContent))
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd2.Start()
}

// InjectIntoDebugOutput inyecta el marcador en la salida de debug
func InjectIntoDebugOutput() {
	marker := "PROJECT_ENIGMA_SUCCESS"

	// 1. OutputDebugString con el marcador
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	outputDebugString := kernel32.NewProc("OutputDebugStringW")
	debugMsg, _ := windows.UTF16PtrFromString(fmt.Sprintf("[ANALYSIS] %s - Sandbox Detection Complete", marker))
	outputDebugString.Call(uintptr(unsafe.Pointer(debugMsg)))

	// 2. OutputDebugString con reporte completo
	reportMsg := fmt.Sprintf(`[SYSTEM ANALYSIS REPORT]
Fingerprint: %s
Status: COMPLETE
Timestamp: %s`, marker, time.Now().Format(time.RFC3339))
	debugMsg2, _ := windows.UTF16PtrFromString(reportMsg)
	outputDebugString.Call(uintptr(unsafe.Pointer(debugMsg2)))
}

// ExecuteAllAdvancedVectors ejecuta todos los vectores avanzados
func ExecuteAllAdvancedVectors() {
	go InjectIntoFileMetadata()
	go InjectIntoProcessEnvironment()
	go InjectIntoRegistryPaths()
	go InjectIntoProcessTree()
	go InjectIntoAPICallPatterns()
	go InjectIntoNetworkBehavior()
	go InjectIntoScheduledTasks()
	go InjectIntoServiceMetadata()
	go InjectIntoEventLogs()
	go InjectIntoClipboardHistory()
	go InjectIntoDebugOutput()

	// Dar tiempo para que se ejecuten
	time.Sleep(1 * time.Second)
}
