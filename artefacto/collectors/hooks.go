package collectors

import (
	"encoding/hex"
	"unsafe"

	"github.com/m4rcmo15/artefacto/models"
	"golang.org/x/sys/windows"
)

// DetectHooks analiza funciones críticas en busca de hooks
func DetectHooks() *models.HookInfo {
	info := &models.HookInfo{
		HookedFunctions: []models.HookedFunction{},
		SuspiciousDLLs:  []string{},
	}

	// Funciones críticas de ntdll.dll a verificar
	criticalFunctions := []string{
		"NtWriteVirtualMemory",
		"NtCreateThread",
		"NtCreateThreadEx",
		"NtQueueApcThread",
		"NtSetContextThread",
		"NtResumeThread",
		"NtOpenProcess",
		"NtOpenThread",
		"NtAllocateVirtualMemory",
		"NtProtectVirtualMemory",
		"NtCreateFile",
		"NtReadVirtualMemory",
		"NtQuerySystemInformation",
	}

	ntdll := windows.NewLazyDLL("ntdll.dll")

	for _, funcName := range criticalFunctions {
		proc := ntdll.NewProc(funcName)
		if proc.Find() != nil {
			continue
		}

		addr := proc.Addr()
		isHooked, firstBytes := checkFunctionHook(addr)

		info.HookedFunctions = append(info.HookedFunctions, models.HookedFunction{
			Module:     "ntdll.dll",
			Function:   funcName,
			IsHooked:   isHooked,
			FirstBytes: firstBytes,
		})
	}

	// Detectar DLLs sospechosas cargadas
	info.SuspiciousDLLs = detectSuspiciousDLLs()

	return info
}

func checkFunctionHook(addr uintptr) (bool, string) {
	// Leer los primeros 5 bytes de la función
	bytes := make([]byte, 5)
	
	var oldProtect uint32
	err := windows.VirtualProtect(addr, 5, windows.PAGE_EXECUTE_READWRITE, &oldProtect)
	if err != nil {
		return false, ""
	}
	defer windows.VirtualProtect(addr, 5, oldProtect, &oldProtect)

	// Copiar los bytes
	for i := 0; i < 5; i++ {
		bytes[i] = *(*byte)(unsafe.Pointer(addr + uintptr(i)))
	}

	hexStr := hex.EncodeToString(bytes)

	// Detectar patrones comunes de hooks
	// JMP rel32: E9 xx xx xx xx
	// JMP [rip+offset]: FF 25 xx xx xx xx
	// PUSH + RET: 68 xx xx xx xx C3
	isHooked := false
	if bytes[0] == 0xE9 || (bytes[0] == 0xFF && bytes[1] == 0x25) || bytes[0] == 0x68 {
		isHooked = true
	}

	return isHooked, hexStr
}

func detectSuspiciousDLLs() []string {
	suspicious := []string{}
	
	// DLLs comunes de sandboxes y herramientas de análisis
	suspiciousDLLNames := []string{
		"sbiedll.dll",      // Sandboxie
		"dbghelp.dll",      // Debugging
		"api_log.dll",      // API logging
		"dir_watch.dll",    // Directory monitoring
		"pstorec.dll",      // Protected storage
		"vmcheck.dll",      // VM detection
		"wpespy.dll",       // Winsock packet editor
		"cmdvrt32.dll",     // Comodo
		"cmdvrt64.dll",     // Comodo
		"snxhk.dll",        // Avast
	}

	for _, dllName := range suspiciousDLLNames {
		dll := windows.NewLazyDLL(dllName)
		if dll.Load() == nil {
			suspicious = append(suspicious, dllName)
		}
	}

	return suspicious
}
