package collectors

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/m4rcmo15/artefacto/models"
	"golang.org/x/sys/windows"
)

// CollectSystemInfo recopila información completa del sistema
func CollectSystemInfo() *models.SystemInfo {
	info := &models.SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		CPUCount:     GetCPUCount(),
		EnvVars:      make(map[string]string),
	}

	// Memoria RAM
	if ram, err := GetMemoryInfo(); err == nil {
		info.TotalRAM = ram
	}

	// Disco
	if disk, err := GetDiskInfo(); err == nil {
		info.TotalDisk = disk
	}

	// BIOS
	info.BIOS = getBIOSInfo()

	// Procesos
	info.Processes = getProcessList()

	// Usuarios y grupos
	info.Users = getUsers()
	info.Groups = getGroups()

	// Conexiones de red
	info.NetworkConns = getNetworkConnections()

	// Servicios
	info.Services = getServices()

	// Variables de entorno
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			info.EnvVars[parts[0]] = parts[1]
		}
	}

	// Pipes
	info.Pipes = getPipes()

	// Posición del mouse
	info.MousePosition = getMousePosition()

	// Aplicaciones instaladas
	info.InstalledApps = getInstalledApps()

	// Archivos recientes
	info.RecentFiles = getRecentFiles()

	// Uptime
	info.UptimeSeconds = getSystemUptime()

	// Language y Timezone
	info.Language = getSystemLanguage()
	info.Timezone = getSystemTimezone()

	// Screenshot
	info.Screenshot = captureScreenshot()

	return info
}

func getBIOSInfo() string {
	// Obtener información del BIOS usando registro de Windows
	var key windows.Handle
	err := windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr("HARDWARE\\DESCRIPTION\\System\\BIOS"),
		0,
		windows.KEY_READ,
		&key,
	)
	if err != nil {
		return "Unknown"
	}
	defer windows.RegCloseKey(key)

	// Leer valores del BIOS
	biosVendor := readRegistryString(key, "BIOSVendor")
	biosVersion := readRegistryString(key, "BIOSVersion")
	biosDate := readRegistryString(key, "BIOSReleaseDate")

	if biosVendor != "" || biosVersion != "" {
		return fmt.Sprintf("%s %s (%s)", biosVendor, biosVersion, biosDate)
	}

	return "Unknown"
}

func readRegistryString(key windows.Handle, valueName string) string {
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

func getProcessList() []models.ProcessInfo {
	processes := []models.ProcessInfo{}

	// Snapshot de procesos
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
		pid := procEntry.ProcessID
		name := windows.UTF16ToString(procEntry.ExeFile[:])

		// Obtener path y owner del proceso
		path := getProcessPath(pid)
		owner := getProcessOwner(pid)

		processes = append(processes, models.ProcessInfo{
			PID:   pid,
			Name:  name,
			Owner: owner,
			Path:  path,
		})

		err = windows.Process32Next(snapshot, &procEntry)
		if err != nil {
			break
		}
	}

	return processes
}

func getProcessPath(pid uint32) string {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(handle)

	var pathBuf [windows.MAX_PATH]uint16
	size := uint32(len(pathBuf))

	err = windows.QueryFullProcessImageName(handle, 0, &pathBuf[0], &size)
	if err != nil {
		return ""
	}

	return windows.UTF16ToString(pathBuf[:size])
}

func getProcessOwner(pid uint32) string {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(handle)

	var token windows.Token
	err = windows.OpenProcessToken(handle, windows.TOKEN_QUERY, &token)
	if err != nil {
		return ""
	}
	defer token.Close()

	tokenUser, err := token.GetTokenUser()
	if err != nil {
		return ""
	}

	account, domain, _, err := tokenUser.User.Sid.LookupAccount("")
	if err != nil {
		return ""
	}

	if domain != "" {
		return fmt.Sprintf("%s\\%s", domain, account)
	}
	return account
}

func getUsers() []string {
	users := []string{}

	// Usar NetUserEnum para enumerar usuarios locales
	netapi32 := windows.NewLazyDLL("netapi32.dll")
	netUserEnum := netapi32.NewProc("NetUserEnum")
	netApiBufferFree := netapi32.NewProc("NetApiBufferFree")

	var (
		dataPointer  uintptr
		entriesRead  uint32
		totalEntries uint32
		resumeHandle uint32
	)

	ret, _, _ := netUserEnum.Call(
		0, // servername (NULL = local)
		0, // level (0 = USER_INFO_0)
		0, // filter (0 = all users)
		uintptr(unsafe.Pointer(&dataPointer)),
		uintptr(0xFFFFFFFF), // prefmaxlen
		uintptr(unsafe.Pointer(&entriesRead)),
		uintptr(unsafe.Pointer(&totalEntries)),
		uintptr(unsafe.Pointer(&resumeHandle)),
	)

	if ret != 0 {
		return users
	}
	defer netApiBufferFree.Call(dataPointer)

	// USER_INFO_0 structure
	type USER_INFO_0 struct {
		Name *uint16
	}

	// Iterar sobre los usuarios
	size := unsafe.Sizeof(USER_INFO_0{})
	for i := uint32(0); i < entriesRead; i++ {
		userInfo := (*USER_INFO_0)(unsafe.Pointer(dataPointer + uintptr(i)*size))
		if userInfo.Name != nil {
			users = append(users, windows.UTF16PtrToString(userInfo.Name))
		}
	}

	return users
}

func getGroups() []string {
	groups := []string{}

	// Usar NetLocalGroupEnum para enumerar grupos locales
	netapi32 := windows.NewLazyDLL("netapi32.dll")
	netLocalGroupEnum := netapi32.NewProc("NetLocalGroupEnum")
	netApiBufferFree := netapi32.NewProc("NetApiBufferFree")

	var (
		dataPointer  uintptr
		entriesRead  uint32
		totalEntries uint32
		resumeHandle uint32
	)

	ret, _, _ := netLocalGroupEnum.Call(
		0, // servername (NULL = local)
		0, // level (0 = LOCALGROUP_INFO_0)
		uintptr(unsafe.Pointer(&dataPointer)),
		uintptr(0xFFFFFFFF), // prefmaxlen
		uintptr(unsafe.Pointer(&entriesRead)),
		uintptr(unsafe.Pointer(&totalEntries)),
		uintptr(unsafe.Pointer(&resumeHandle)),
	)

	if ret != 0 {
		return groups
	}
	defer netApiBufferFree.Call(dataPointer)

	// LOCALGROUP_INFO_0 structure
	type LOCALGROUP_INFO_0 struct {
		Name *uint16
	}

	size := unsafe.Sizeof(LOCALGROUP_INFO_0{})
	for i := uint32(0); i < entriesRead; i++ {
		groupInfo := (*LOCALGROUP_INFO_0)(unsafe.Pointer(dataPointer + uintptr(i)*size))
		if groupInfo.Name != nil {
			groups = append(groups, windows.UTF16PtrToString(groupInfo.Name))
		}
	}

	return groups
}

func getNetworkConnections() []models.NetworkConn {
	connections := []models.NetworkConn{}

	// Obtener conexiones TCP
	tcpConns := getTCPConnections()
	connections = append(connections, tcpConns...)

	// Obtener conexiones UDP
	udpConns := getUDPConnections()
	connections = append(connections, udpConns...)

	return connections
}

func getTCPConnections() []models.NetworkConn {
	connections := []models.NetworkConn{}

	iphlpapi := windows.NewLazyDLL("iphlpapi.dll")
	getExtendedTcpTable := iphlpapi.NewProc("GetExtendedTcpTable")

	var size uint32
	// Primera llamada para obtener el tamaño necesario
	getExtendedTcpTable.Call(
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		windows.AF_INET,
		5, // TCP_TABLE_OWNER_PID_ALL
		0,
	)

	if size == 0 {
		return connections
	}

	buffer := make([]byte, size)
	ret, _, _ := getExtendedTcpTable.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		windows.AF_INET,
		5, // TCP_TABLE_OWNER_PID_ALL
		0,
	)

	if ret != 0 {
		return connections
	}

	// MIB_TCPTABLE_OWNER_PID structure
	type MIB_TCPROW_OWNER_PID struct {
		State      uint32
		LocalAddr  uint32
		LocalPort  uint32
		RemoteAddr uint32
		RemotePort uint32
		OwningPid  uint32
	}

	numEntries := *(*uint32)(unsafe.Pointer(&buffer[0]))
	rowSize := unsafe.Sizeof(MIB_TCPROW_OWNER_PID{})

	for i := uint32(0); i < numEntries && i < 100; i++ { // Limitar a 100 conexiones
		offset := 4 + uintptr(i)*rowSize
		row := (*MIB_TCPROW_OWNER_PID)(unsafe.Pointer(&buffer[offset]))

		localAddr := fmt.Sprintf("%d.%d.%d.%d:%d",
			byte(row.LocalAddr),
			byte(row.LocalAddr>>8),
			byte(row.LocalAddr>>16),
			byte(row.LocalAddr>>24),
			ntohs(uint16(row.LocalPort)),
		)

		remoteAddr := fmt.Sprintf("%d.%d.%d.%d:%d",
			byte(row.RemoteAddr),
			byte(row.RemoteAddr>>8),
			byte(row.RemoteAddr>>16),
			byte(row.RemoteAddr>>24),
			ntohs(uint16(row.RemotePort)),
		)

		connections = append(connections, models.NetworkConn{
			Protocol:   "TCP",
			LocalAddr:  localAddr,
			RemoteAddr: remoteAddr,
			State:      getTCPState(row.State),
		})
	}

	return connections
}

func getUDPConnections() []models.NetworkConn {
	connections := []models.NetworkConn{}

	iphlpapi := windows.NewLazyDLL("iphlpapi.dll")
	getExtendedUdpTable := iphlpapi.NewProc("GetExtendedUdpTable")

	var size uint32
	getExtendedUdpTable.Call(
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		windows.AF_INET,
		1, // UDP_TABLE_OWNER_PID
		0,
	)

	if size == 0 {
		return connections
	}

	buffer := make([]byte, size)
	ret, _, _ := getExtendedUdpTable.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		windows.AF_INET,
		1, // UDP_TABLE_OWNER_PID
		0,
	)

	if ret != 0 {
		return connections
	}

	type MIB_UDPROW_OWNER_PID struct {
		LocalAddr uint32
		LocalPort uint32
		OwningPid uint32
	}

	numEntries := *(*uint32)(unsafe.Pointer(&buffer[0]))
	rowSize := unsafe.Sizeof(MIB_UDPROW_OWNER_PID{})

	for i := uint32(0); i < numEntries && i < 100; i++ {
		offset := 4 + uintptr(i)*rowSize
		row := (*MIB_UDPROW_OWNER_PID)(unsafe.Pointer(&buffer[offset]))

		localAddr := fmt.Sprintf("%d.%d.%d.%d:%d",
			byte(row.LocalAddr),
			byte(row.LocalAddr>>8),
			byte(row.LocalAddr>>16),
			byte(row.LocalAddr>>24),
			ntohs(uint16(row.LocalPort)),
		)

		connections = append(connections, models.NetworkConn{
			Protocol:   "UDP",
			LocalAddr:  localAddr,
			RemoteAddr: "*:*",
			State:      "LISTENING",
		})
	}

	return connections
}

func ntohs(port uint16) uint16 {
	return (port>>8)&0xff | (port<<8)&0xff00
}

func getTCPState(state uint32) string {
	states := map[uint32]string{
		1:  "CLOSED",
		2:  "LISTENING",
		3:  "SYN_SENT",
		4:  "SYN_RECEIVED",
		5:  "ESTABLISHED",
		6:  "FIN_WAIT_1",
		7:  "FIN_WAIT_2",
		8:  "CLOSE_WAIT",
		9:  "CLOSING",
		10: "LAST_ACK",
		11: "TIME_WAIT",
		12: "DELETE_TCB",
	}

	if s, ok := states[state]; ok {
		return s
	}
	return "UNKNOWN"
}

func getServices() []string {
	services := []string{}

	// Abrir Service Control Manager
	scm, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_ENUMERATE_SERVICE)
	if err != nil {
		return services
	}
	defer windows.CloseServiceHandle(scm)

	var (
		bytesNeeded      uint32
		servicesReturned uint32
		resumeHandle     uint32
	)

	// Primera llamada para obtener el tamaño necesario
	windows.EnumServicesStatusEx(
		scm,
		windows.SC_ENUM_PROCESS_INFO,
		windows.SERVICE_WIN32,
		windows.SERVICE_STATE_ALL,
		nil,
		0,
		&bytesNeeded,
		&servicesReturned,
		&resumeHandle,
		nil,
	)

	if bytesNeeded == 0 {
		return services
	}

	buffer := make([]byte, bytesNeeded)
	err = windows.EnumServicesStatusEx(
		scm,
		windows.SC_ENUM_PROCESS_INFO,
		windows.SERVICE_WIN32,
		windows.SERVICE_STATE_ALL,
		&buffer[0],
		bytesNeeded,
		&bytesNeeded,
		&servicesReturned,
		&resumeHandle,
		nil,
	)

	if err != nil {
		return services
	}

	// Parsear los servicios
	type ENUM_SERVICE_STATUS_PROCESS struct {
		ServiceName          *uint16
		DisplayName          *uint16
		ServiceStatusProcess windows.SERVICE_STATUS_PROCESS
	}

	size := unsafe.Sizeof(ENUM_SERVICE_STATUS_PROCESS{})
	for i := uint32(0); i < servicesReturned && i < 200; i++ { // Limitar a 200 servicios
		offset := uintptr(i) * size
		service := (*ENUM_SERVICE_STATUS_PROCESS)(unsafe.Pointer(&buffer[offset]))

		if service.ServiceName != nil {
			serviceName := windows.UTF16PtrToString(service.ServiceName)
			services = append(services, serviceName)
		}
	}

	return services
}

func getPipes() []string {
	pipes := []string{}

	// Buscar named pipes en \\.\pipe\
	pipePath := `\\.\pipe\*`

	var findData windows.Win32finddata
	handle, err := windows.FindFirstFile(windows.StringToUTF16Ptr(pipePath), &findData)
	if err != nil {
		return pipes
	}
	defer windows.FindClose(handle)

	pipes = append(pipes, windows.UTF16ToString(findData.FileName[:]))

	// Continuar buscando
	for i := 0; i < 100; i++ { // Limitar a 100 pipes
		err = windows.FindNextFile(handle, &findData)
		if err != nil {
			break
		}
		pipeName := windows.UTF16ToString(findData.FileName[:])
		if pipeName != "." && pipeName != ".." {
			pipes = append(pipes, pipeName)
		}
	}

	return pipes
}

func getMousePosition() models.Point {
	user32 := windows.NewLazyDLL("user32.dll")
	getCursorPos := user32.NewProc("GetCursorPos")

	var point models.Point
	getCursorPos.Call(uintptr(unsafe.Pointer(&point)))

	return point
}

func getInstalledApps() []string {
	apps := []string{}

	// Rutas de registro donde se almacenan las aplicaciones instaladas
	registryPaths := []struct {
		root windows.Handle
		path string
	}{
		{windows.HKEY_LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
		{windows.HKEY_LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`},
		{windows.HKEY_CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
	}

	for _, regPath := range registryPaths {
		apps = append(apps, getAppsFromRegistry(regPath.root, regPath.path)...)
	}

	return apps
}

func getAppsFromRegistry(root windows.Handle, path string) []string {
	apps := []string{}

	var key windows.Handle
	err := windows.RegOpenKeyEx(root, windows.StringToUTF16Ptr(path), 0, windows.KEY_READ, &key)
	if err != nil {
		return apps
	}
	defer windows.RegCloseKey(key)

	// Enumerar subclaves
	var index uint32
	for i := 0; i < 500; i++ { // Limitar a 500 aplicaciones
		var nameBuf [256]uint16
		nameLen := uint32(len(nameBuf))

		err := windows.RegEnumKeyEx(key, index, &nameBuf[0], &nameLen, nil, nil, nil, nil)
		if err != nil {
			break
		}

		subKeyName := windows.UTF16ToString(nameBuf[:nameLen])

		// Abrir subclave para leer DisplayName
		var subKey windows.Handle
		subKeyPath := path + `\` + subKeyName
		err = windows.RegOpenKeyEx(root, windows.StringToUTF16Ptr(subKeyPath), 0, windows.KEY_READ, &subKey)
		if err == nil {
			displayName := readRegistryString(subKey, "DisplayName")
			if displayName != "" {
				apps = append(apps, displayName)
			}
			windows.RegCloseKey(subKey)
		}

		index++
	}

	return apps
}

func getRecentFiles() []string {
	files := []string{}

	// Obtener la carpeta Recent del usuario
	recentPath := os.Getenv("APPDATA") + `\Microsoft\Windows\Recent`

	entries, err := os.ReadDir(recentPath)
	if err != nil {
		return files
	}

	for i, entry := range entries {
		if i >= 50 { // Limitar a 50 archivos recientes
			break
		}
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files
}

func getSystemUptime() int64 {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	getTickCount64 := kernel32.NewProc("GetTickCount64")

	ret, _, _ := getTickCount64.Call()
	milliseconds := int64(ret)

	return milliseconds / 1000 // Convertir a segundos
}

func captureScreenshot() string {
	user32 := windows.NewLazyDLL("user32.dll")
	gdi32 := windows.NewLazyDLL("gdi32.dll")

	getDC := user32.NewProc("GetDC")
	createCompatibleDC := gdi32.NewProc("CreateCompatibleDC")
	createCompatibleBitmap := gdi32.NewProc("CreateCompatibleBitmap")
	selectObject := gdi32.NewProc("SelectObject")
	bitBlt := gdi32.NewProc("BitBlt")
	getSystemMetrics := user32.NewProc("GetSystemMetrics")
	getDIBits := gdi32.NewProc("GetDIBits")
	deleteDC := gdi32.NewProc("DeleteDC")
	deleteObject := gdi32.NewProc("DeleteObject")
	releaseDC := user32.NewProc("ReleaseDC")

	screenWidth, _, _ := getSystemMetrics.Call(0)
	screenHeight, _, _ := getSystemMetrics.Call(1)

	if screenWidth == 0 || screenHeight == 0 {
		return ""
	}

	hdcScreen, _, _ := getDC.Call(0)
	if hdcScreen == 0 {
		return ""
	}
	defer releaseDC.Call(0, hdcScreen)

	hdcMem, _, _ := createCompatibleDC.Call(hdcScreen)
	if hdcMem == 0 {
		return ""
	}
	defer deleteDC.Call(hdcMem)

	hBitmap, _, _ := createCompatibleBitmap.Call(hdcScreen, screenWidth, screenHeight)
	if hBitmap == 0 {
		return ""
	}
	defer deleteObject.Call(hBitmap)

	selectObject.Call(hdcMem, hBitmap)

	const SRCCOPY = 0x00CC0020
	bitBlt.Call(hdcMem, 0, 0, screenWidth, screenHeight, hdcScreen, 0, 0, SRCCOPY)

	type BITMAPINFOHEADER struct {
		BiSize          uint32
		BiWidth         int32
		BiHeight        int32
		BiPlanes        uint16
		BiBitCount      uint16
		BiCompression   uint32
		BiSizeImage     uint32
		BiXPelsPerMeter int32
		BiYPelsPerMeter int32
		BiClrUsed       uint32
		BiClrImportant  uint32
	}

	type BITMAPINFO struct {
		BmiHeader BITMAPINFOHEADER
		BmiColors [1]uint32
	}

	bi := BITMAPINFO{}
	bi.BmiHeader.BiSize = uint32(unsafe.Sizeof(bi.BmiHeader))
	bi.BmiHeader.BiWidth = int32(screenWidth)
	bi.BmiHeader.BiHeight = -int32(screenHeight)
	bi.BmiHeader.BiPlanes = 1
	bi.BmiHeader.BiBitCount = 32
	bi.BmiHeader.BiCompression = 0

	bufferSize := int(screenWidth) * int(screenHeight) * 4
	buffer := make([]byte, bufferSize)

	ret, _, _ := getDIBits.Call(
		hdcMem,
		hBitmap,
		0,
		uintptr(screenHeight),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&bi)),
		0,
	)

	if ret == 0 {
		return ""
	}

	img := image.NewRGBA(image.Rect(0, 0, int(screenWidth), int(screenHeight)))

	for y := 0; y < int(screenHeight); y++ {
		for x := 0; x < int(screenWidth); x++ {
			offset := (y*int(screenWidth) + x) * 4
			img.Pix[offset+0] = buffer[offset+2]
			img.Pix[offset+1] = buffer[offset+1]
			img.Pix[offset+2] = buffer[offset+0]
			img.Pix[offset+3] = buffer[offset+3]
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func getSystemLanguage() string {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	getUserDefaultLocaleName := kernel32.NewProc("GetUserDefaultLocaleName")

	var localeName [85]uint16
	ret, _, _ := getUserDefaultLocaleName.Call(
		uintptr(unsafe.Pointer(&localeName[0])),
		uintptr(len(localeName)),
	)

	if ret == 0 {
		return "Unknown"
	}

	return windows.UTF16ToString(localeName[:])
}

func getSystemTimezone() string {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	getTimeZoneInformation := kernel32.NewProc("GetTimeZoneInformation")

	type TIME_ZONE_INFORMATION struct {
		Bias         int32
		StandardName [32]uint16
		StandardDate [16]byte
		StandardBias int32
		DaylightName [32]uint16
		DaylightDate [16]byte
		DaylightBias int32
	}

	var tzi TIME_ZONE_INFORMATION
	ret, _, _ := getTimeZoneInformation.Call(uintptr(unsafe.Pointer(&tzi)))

	if ret == 0xFFFFFFFF {
		return "Unknown"
	}

	return windows.UTF16ToString(tzi.StandardName[:])
}
