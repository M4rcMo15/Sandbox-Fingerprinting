#!/bin/bash

# Nombre del ejecutable por defecto o pasado por parámetro
OUTPUT_NAME=${1:-agent.exe}
TARGET_SANDBOX=${2:-ANY_RUN}

echo "[+] Compilando agente como: $OUTPUT_NAME..."
echo "[+] Target Sandbox: $TARGET_SANDBOX"
echo "[+] Modo: MÁXIMA DETECCIÓN (sin ofuscación)"

# Asegurar que estamos en el directorio del script
cd "$(dirname "$0")"

# Variables de entorno para Cross-Compilation a Windows
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=0

# Comando de compilación SIN ofuscación para MÁXIMA DETECCIÓN
# NO usamos -s -w -trimpath para mantener símbolos, debug info y strings visibles
# Inyectamos el target sandbox en tiempo de compilación
go build -ldflags "-X main.targetSandbox=$TARGET_SANDBOX" -o "$OUTPUT_NAME" .

if [ $? -eq 0 ]; then
    echo "[+] Compilación completada: $OUTPUT_NAME"
    ls -lh "$OUTPUT_NAME"
    echo "[+] Binario compilado en modo MÁXIMA DETECCIÓN"
    echo "[+] Target Sandbox: $TARGET_SANDBOX"
    echo "[+] Símbolos: PRESENTES"
    echo "[+] Debug Info: PRESENTE"
    echo "[+] Strings: VISIBLES"
    echo ""
    echo "[+] Nuevos vectores de inyección añadidos:"
    echo "    - Registry NLS (T1012 - Checks supported languages)"
    echo "    - Computer Name queries (Reads the computer name)"
    echo "    - DNS queries con nslookup (Uses NSLOOKUP.EXE)"
    echo "    - Service Management con SC.EXE"
    echo "    - Task Scheduler avanzado (schtasks)"
    echo "    - WMI queries extensivas"
    echo "    - NET commands (user, localgroup, share, etc.)"
    echo "    - NETSH commands (firewall, interface, wlan)"
    echo "    - PowerShell commands adicionales (Get-ComputerInfo, Get-Process, etc.)"
    echo "    - CMD commands adicionales (systeminfo, tasklist, whoami, etc.)"
else
    echo "[-] Error al compilar"
    exit 1
fi
