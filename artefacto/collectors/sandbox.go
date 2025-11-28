package collectors

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/m4rcmo15/artefacto/models"
	"golang.org/x/sys/windows"
)

// CheckSandbox detecta si el sistema está corriendo en un entorno virtualizado
func CheckSandbox() *models.SandboxInfo {
	info := &models.SandboxInfo{
		VMIndicators:       []string{},
		RegistryIndicators: []string{},
		DiskIndicators:     []string{},
	}

	// Verificar archivos de VM
	checkVMFiles(info)

	// Verificar registro de Windows
	checkRegistry(info)

	// Verificar nombres de disco
	checkDiskNames(info)

	// Verificar temperatura de CPU
	info.CPUTemperature = getCPUTemperature()

	// Contar ventanas abiertas
	info.WindowCount = getWindowCount()

	// Verificar privilegios de depuración
	info.HasDebugPrivilege = hasDebugPrivilege()

	// Determinar si es VM
	info.IsVM = len(info.VMIndicators) > 0 || 
		len(info.RegistryIndicators) > 0 || 
		len(info.DiskIndicators) > 0 ||
		info.CPUTemperature == 0.0 ||
		info.WindowCount < 10

	return info
}

func checkVMFiles(info *models.SandboxInfo) {
	vmFiles := []string{
		"C:\\Windows\\System32\\drivers\\VBoxMouse.sys",
		"C:\\Windows\\System32\\drivers\\VBoxGuest.sys",
		"C:\\Windows\\System32\\drivers\\VBoxSF.sys",
		"C:\\Windows\\System32\\drivers\\VBoxVideo.sys",
		"C:\\Windows\\System32\\vboxdisp.dll",
		"C:\\Windows\\System32\\vboxhook.dll",
		"C:\\Windows\\System32\\drivers\\vmmouse.sys",
		"C:\\Windows\\System32\\drivers\\vmhgfs.sys",
		"C:\\Windows\\System32\\drivers\\vmmemctl.sys",
		"C:\\Program Files\\VMware\\VMware Tools\\",
		"C:\\Program Files\\Oracle\\VirtualBox Guest Additions\\",
	}

	for _, file := range vmFiles {
		if _, err := os.Stat(file); err == nil {
			info.VMIndicators = append(info.VMIndicators, file)
		}
	}
}

func checkRegistry(info *models.SandboxInfo) {
	// Claves de registro comunes en VMs
	registryKeys := []struct {
		root windows.Handle
		path string
		key  string
	}{
		{windows.HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\VBoxGuest", ""},
		{windows.HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\VBoxMouse", ""},
		{windows.HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\VBoxService", ""},
		{windows.HKEY_LOCAL_MACHINE, "SOFTWARE\\VMware, Inc.\\VMware Tools", ""},
		{windows.HKEY_LOCAL_MACHINE, "HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 0\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", "Identifier"},
	}

	for _, reg := range registryKeys {
		var key windows.Handle
		err := windows.RegOpenKeyEx(reg.root, windows.StringToUTF16Ptr(reg.path), 0, windows.KEY_READ, &key)
		if err == nil {
			info.RegistryIndicators = append(info.RegistryIndicators, reg.path)
			windows.RegCloseKey(key)
		}
	}
}

func checkDiskNames(info *models.SandboxInfo) {
	// Nombres de disco sospechosos
	suspiciousNames := []string{"VBOX", "VMware", "QEMU", "Virtual", "DADY HARDDISK"}

	// Leer el identificador del disco desde el registro
	var key windows.Handle
	err := windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`HARDWARE\DEVICEMAP\Scsi\Scsi Port 0\Scsi Bus 0\Target Id 0\Logical Unit Id 0`),
		0,
		windows.KEY_READ,
		&key,
	)
	if err != nil {
		return
	}
	defer windows.RegCloseKey(key)

	// Leer el valor "Identifier"
	var bufLen uint32 = 256
	buf := make([]uint16, bufLen)

	err = windows.RegQueryValueEx(
		key,
		windows.StringToUTF16Ptr("Identifier"),
		nil,
		nil,
		(*byte)(unsafe.Pointer(&buf[0])),
		&bufLen,
	)
	if err != nil {
		return
	}

	identifier := windows.UTF16ToString(buf)

	// Verificar si contiene nombres sospechosos
	identifierUpper := strings.ToUpper(identifier)
	for _, suspicious := range suspiciousNames {
		if strings.Contains(identifierUpper, strings.ToUpper(suspicious)) {
			info.DiskIndicators = append(info.DiskIndicators, identifier)
			break
		}
	}
}

func getCPUTemperature() float64 {
	// Los sandboxes a menudo no reportan temperatura
	// Intentar leer temperatura del registro (MSR - Model Specific Register)
	// Nota: Esto requiere privilegios elevados y drivers específicos
	// En la mayoría de casos retornará 0.0 en VMs
	
	// Verificar si existe el registro de temperatura
	var key windows.Handle
	err := windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`HARDWARE\DESCRIPTION\System\CentralProcessor\0`),
		0,
		windows.KEY_READ,
		&key,
	)
	if err != nil {
		return 0.0
	}
	defer windows.RegCloseKey(key)

	// En sistemas reales, la temperatura estaría disponible
	// En VMs, típicamente no hay sensores de temperatura
	// Por ahora retornamos 0.0 como indicador de VM
	return 0.0
}

func getWindowCount() int {
	// Cuenta las ventanas visibles del sistema
	// Los sandboxes suelen tener pocas ventanas abiertas
	count := 0
	
	user32 := windows.NewLazyDLL("user32.dll")
	enumWindows := user32.NewProc("EnumWindows")
	
	enumWindowsProc := windows.NewCallback(func(hwnd windows.Handle, lParam uintptr) uintptr {
		if isWindowVisible(hwnd) {
			count++
		}
		return 1 // Continuar enumeración
	})

	enumWindows.Call(enumWindowsProc, 0)
	return count
}

func isWindowVisible(hwnd windows.Handle) bool {
	ret, _, _ := windows.NewLazySystemDLL("user32.dll").
		NewProc("IsWindowVisible").
		Call(uintptr(hwnd))
	return ret != 0
}

func hasDebugPrivilege() bool {
	var token windows.Token
	proc, _ := windows.GetCurrentProcess()
	
	err := windows.OpenProcessToken(proc, windows.TOKEN_QUERY, &token)
	if err != nil {
		return false
	}
	defer token.Close()

	// Intentar habilitar SeDebugPrivilege
	var luid windows.LUID
	err = windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr("SeDebugPrivilege"), &luid)
	if err != nil {
		return false
	}

	// Verificar si el privilegio está disponible
	privileges := windows.Tokenprivileges{
		PrivilegeCount: 1,
		Privileges: [1]windows.LUIDAndAttributes{
			{
				Luid:       luid,
				Attributes: windows.SE_PRIVILEGE_ENABLED,
			},
		},
	}

	err = windows.AdjustTokenPrivileges(token, false, &privileges, 0, nil, nil)
	return err == nil
}

// GetMemoryInfo obtiene información de memoria del sistema
func GetMemoryInfo() (totalRAM uint64, err error) {
	type MEMORYSTATUSEX struct {
		DwLength                uint32
		DwMemoryLoad            uint32
		UllTotalPhys            uint64
		UllAvailPhys            uint64
		UllTotalPageFile        uint64
		UllAvailPageFile        uint64
		UllTotalVirtual         uint64
		UllAvailVirtual         uint64
		UllAvailExtendedVirtual uint64
	}

	kernel32 := windows.NewLazyDLL("kernel32.dll")
	globalMemoryStatusEx := kernel32.NewProc("GlobalMemoryStatusEx")

	memStatus := MEMORYSTATUSEX{}
	memStatus.DwLength = uint32(unsafe.Sizeof(memStatus))

	ret, _, err := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))
	if ret == 0 {
		return 0, fmt.Errorf("GlobalMemoryStatusEx failed: %v", err)
	}

	// Convertir bytes a MB
	totalRAM = memStatus.UllTotalPhys / (1024 * 1024)
	return totalRAM, nil
}

// GetDiskInfo obtiene información del disco
func GetDiskInfo() (totalBytes int64, err error) {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes int64

	ret, _, err := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr("C:\\"))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
	)

	if ret == 0 {
		return 0, fmt.Errorf("GetDiskFreeSpaceExW failed: %v", err)
	}

	return totalNumberOfBytes, nil
}

// GetCPUCount retorna el número de CPUs
func GetCPUCount() int {
	return runtime.NumCPU()
}
