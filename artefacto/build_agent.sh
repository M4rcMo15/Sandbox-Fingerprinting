#!/bin/bash

# Nombre del ejecutable por defecto o pasado por parámetro
OUTPUT_NAME=${1:-agent.exe}

echo "[+] Compilando agente como: $OUTPUT_NAME..."
echo "[+] Modo: MÁXIMA DETECCIÓN (sin ofuscación)"

# Asegurar que estamos en el directorio del script
cd "$(dirname "$0")"

# Variables de entorno para Cross-Compilation a Windows
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=0

# Comando de compilación SIN ofuscación para MÁXIMA DETECCIÓN
# NO usamos -s -w -trimpath para mantener símbolos, debug info y strings visibles
go build -o "$OUTPUT_NAME" .

if [ $? -eq 0 ]; then
    echo "[+] Compilación completada: $OUTPUT_NAME"
    ls -lh "$OUTPUT_NAME"
    echo "[+] Binario compilado en modo MÁXIMA DETECCIÓN"
    echo "[+] Símbolos: PRESENTES"
    echo "[+] Debug Info: PRESENTE"
    echo "[+] Strings: VISIBLES"
else
    echo "[-] Error al compilar"
    exit 1
fi
