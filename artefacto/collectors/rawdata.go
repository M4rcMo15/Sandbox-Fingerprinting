package collectors

import (
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/m4rcmo15/artefacto/models"
	"golang.org/x/sys/windows"
)

// CollectRawData recopila datos en bruto sin procesamiento
// El análisis se hará en el servidor
func CollectRawData() *models.RawData {
	data := &models.RawData{
		VMFiles:           []string{},
		RegistryKeys:      []models.RegistryKey{},
		SecurityProcesses: []string{},
		Drivers:           []string{},
		DiskInfo:          models.DiskInfo{},
		CPUInfo:           models.CPUInfo{},
		MouseHistory:      []models.MousePoint{},
	}

	// Recopilar archivos relacionados con VMs (sin determinar si es VM)
	data.VMFiles = collectVMFiles()

	// Recopilar claves de registro (sin analizar)
	data.RegistryKeys = collectRegistryKeys()

	// Recopilar procesos de seguridad (sin identificar productos)
	data.SecurityProcesses = collectSecurityProcesses()

	// Recopilar drivers (sin analizar)
	data.Drivers = collectDrivers()

	// Información del disco (datos en bruto)
	data.DiskInfo = collectDiskInfo()

	// Información de CPU (datos en bruto)
	data.CPUInfo = collectCPUInfo()

	// Temperatura de CPU (si está disponible)
	data.CPUInfo.Temperature = getCPUTemperatureRaw()

	// Contador de ventanas
	data.WindowCount = countWindows()

	// A. Timing Attacks (Diferencia entre reloj OS y CPU)
	data.TimingDiscrepancy = checkTimingDiscrepancy()

	// B. Human Interaction (Historial de mouse)
	data.MouseHistory = collectMouseHistory()

	// C. Hardware Profundo (CPUID y MAC OUI)
	data.CPUIDHypervisor = checkCPUIDHypervisor()
	data.MACOUI = collectMACOUI()

	// D. Clipboard
	data.ClipboardPreview = collectClipboard()

	return data
}

func collectVMFiles() []string {
	files := []string{}

	// Lista de archivos a verificar (sin determinar si es VM)
	checkPaths := []string{
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
		"C:\\Windows\\System32\\drivers\\prleth.sys",
		"C:\\Windows\\System32\\drivers\\prlfs.sys",
		"C:\\Windows\\System32\\drivers\\prlmouse.sys",
		"C:\\Windows\\System32\\drivers\\prlvideo.sys",
		"C:\\Windows\\System32\\drivers\\prl_pv32.sys",
		"C:\\Windows\\System32\\drivers\\prl_paravirt_32.sys",
	}

	for _, path := range checkPaths {
		if _, err := os.Stat(path); err == nil {
			files = append(files, path)
		}
	}

	return files
}

func collectRegistryKeys() []models.RegistryKey {
	keys := []models.RegistryKey{}

	// Lista de claves de registro a recopilar (sin analizar)
	registryPaths := []struct {
		root windows.Handle
		path string
		name string
	}{
		// VirtualBox
		{windows.HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\VBoxGuest", "VBoxGuest"},
		{windows.HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\VBoxMouse", "VBoxMouse"},
		{windows.HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\VBoxService", "VBoxService"},
		{windows.HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\VBoxSF", "VBoxSF"},

		// VMware
		{windows.HKEY_LOCAL_MACHINE, "SOFTWARE\\VMware, Inc.\\VMware Tools", "VMware Tools"},
		{windows.HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\vmmouse", "vmmouse"},
		{windows.HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\vmhgfs", "vmhgfs"},

		// Hyper-V
		{windows.HKEY_LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Hyper-V", "Hyper-V"},
		{windows.HKEY_LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Virtual Machine\\Guest", "VM Guest"},

		// QEMU
		{windows.HKEY_LOCAL_MACHINE, "HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 0\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", "SCSI"},

		// Parallels
		{windows.HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\prl_pv32", "Parallels"},

		// Xen
		{windows.HKEY_LOCAL_MACHINE, "HARDWARE\\ACPI\\DSDT\\Xen", "Xen"},
	}

	for _, reg := range registryPaths {
		var key windows.Handle
		err := windows.RegOpenKeyEx(reg.root, windows.StringToUTF16Ptr(reg.path), 0, windows.KEY_READ, &key)
		if err == nil {
			// Leer valores de la clave
			values := readRegistryValues(key)
			keys = append(keys, models.RegistryKey{
				Path:   reg.path,
				Name:   reg.name,
				Exists: true,
				Values: values,
			})
			windows.RegCloseKey(key)
		} else {
			keys = append(keys, models.RegistryKey{
				Path:   reg.path,
				Name:   reg.name,
				Exists: false,
			})
		}
	}

	return keys
}

func readRegistryValues(key windows.Handle) map[string]string {
	values := make(map[string]string)

	// Lista de valores comunes a leer
	commonValues := []string{
		"DisplayName",
		"Description",
		"ImagePath",
		"Start",
		"Type",
		"ErrorControl",
	}

	for _, valueName := range commonValues {
		var bufLen uint32 = 1024
		buf := make([]uint16, bufLen)

		err := windows.RegQueryValueEx(
			key,
			windows.StringToUTF16Ptr(valueName),
			nil,
			nil,
			(*byte)(unsafe.Pointer(&buf[0])),
			&bufLen,
		)

		if err == nil {
			values[valueName] = windows.UTF16ToString(buf)
		}
	}

	return values
}

func collectSecurityProcesses() []string {
	processes := []string{}

	// Lista de procesos de seguridad conocidos (sin identificar productos)
	securityKeywords := []string{
		// EDR/AV general
		"defender", "antivirus", "security", "edr", "endpoint",
		// Productos específicos
		"crowdstrike", "falcon", "sentinel", "carbon", "cylance",
		"symantec", "mcafee", "kaspersky", "trend", "eset",
		"palo", "traps", "fireeye", "xagt", "sophos",
		"avast", "avg", "bitdefender", "norton", "webroot",
		// Procesos comunes
		"mssense", "msmpeng", "nissrv", "securityhealth",
		"csfalcon", "sentinelagent", "cb.exe", "cylancesvc",
	}

	// Obtener todos los procesos
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return processes
	}
	defer windows.CloseHandle(snapshot)

	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	err = windows.Process32First(snapshot, &procEntry)
	if err != nil {
		return processes
	}

	for {
		name := windows.UTF16ToString(procEntry.ExeFile[:])
		nameLower := strings.ToLower(name)

		// Verificar si contiene palabras clave de seguridad
		for _, keyword := range securityKeywords {
			if strings.Contains(nameLower, keyword) {
				processes = append(processes, name)
				break
			}
		}

		err = windows.Process32Next(snapshot, &procEntry)
		if err != nil {
			break
		}
	}

	return processes
}

func collectDrivers() []string {
	drivers := []string{}

	// Leer todos los drivers de System32\drivers
	driverPath := "C:\\Windows\\System32\\drivers"

	entries, err := os.ReadDir(driverPath)
	if err != nil {
		return drivers
	}

	// Recopilar todos los .sys (sin filtrar)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".sys") {
			drivers = append(drivers, entry.Name())
		}
	}

	return drivers
}

func collectDiskInfo() models.DiskInfo {
	info := models.DiskInfo{}

	// Leer identificador del disco desde el registro
	var key windows.Handle
	err := windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`HARDWARE\DEVICEMAP\Scsi\Scsi Port 0\Scsi Bus 0\Target Id 0\Logical Unit Id 0`),
		0,
		windows.KEY_READ,
		&key,
	)
	if err == nil {
		defer windows.RegCloseKey(key)

		// Leer valores
		info.Identifier = readRegistryStringRaw(key, "Identifier")
		info.SerialNumber = readRegistryStringRaw(key, "SerialNumber")
	}

	// Obtener tamaño del disco
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes int64

	ret, _, _ := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr("C:\\"))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
	)

	if ret != 0 {
		info.TotalBytes = totalNumberOfBytes
		info.FreeBytes = totalNumberOfFreeBytes
	}

	return info
}

func collectCPUInfo() models.CPUInfo {
	info := models.CPUInfo{}

	// Leer información de CPU desde el registro
	var key windows.Handle
	err := windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`HARDWARE\DESCRIPTION\System\CentralProcessor\0`),
		0,
		windows.KEY_READ,
		&key,
	)
	if err == nil {
		defer windows.RegCloseKey(key)

		info.ProcessorName = readRegistryStringRaw(key, "ProcessorNameString")
		info.Vendor = readRegistryStringRaw(key, "VendorIdentifier")
		info.Identifier = readRegistryStringRaw(key, "Identifier")
	}

	return info
}

func readRegistryStringRaw(key windows.Handle, valueName string) string {
	var bufLen uint32 = 256
	buf := make([]uint16, bufLen)

	err := windows.RegQueryValueEx(
		key,
		windows.StringToUTF16Ptr(valueName),
		nil,
		nil,
		(*byte)(unsafe.Pointer(&buf[0])),
		&bufLen,
	)
	if err != nil {
		return ""
	}

	return windows.UTF16ToString(buf)
}

func getCPUTemperatureRaw() float64 {
	// Intento 1: MSAcpi_ThermalZoneTemperature (Deci-Kelvin)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
		"try { (Get-WmiObject -Namespace root/wmi -Class MSAcpi_ThermalZoneTemperature -ErrorAction Stop | Select-Object -ExpandProperty CurrentTemperature | Select-Object -First 1) } catch { 0 }")

	// Ocultar ventana de consola para evitar parpadeos
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}

	out, err := cmd.Output()
	if err == nil {
		valStr := strings.TrimSpace(string(out))
		val, err := strconv.ParseFloat(valStr, 64)
		if err == nil && val > 0 {
			// Convertir deci-Kelvin a Celsius: C = (dK / 10) - 273.15
			return (val / 10.0) - 273.15
		}
	}

	// Intento 2: Win32_PerfFormattedData_Counters_ThermalZoneInformation (Kelvin)
	// Este suele funcionar en máquinas donde MSAcpi falla
	cmd2 := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
		"try { (Get-WmiObject -Class Win32_PerfFormattedData_Counters_ThermalZoneInformation -ErrorAction Stop | Select-Object -ExpandProperty Temperature | Select-Object -First 1) } catch { 0 }")
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}

	out2, err2 := cmd2.Output()
	if err2 == nil {
		valStr2 := strings.TrimSpace(string(out2))
		val2, err2 := strconv.ParseFloat(valStr2, 64)
		if err2 == nil && val2 > 0 {
			// Convertir Kelvin a Celsius: C = K - 273.15
			return val2 - 273.15
		}
	}

	return 0.0
}

func checkTimingDiscrepancy() float64 {
	// Medir discrepancia entre reloj del OS y ciclos de CPU (QPC)
	// Las sandboxes que aceleran el tiempo suelen desincronizar estos relojes
	startOS := time.Now()

	var startQPC, endQPC, freq int64
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	qpc := kernel32.NewProc("QueryPerformanceCounter")
	qpf := kernel32.NewProc("QueryPerformanceFrequency")

	qpf.Call(uintptr(unsafe.Pointer(&freq)))
	qpc.Call(uintptr(unsafe.Pointer(&startQPC)))

	// Dormir 5s (si la sandbox acelera el sleep, el QPC revelará que pasó menos tiempo real)
	time.Sleep(200 * time.Millisecond) // Reducido de 2s a 200ms para velocidad

	endOS := time.Now()
	qpc.Call(uintptr(unsafe.Pointer(&endQPC)))

	diffOS := endOS.Sub(startOS).Seconds()
	diffQPC := float64(endQPC-startQPC) / float64(freq)

	// Retornar la diferencia absoluta
	diff := diffOS - diffQPC
	if diff < 0 {
		return -diff
	}
	return diff
}

func collectMouseHistory() []models.MousePoint {
	history := []models.MousePoint{}

	// Recopilar 50 muestras en ~500ms para analizar entropía de movimiento
	for i := 0; i < 10; i++ { // Reducido de 50 a 10 muestras para velocidad
		pos := getMousePosition() // Reutiliza función existente en el paquete collectors
		history = append(history, models.MousePoint{
			X:    pos.X,
			Y:    pos.Y,
			Time: time.Now().UnixMilli(),
		})
		time.Sleep(10 * time.Millisecond)
	}

	return history
}

func checkCPUIDHypervisor() bool {
	// Verificar bit de hipervisor usando WMI (Win32_ComputerSystem.HypervisorPresent)
	// Es más robusto en Go puro que intentar inyectar ensamblador inline
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
		"(Get-WmiObject Win32_ComputerSystem).HypervisorPresent")
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "True"
}

func collectMACOUI() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		// Ignorar loopback y interfaces caídas
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// Obtener MAC y devolver los primeros 3 bytes (OUI)
		mac := iface.HardwareAddr
		if len(mac) >= 3 {
			return hex.EncodeToString(mac[:3])
		}
	}

	return ""
}

func collectClipboard() string {
	user32 := windows.NewLazyDLL("user32.dll")
	kernel32 := windows.NewLazyDLL("kernel32.dll")

	openClipboard := user32.NewProc("OpenClipboard")
	closeClipboard := user32.NewProc("CloseClipboard")
	getClipboardData := user32.NewProc("GetClipboardData")
	globalLock := kernel32.NewProc("GlobalLock")
	globalUnlock := kernel32.NewProc("GlobalUnlock")

	// Intentar abrir clipboard
	if ret, _, _ := openClipboard.Call(0); ret == 0 {
		return ""
	}
	defer closeClipboard.Call()

	// CF_UNICODETEXT = 13 (Mejor compatibilidad que CF_TEXT)
	hData, _, _ := getClipboardData.Call(13)
	if hData == 0 {
		return ""
	}

	ptr, _, _ := globalLock.Call(hData)
	if ptr == 0 {
		return ""
	}
	defer globalUnlock.Call(hData)

	// Leer string UTF-16
	return windows.UTF16PtrToString((*uint16)(unsafe.Pointer(ptr)))
}

func countWindows() int {
	count := 0

	user32 := windows.NewLazyDLL("user32.dll")
	enumWindows := user32.NewProc("EnumWindows")

	enumWindowsProc := windows.NewCallback(func(hwnd windows.Handle, lParam uintptr) uintptr {
		if isWindowVisibleRaw(hwnd) {
			count++
		}
		return 1
	})

	enumWindows.Call(enumWindowsProc, 0)
	return count
}

func isWindowVisibleRaw(hwnd windows.Handle) bool {
	ret, _, _ := windows.NewLazySystemDLL("user32.dll").
		NewProc("IsWindowVisible").
		Call(uintptr(hwnd))
	return ret != 0
}

// GetCPUCount retorna el número de CPUs
func GetCPUCount() int {
	return runtime.NumCPU()
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

// GetPublicIP obtiene la IP pública del sistema
func GetPublicIP() string {
	client := &http.Client{Timeout: 5 * time.Second}

	services := []string{
		"https://api.ipify.org",
		"https://icanhazip.com",
		"https://ifconfig.me/ip",
	}

	for _, service := range services {
		resp, err := client.Get(service)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		return string(body)
	}

	return ""
}
