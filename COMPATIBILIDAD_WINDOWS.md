# 🪟 Compatibilidad con Windows

## ✅ Versiones de Windows Soportadas

El agente está diseñado para Windows y debería funcionar en:

### Completamente Soportado ✅
- **Windows 10** (32-bit y 64-bit)
- **Windows 11** (64-bit)
- **Windows Server 2016**
- **Windows Server 2019**
- **Windows Server 2022**

### Parcialmente Soportado ⚠️
- **Windows 8.1** - Funciona pero algunas APIs pueden fallar
- **Windows 7** - Funciona pero con limitaciones
- **Windows Server 2012 R2** - Funciona con limitaciones

### No Soportado ❌
- **Windows XP** - APIs no disponibles
- **Windows Vista** - APIs no disponibles
- **Linux** - Código específico de Windows
- **macOS** - Código específico de Windows

## 🔍 Funcionalidades por Versión

### Windows 10/11 (Recomendado)

| Funcionalidad | Windows 10 | Windows 11 | Notas |
|---------------|------------|------------|-------|
| Detección de Sandbox | ✅ | ✅ | Completo |
| Información del Sistema | ✅ | ✅ | Completo |
| Detección de Hooks | ✅ | ✅ | Completo |
| Screenshot | ✅ | ✅ | Completo |
| Detección de EDR | ✅ | ✅ | Completo |
| Crawler de Archivos | ✅ | ✅ | Completo |
| Detección de Herramientas | ✅ | ✅ | Completo |
| Geolocalización | ✅ | ✅ | Completo |
| Exfiltración | ✅ | ✅ | Completo |

### Windows 7/8.1

| Funcionalidad | Windows 7 | Windows 8.1 | Notas |
|---------------|-----------|-------------|-------|
| Detección de Sandbox | ✅ | ✅ | Completo |
| Información del Sistema | ⚠️ | ✅ | Algunas APIs limitadas |
| Detección de Hooks | ✅ | ✅ | Completo |
| Screenshot | ⚠️ | ✅ | Puede fallar en algunos casos |
| Detección de EDR | ✅ | ✅ | Completo |
| Crawler de Archivos | ✅ | ✅ | Completo |
| Detección de Herramientas | ✅ | ✅ | Completo |
| Geolocalización | ✅ | ✅ | Completo |
| Exfiltración | ✅ | ✅ | Completo |

## 🏗️ Arquitecturas Soportadas

### 64-bit (x64) ✅ Recomendado
```bash
# Compilar para 64-bit
set GOARCH=amd64
go build -o agent_x64.exe -ldflags="-s -w"
```

### 32-bit (x86) ✅ Soportado
```bash
# Compilar para 32-bit
set GOARCH=386
go build -o agent_x86.exe -ldflags="-s -w"
```

### ARM64 ⚠️ Experimental
```bash
# Compilar para ARM64 (Windows 11 ARM)
set GOARCH=arm64
go build -o agent_arm64.exe -ldflags="-s -w"
```

## 🧪 Testing en Diferentes Versiones

### Windows 10 (Tu Sistema Actual)

```powershell
# Compilar
cd artefacto
go build -o agent.exe -ldflags="-s -w"

# Ejecutar
.\agent.exe

# Debería funcionar perfectamente ✅
```

### Windows 11

```powershell
# Mismo proceso que Windows 10
# Sin cambios necesarios
# Funciona igual ✅
```

### Windows 7 (Si necesitas probarlo)

```powershell
# Compilar con compatibilidad
set GOOS=windows
set GOARCH=amd64
go build -o agent_win7.exe -ldflags="-s -w"

# Ejecutar en Windows 7
# Puede mostrar advertencias pero debería funcionar ⚠️
```

## ⚠️ Limitaciones Conocidas

### Windows 7
1. **Screenshot** - Puede fallar en algunos casos
   - API de captura de pantalla limitada
   - Solución: El agente continúa sin screenshot

2. **Información del Sistema** - Algunas métricas no disponibles
   - Uptime puede ser incorrecto
   - Algunas variables de entorno pueden faltar

3. **Detección de Herramientas** - Lista limitada
   - Menos herramientas detectables
   - Paths diferentes

### Windows Server
1. **Screenshot** - Puede fallar si no hay GUI
   - Servidores sin escritorio
   - Solución: El agente continúa sin screenshot

2. **Aplicaciones Instaladas** - Lista diferente
   - Menos aplicaciones típicas
   - Más servicios y roles

## 🔧 Requisitos del Sistema

### Mínimos
- **OS:** Windows 7 SP1 o superior
- **RAM:** 512 MB
- **Disco:** 50 MB libres
- **Red:** Conexión a Internet

### Recomendados
- **OS:** Windows 10/11
- **RAM:** 2 GB
- **Disco:** 100 MB libres
- **Red:** Conexión estable

## 🛡️ Permisos Necesarios

### Usuario Normal ✅
El agente funciona con permisos de usuario normal:
- ✅ Recopila información básica
- ✅ Detecta sandbox
- ✅ Captura screenshot
- ✅ Exfiltra datos
- ⚠️ Algunas funciones limitadas

### Administrador ⭐ Recomendado
Con permisos de administrador:
- ✅ Información completa del sistema
- ✅ Todos los procesos visibles
- ✅ Todos los servicios visibles
- ✅ Acceso a más archivos
- ✅ Detección completa de EDR

## 🧪 Verificar Compatibilidad

### Script de Verificación

```powershell
# Verificar versión de Windows
systeminfo | findstr /B /C:"OS Name" /C:"OS Version"

# Verificar arquitectura
wmic os get osarchitecture

# Verificar .NET (si es necesario)
reg query "HKLM\SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full" /v Version

# Verificar conectividad
ping 54.37.226.179
curl http://54.37.226.179
```

## 🐛 Problemas Comunes

### Error: "The program can't start because api-ms-win-*.dll is missing"

**Causa:** Windows 7 sin actualizaciones

**Solución:**
```powershell
# Instalar actualizaciones de Windows 7
# O compilar con compatibilidad:
set CGO_ENABLED=0
go build -o agent.exe -ldflags="-s -w"
```

### Error: "Screenshot failed"

**Causa:** No hay GUI o permisos insuficientes

**Solución:**
- El agente continúa sin screenshot
- No afecta otras funcionalidades
- Normal en Windows Server Core

### Error: "Access denied" al leer procesos

**Causa:** Permisos insuficientes

**Solución:**
```powershell
# Ejecutar como administrador
# Click derecho > "Run as administrator"
```

## 📊 Sandboxes por Versión de Windows

### VirusTotal
- Windows 7 32-bit
- Windows 7 64-bit
- Windows 10 64-bit

### Hybrid Analysis
- Windows 7 32-bit
- Windows 7 64-bit
- Windows 10 64-bit ⭐ Recomendado

### Any.Run
- Windows 7 32-bit
- Windows 7 64-bit
- Windows 10 64-bit ⭐ Recomendado
- Windows 11 64-bit

## ✅ Recomendaciones

### Para Desarrollo y Testing
```
Usa: Windows 10/11 64-bit
Por qué: Máxima compatibilidad y funcionalidades
```

### Para Sandboxes
```
Elige: Windows 10 64-bit en la sandbox
Por qué: Más común, mejor soporte, más datos
```

### Para Producción
```
Compila: Para 64-bit
Por qué: Más común en sistemas modernos
Fallback: Compila también 32-bit para compatibilidad
```

## 🎯 Compilación Multi-Versión

```bash
# Script para compilar todas las versiones
# compile_all.bat

@echo off
echo Compilando para todas las versiones...

REM Windows 10/11 64-bit (principal)
set GOOS=windows
set GOARCH=amd64
go build -o agent_win10_x64.exe -ldflags="-s -w"
echo [OK] Windows 10/11 64-bit

REM Windows 10/11 32-bit
set GOARCH=386
go build -o agent_win10_x86.exe -ldflags="-s -w"
echo [OK] Windows 10/11 32-bit

REM Windows 7 64-bit (compatibilidad)
set GOARCH=amd64
set CGO_ENABLED=0
go build -o agent_win7_x64.exe -ldflags="-s -w"
echo [OK] Windows 7 64-bit

echo.
echo Compilacion completada!
dir *.exe
```

## 📝 Resumen

### ✅ Funciona Perfectamente
- Windows 10 (32-bit y 64-bit)
- Windows 11 (64-bit)
- Windows Server 2016+

### ⚠️ Funciona con Limitaciones
- Windows 7 (algunas APIs limitadas)
- Windows 8.1 (funciona bien)
- Windows Server 2012 R2

### ❌ No Funciona
- Windows XP/Vista
- Linux/macOS

### 🎯 Recomendación
**Usa Windows 10 64-bit** para máxima compatibilidad y funcionalidades completas.

---

**Tu sistema actual:** Windows 10 ✅  
**Compatibilidad:** 100% ✅  
**Funcionalidades:** Todas disponibles ✅
