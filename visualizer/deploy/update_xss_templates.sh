#!/bin/bash

# Script de actualización de templates XSS - Phase 1
# Actualiza las templates para soportar los nuevos vectores XSS

set -e  # Exit on error

echo "================================================"
echo "  Actualización Templates XSS - Phase 1"
echo "================================================"
echo ""

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Directorio base
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
TEMPLATES_DIR="$PROJECT_DIR/xss_audit/templates/xss_audit"

echo -e "${YELLOW}[1/6]${NC} Verificando directorio de templates..."
if [ ! -d "$TEMPLATES_DIR" ]; then
    echo -e "${RED}Error: No se encuentra el directorio de templates${NC}"
    echo "Buscado en: $TEMPLATES_DIR"
    exit 1
fi
echo -e "${GREEN}✓${NC} Directorio encontrado: $TEMPLATES_DIR"
echo ""

echo -e "${YELLOW}[2/6]${NC} Creando backup de templates actuales..."
BACKUP_DIR="$PROJECT_DIR/backups/templates_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"
cp -r "$TEMPLATES_DIR"/* "$BACKUP_DIR/" 2>/dev/null || true
echo -e "${GREEN}✓${NC} Backup creado en: $BACKUP_DIR"
echo ""

echo -e "${YELLOW}[3/6]${NC} Verificando archivo dashboard.html..."
DASHBOARD_FILE="$TEMPLATES_DIR/dashboard.html"
if [ ! -f "$DASHBOARD_FILE" ]; then
    echo -e "${RED}Error: No se encuentra dashboard.html${NC}"
    exit 1
fi
echo -e "${GREEN}✓${NC} Archivo encontrado"
echo ""

echo -e "${YELLOW}[4/6]${NC} Verificando si ya están los nuevos badges..."
if grep -q "badge-dns" "$DASHBOARD_FILE"; then
    echo -e "${YELLOW}⚠${NC} Los badges ya parecen estar actualizados"
    read -p "¿Deseas continuar de todas formas? (s/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Ss]$ ]]; then
        echo "Operación cancelada"
        exit 0
    fi
else
    echo -e "${GREEN}✓${NC} Archivo necesita actualización"
fi
echo ""

echo -e "${YELLOW}[5/6]${NC} Actualizando badges CSS..."

# Buscar la línea donde está .badge-cmdline y añadir los nuevos badges después
if grep -q "\.badge-cmdline" "$DASHBOARD_FILE"; then
    # Crear archivo temporal con los nuevos badges
    NEW_BADGES=".badge-dns { background: #0969da; color: white; }
.badge-http { background: #1a7f37; color: white; }
.badge-file-content { background: #bf3989; color: white; }
.badge-pe-metadata { background: #bc4c00; color: white; }
.badge-environment { background: #8250df; color: white; }"
    
    # Usar sed para insertar después de .badge-cmdline
    sed -i.bak "/\.badge-cmdline.*{.*}/a\\
$NEW_BADGES" "$DASHBOARD_FILE"
    
    # Verificar que se añadieron
    if grep -q "badge-dns" "$DASHBOARD_FILE"; then
        echo -e "${GREEN}✓${NC} Badges CSS actualizados correctamente"
        rm -f "$DASHBOARD_FILE.bak"
    else
        echo -e "${RED}Error: No se pudieron añadir los badges${NC}"
        echo "Restaurando desde backup..."
        cp "$BACKUP_DIR/dashboard.html" "$DASHBOARD_FILE"
        exit 1
    fi
else
    echo -e "${RED}Error: No se encuentra .badge-cmdline en el archivo${NC}"
    echo "El archivo puede tener una estructura diferente"
    exit 1
fi
echo ""

echo -e "${YELLOW}[6/6]${NC} Reiniciando servidor Django..."

# Detectar cómo está corriendo el servidor
if systemctl is-active --quiet gunicorn; then
    echo "Detectado: systemd + gunicorn"
    sudo systemctl restart gunicorn
    echo -e "${GREEN}✓${NC} Gunicorn reiniciado"
elif pgrep -f "gunicorn" > /dev/null; then
    echo "Detectado: gunicorn manual"
    pkill -HUP gunicorn
    echo -e "${GREEN}✓${NC} Gunicorn recargado (HUP signal)"
elif pgrep -f "manage.py runserver" > /dev/null; then
    echo "Detectado: Django development server"
    echo -e "${YELLOW}⚠${NC} Servidor de desarrollo detectado"
    echo "Por favor reinicia manualmente el servidor"
else
    echo -e "${YELLOW}⚠${NC} No se detectó servidor corriendo"
    echo "Por favor inicia el servidor manualmente"
fi
echo ""

echo "================================================"
echo -e "${GREEN}✓ Actualización completada exitosamente${NC}"
echo "================================================"
echo ""
echo "Nuevos vectores añadidos:"
echo "  • dns (DNS Queries) - Azul oscuro"
echo "  • http (HTTP Requests) - Verde oscuro"
echo "  • file-content (File Content) - Magenta"
echo "  • pe-metadata (PE Metadata) - Naranja"
echo "  • environment (Environment Vars) - Morado"
echo ""
echo "Backup guardado en:"
echo "  $BACKUP_DIR"
echo ""
echo "Para verificar:"
echo "  1. Abre http://tu-servidor/dashboard/"
echo "  2. Ejecuta el agente Phase 1"
echo "  3. Verifica que aparecen los nuevos badges de colores"
echo ""
echo "Para rollback:"
echo "  cp -r $BACKUP_DIR/* $TEMPLATES_DIR/"
echo "  sudo systemctl restart gunicorn"
echo ""
