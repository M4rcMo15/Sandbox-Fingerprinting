#!/bin/bash

# Script de despliegue del módulo XSS Audit a producción
# Uso: ./deploy_xss_audit.sh

set -e  # Salir si hay algún error

echo "========================================="
echo "🎯 Despliegue de XSS Audit Module"
echo "========================================="
echo ""

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Directorio del proyecto
PROJECT_DIR="/opt/artefacto-visualizer"
VENV_DIR="/opt/venv"

# Verificar que estamos en el servidor correcto
if [ ! -d "$PROJECT_DIR" ]; then
    echo -e "${RED}❌ Error: Directorio $PROJECT_DIR no encontrado${NC}"
    echo "¿Estás en el servidor correcto?"
    exit 1
fi

echo -e "${YELLOW}📦 Paso 1: Backup de la base de datos${NC}"
BACKUP_DIR="$PROJECT_DIR/../backups"
mkdir -p $BACKUP_DIR
BACKUP_FILE="$BACKUP_DIR/db_backup_$(date +%Y%m%d_%H%M%S).sqlite3"
cp $PROJECT_DIR/db.sqlite3 $BACKUP_FILE
echo -e "${GREEN}✅ Backup creado: $BACKUP_FILE${NC}"
echo ""

echo -e "${YELLOW}📦 Paso 2: Activar entorno virtual${NC}"
source $VENV_DIR/bin/activate
echo -e "${GREEN}✅ Entorno virtual activado${NC}"
echo ""

echo -e "${YELLOW}📦 Paso 3: Aplicar migraciones${NC}"
cd $PROJECT_DIR
python manage.py makemigrations xss_audit
python manage.py migrate
echo -e "${GREEN}✅ Migraciones aplicadas${NC}"
echo ""

echo -e "${YELLOW}📦 Paso 4: Recolectar archivos estáticos${NC}"
python manage.py collectstatic --noinput
echo -e "${GREEN}✅ Archivos estáticos recolectados${NC}"
echo ""

echo -e "${YELLOW}📦 Paso 5: Verificar configuración${NC}"
# Verificar que xss_audit está en INSTALLED_APPS
if grep -q "xss_audit" $PROJECT_DIR/visualizer/settings.py; then
    echo -e "${GREEN}✅ xss_audit está en INSTALLED_APPS${NC}"
else
    echo -e "${RED}❌ Error: xss_audit NO está en INSTALLED_APPS${NC}"
    echo "Añadir 'xss_audit' a INSTALLED_APPS en settings.py"
    exit 1
fi
echo ""

echo -e "${YELLOW}📦 Paso 6: Reiniciar servicios${NC}"
sudo systemctl restart gunicorn
echo -e "${GREEN}✅ Gunicorn reiniciado${NC}"

sudo systemctl restart nginx
echo -e "${GREEN}✅ Nginx reiniciado${NC}"
echo ""

echo -e "${YELLOW}📦 Paso 7: Verificar estado de servicios${NC}"
if systemctl is-active --quiet gunicorn; then
    echo -e "${GREEN}✅ Gunicorn está activo${NC}"
else
    echo -e "${RED}❌ Error: Gunicorn no está activo${NC}"
    sudo systemctl status gunicorn
    exit 1
fi

if systemctl is-active --quiet nginx; then
    echo -e "${GREEN}✅ Nginx está activo${NC}"
else
    echo -e "${RED}❌ Error: Nginx no está activo${NC}"
    sudo systemctl status nginx
    exit 1
fi
echo ""

echo -e "${YELLOW}📦 Paso 8: Probar endpoint de callback${NC}"
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost/xss-callback?id=test123&v=test")
if [ "$RESPONSE" = "200" ]; then
    echo -e "${GREEN}✅ Endpoint /xss-callback responde correctamente (HTTP $RESPONSE)${NC}"
else
    echo -e "${RED}❌ Error: Endpoint /xss-callback no responde correctamente (HTTP $RESPONSE)${NC}"
    exit 1
fi
echo ""

echo -e "${YELLOW}📦 Paso 9: Verificar dashboard${NC}"
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost/xss-audit/dashboard/")
if [ "$RESPONSE" = "200" ]; then
    echo -e "${GREEN}✅ Dashboard XSS Audit accesible (HTTP $RESPONSE)${NC}"
else
    echo -e "${RED}❌ Error: Dashboard no accesible (HTTP $RESPONSE)${NC}"
    exit 1
fi
echo ""

echo "========================================="
echo -e "${GREEN}✅ Despliegue completado exitosamente${NC}"
echo "========================================="
echo ""
echo "📊 Acceder al dashboard:"
echo "   https://releases.life/xss-audit/dashboard/"
echo ""
echo "🔍 Ver logs:"
echo "   sudo journalctl -u gunicorn -n 50 -f"
echo "   sudo tail -f /var/log/nginx/artefacto-visualizer-error.log"
echo ""
echo "📝 Backup de la base de datos:"
echo "   $BACKUP_FILE"
echo ""
