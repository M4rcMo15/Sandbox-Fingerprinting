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
		switch payload.Vector {
		case "hostname":
			// El hostname ya se modifica en main.go
			fmt.Printf("[XSS] Payload %s inyectado en hostname\n", payload.ID[:8])

		case "filename":
			injectIntoFilename(payload)

		case "process":
			injectIntoProcess(payload)

		case "registry":
			injectIntoRegistry(payload)

		case "window":
			injectIntoWindow(payload)

		case "cmdline":
			injectIntoCmdLine(payload)
		}
	}
}

// injectIntoFilename crea archivos con nombres que contienen XSS
func injectIntoFilename(payload XSSPayload) {
	tempDir := os.TempDir()
	filename := filepath.Join(tempDir, payload.Content)

	// Crear archivo con contenido inocuo
	content := []byte("This is a test file for XSS audit")
	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		fmt.Printf("[XSS] Error creando archivo: %v\n", err)
		return
	}

	fmt.Printf("[XSS] Payload %s inyectado en filename: %s\n", payload.ID[:8], filename)
}

// injectIntoProcess ejecuta un proceso con nombre/argumento XSS
func injectIntoProcess(payload XSSPayload) {
	// Ejecutar cmd.exe con argumento que contiene el payload
	cmd := exec.Command("cmd.exe", "/c", "echo", payload.Content)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("[XSS] Error ejecutando proceso: %v\n", err)
		return
	}

	// No esperamos a que termine, solo lo lanzamos
	fmt.Printf("[XSS] Payload %s inyectado en proceso (PID: %d)\n", payload.ID[:8], cmd.Process.Pid)
}

// injectIntoRegistry crea una clave de registro con XSS
func injectIntoRegistry(payload XSSPayload) {
	// Crear clave en HKEY_CURRENT_USER\Software\XSSAudit
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software`, registry.CREATE_SUB_KEY)
	if err != nil {
		fmt.Printf("[XSS] Error abriendo registro: %v\n", err)
		return
	}
	defer key.Close()

	// Crear subclave con nombre XSS
	subKey, _, err := registry.CreateKey(key, "XSSAudit_"+payload.ID[:8], registry.SET_VALUE)
	if err != nil {
		fmt.Printf("[XSS] Error creando subclave: %v\n", err)
		return
	}
	defer subKey.Close()

	// Establecer valor con payload
	err = subKey.SetStringValue("TestValue", payload.Content)
	if err != nil {
		fmt.Printf("[XSS] Error estableciendo valor: %v\n", err)
		return
	}

	fmt.Printf("[XSS] Payload %s inyectado en registro\n", payload.ID[:8])
}

// injectIntoWindow crea una ventana con título XSS
func injectIntoWindow(payload XSSPayload) {
	// Crear ventana invisible con título XSS
	user32 := windows.NewLazyDLL("user32.dll")
	createWindowEx := user32.NewProc("CreateWindowExW")

	className, _ := windows.UTF16PtrFromString("STATIC")
	windowName, _ := windows.UTF16PtrFromString(payload.Content)

	hwnd, _, _ := createWindowEx.Call(
		0,                          // dwExStyle
		uintptr(unsafe.Pointer(className)),   // lpClassName
		uintptr(unsafe.Pointer(windowName)),  // lpWindowName
		0,                          // dwStyle
		0, 0, 0, 0,                // x, y, width, height
		0,                          // hWndParent
		0,                          // hMenu
		0,                          // hInstance
		0,                          // lpParam
	)

	if hwnd != 0 {
		fmt.Printf("[XSS] Payload %s inyectado en ventana (HWND: %d)\n", payload.ID[:8], hwnd)
	}
}

// injectIntoCmdLine ejecuta comando con argumento XSS
func injectIntoCmdLine(payload XSSPayload) {
	// Ejecutar notepad con argumento XSS
	cmd := exec.Command("notepad.exe", payload.Content)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("[XSS] Error ejecutando comando: %v\n", err)
		return
	}

	fmt.Printf("[XSS] Payload %s inyectado en command line\n", payload.ID[:8])
}
