#!/bin/bash
# Script de despliegue rápido para Artefacto Visualizer
# Ejecutar como root en el servidor: bash quick_deploy.sh

set -e

echo "🚀 Despliegue Rápido de Artefacto Visualizer"
echo "=============================================="
echo ""

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Variables
PROJECT_DIR="/opt/artefacto-visualizer"
VENV_DIR="$PROJECT_DIR/venv"
LOG_DIR="/var/log/gunicorn"

# Verificar que se ejecuta como root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}[✗]${NC} Este script debe ejecutarse como root"
    echo "Usa: sudo bash quick_deploy.sh"
    exit 1
fi

echo -e "${GREEN}[✓]${NC} Ejecutando como root"

# Verificar que el proyecto existe
if [ ! -f "$PROJECT_DIR/manage.py" ]; then
    echo -e "${RED}[✗]${NC} No se encuentra el proyecto en $PROJECT_DIR"
    echo "Asegúrate de haber transferido y descomprimido el proyecto primero"
    exit 1
fi

echo -e "${GREEN}[✓]${NC} Proyecto encontrado en $PROJECT_DIR"

# Actualizar sistema
echo -e "${YELLOW}[→]${NC} Actualizando sistema..."
apt update -qq

# Instalar dependencias
echo -e "${YELLOW}[→]${NC} Instalando dependencias..."
apt install -y python3 python3-pip python3-venv nginx certbot python3-certbot-nginx > /dev/null 2>&1

# Crear entorno virtual
if [ ! -d "$VENV_DIR" ]; then
    echo -e "${YELLOW}[→]${NC} Creando entorno virtual..."
    python3 -m venv $VENV_DIR
fi

# Instalar dependencias Python
echo -e "${YELLOW}[→]${NC} Instalando dependencias Python..."
cd $PROJECT_DIR
source $VENV_DIR/bin/activate
pip install --upgrade pip -q
pip install -r requirements.txt -q

# Aplicar migraciones
echo -e "${YELLOW}[→]${NC} Aplicando migraciones..."
python manage.py makemigrations --noinput
python manage.py makemigrations collector --noinput
python manage.py makemigrations xss_audit --noinput
python manage.py migrate --noinput

# Recolectar estáticos
echo -e "${YELLOW}[→]${NC} Recolectando archivos estáticos..."
python manage.py collectstatic --noinput

# Configurar permisos
echo -e "${YELLOW}[→]${NC} Configurando permisos..."
chown -R www-data:www-data $PROJECT_DIR
chmod 664 $PROJECT_DIR/db.sqlite3 2>/dev/null || true

# Crear directorios de logs
echo -e "${YELLOW}[→]${NC} Creando directorios de logs..."
mkdir -p $LOG_DIR
mkdir -p /var/run/gunicorn
chown -R www-data:www-data $LOG_DIR
chown -R www-data:www-data /var/run/gunicorn

# Configurar systemd service
echo -e "${YELLOW}[→]${NC} Configurando servicio systemd..."
cat > /etc/systemd/system/artefacto-visualizer.service << 'EOF'
[Unit]
Description=Artefacto Visualizer Gunicorn daemon
After=network.target

[Service]
Type=notify
User=www-data
Group=www-data
RuntimeDirectory=gunicorn
WorkingDirectory=/opt/artefacto-visualizer
Environment="PATH=/opt/artefacto-visualizer/venv/bin"
ExecStart=/opt/artefacto-visualizer/venv/bin/gunicorn \
          --config /opt/artefacto-visualizer/deploy/gunicorn_config.py \
          visualizer.wsgi:application
ExecReload=/bin/kill -s HUP $MAINPID
KillMode=mixed
TimeoutStopSec=5
PrivateTmp=true
Restart=always

[Install]
WantedBy=multi-user.target
EOF

# Recargar systemd
systemctl daemon-reload

# Configurar Nginx
echo -e "${YELLOW}[→]${NC} Configurando Nginx..."
cp $PROJECT_DIR/deploy/nginx.conf /etc/nginx/sites-available/artefacto-visualizer
ln -sf /etc/nginx/sites-available/artefacto-visualizer /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Verificar configuración de Nginx
if nginx -t 2>/dev/null; then
    echo -e "${GREEN}[✓]${NC} Configuración de Nginx válida"
else
    echo -e "${RED}[✗]${NC} Error en configuración de Nginx"
    nginx -t
    exit 1
fi

# Reiniciar servicios
echo -e "${YELLOW}[→]${NC} Reiniciando servicios..."
systemctl enable artefacto-visualizer
systemctl restart artefacto-visualizer
systemctl enable nginx
systemctl restart nginx

# Esperar a que los servicios inicien
sleep 3

# Verificar servicios
echo ""
echo "Verificando servicios..."
if systemctl is-active --quiet artefacto-visualizer; then
    echo -e "${GREEN}[✓]${NC} Gunicorn corriendo"
else
    echo -e "${RED}[✗]${NC} Gunicorn no está corriendo"
    echo "Ver logs: journalctl -u artefacto-visualizer -n 50"
fi

if systemctl is-active --quiet nginx; then
    echo -e "${GREEN}[✓]${NC} Nginx corriendo"
else
    echo -e "${RED}[✗]${NC} Nginx no está corriendo"
fi

echo ""
echo "=============================================="
echo -e "${GREEN}✅ Despliegue completado!${NC}"
echo "=============================================="
echo ""
echo "🌐 Accede a la aplicación en:"
echo "   http://releases.life/"
echo "   http://79.137.76.235/"
echo ""
echo "🔒 Para configurar HTTPS, ejecuta:"
echo "   certbot --nginx -d releases.life -d www.releases.life"
echo ""
echo "📝 Logs disponibles en:"
echo "   Gunicorn: $LOG_DIR/error.log"
echo "   Nginx: /var/log/nginx/artefacto-error.log"
echo "   Systemd: journalctl -u artefacto-visualizer -f"
echo ""
echo "🔄 Comandos útiles:"
echo "   Reiniciar: systemctl restart artefacto-visualizer"
echo "   Ver logs: tail -f $LOG_DIR/error.log"
echo "   Estado: systemctl status artefacto-visualizer"
echo ""
