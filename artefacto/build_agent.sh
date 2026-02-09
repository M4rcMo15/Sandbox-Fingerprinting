#!/bin/bash

# Nombre del ejecutable por defecto o pasado por parámetro
OUTPUT_NAME=${1:-agent.exe}

echo "[+] Compilando agente como: $OUTPUT_NAME..."

# Asegurar que estamos en el directorio del script
cd "$(dirname "$0")"

# Variables de entorno para Cross-Compilation a Windows
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=0

# Comando de compilación "Malware Style" (Stripped)
# -ldflags="-s -w": Elimina símbolos y debug info (reduce tamaño y aumenta entropía/sospecha).
# -trimpath: Oculta rutas de compilación.
go build -ldflags="-s -w" -trimpath -o "$OUTPUT_NAME" .

if [ $? -eq 0 ]; then
    echo "[+] Compilación completada: $OUTPUT_NAME"
    ls -lh "$OUTPUT_NAME"
else
    echo "[-] Error al compilar"
    exit 1
fi
