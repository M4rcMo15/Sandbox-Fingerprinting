#!/bin/bash

# Script de compilación para el agente de detección de sandbox

echo "[*] Compilando agente de detección de sandbox..."

# Colores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Limpiar builds anteriores
echo -e "${BLUE}[+] Limpiando builds anteriores...${NC}"
rm -f agent.exe agent_optimized.exe agent_obfuscated.exe

# Build básico
echo -e "${BLUE}[+] Compilación básica...${NC}"
GOOS=windows GOARCH=amd64 go build -o agent.exe
if [ $? -eq 0 ]; then
    SIZE=$(du -h agent.exe | cut -f1)
    echo -e "${GREEN}[✓] agent.exe creado (${SIZE})${NC}"
else
    echo "[!] Error en compilación básica"
    exit 1
fi

# Build optimizado
echo -e "${BLUE}[+] Compilación optimizada (stripped)...${NC}"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o agent_optimized.exe
if [ $? -eq 0 ]; then
    SIZE=$(du -h agent_optimized.exe | cut -f1)
    echo -e "${GREEN}[✓] agent_optimized.exe creado (${SIZE})${NC}"
else
    echo "[!] Error en compilación optimizada"
fi

# Build con UPX (si está instalado)
if command -v upx &> /dev/null; then
    echo -e "${BLUE}[+] Comprimiendo con UPX...${NC}"
    cp agent_optimized.exe agent_compressed.exe
    upx --best --lzma agent_compressed.exe 2>/dev/null
    if [ $? -eq 0 ]; then
        SIZE=$(du -h agent_compressed.exe | cut -f1)
        echo -e "${GREEN}[✓] agent_compressed.exe creado (${SIZE})${NC}"
    fi
else
    echo "[i] UPX no instalado, saltando compresión"
fi

# Build ofuscado (si garble está instalado)
if command -v garble &> /dev/null; then
    echo -e "${BLUE}[+] Compilación ofuscada con garble...${NC}"
    GOOS=windows GOARCH=amd64 garble -literals -tiny build -o agent_obfuscated.exe
    if [ $? -eq 0 ]; then
        SIZE=$(du -h agent_obfuscated.exe | cut -f1)
        echo -e "${GREEN}[✓] agent_obfuscated.exe creado (${SIZE})${NC}"
    fi
else
    echo "[i] Garble no instalado, saltando ofuscación"
    echo "[i] Instalar con: go install mvdan.cc/garble@latest"
fi

echo ""
echo -e "${GREEN}[✓] Compilación completada${NC}"
echo ""
echo "Archivos generados:"
ls -lh *.exe 2>/dev/null | awk '{print "  - " $9 " (" $5 ")"}'
echo ""
echo -e "${BLUE}[i] Configuración actual:${NC}"
if [ -f .env ]; then
    grep "SERVER_URL" .env | head -1
else
    echo "  No se encontró archivo .env"
fi
echo ""
echo -e "${GREEN}[✓] Listo para ejecutar${NC}"
echo "  Verifica los resultados en: http://54.37.226.179"
