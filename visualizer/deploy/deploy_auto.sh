#!/bin/bash
# Script de despliegue automático para Artefacto Visualizer
# Ejecutar como root en Ubuntu Server

set -e  # Salir si hay algún error

echo "🚀 Iniciando despliegue de Artefacto Visualizer..."
echo "=================================================="

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Variables
PROJECT_DIR="/opt/artefacto-visualizer"
VENV_DIR="$PROJECT_DIR/venv"
LOG_DIR="/var/log/gunicorn"

# Función para imprimir mensajes
print_step() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

# 1. Actualizar sistema
print_step "Actualizando sistema..."
apt update -qq

# 2. Instalar dependencias del sistema
print_step "Instalando dependencias del sistema..."
apt install -y python3 python3-pip python3-venv nginx > /dev/null 2>&1

# 3. Verificar que el proyecto esté en el directorio correcto
if [ ! -f "$PROJECT_DIR/manage.py" ]; then
    print_warning "No se encuentra manage.py en $PROJECT_DIR"
    echo "Asegúrate de haber descomprimido el proyecto en $PROJECT_DIR"
    exit 1
fi

cd $PROJECT_DIR

# 4. Crear entorno virtual si no existe
if [ ! -d "$VENV_DIR" ]; then
    print_step "Creando entorno virtual..."
    python3 -m venv $VENV_DIR
fi

# 5. Activar entorno virtual e instalar dependencias
print_step "Instalando dependencias de Python..."
source $VENV_DIR/bin/activate
python -m pip install --upgrade pip -q
pip install -r requirements.txt -q

# 6. Configurar base de datos
print_step "Configurando base de datos..."
python manage.py makemigrations --noinput
python manage.py makemigrations collector --noinput
python manage.py makemigrations xss_audit --noinput
python manage.py migrate --noinput

# 7. Recolectar archivos estáticos
print_step "Recolectando archivos estáticos..."
python manage.py collectstatic --noinput

# 8. Configurar permisos
print_step "Configurando permisos..."
chown -R www-data:www-data $PROJECT_DIR
chmod 664 $PROJECT_DIR/db.sqlite3 2>/dev/null || true
chown www-data:www-data $PROJECT_DIR/db.sqlite3 2>/dev/null || true

# 9. Crear directorio de logs de Gunicorn
print_step "Creando directorio de logs..."
mkdir -p $LOG_DIR
chown -R www-data:www-data $LOG_DIR

# 10. Configurar Nginx
print_step "Configurando Nginx..."
cp $PROJECT_DIR/deploy/nginx.conf /etc/nginx/sites-available/artefacto-visualizer
ln -sf /etc/nginx/sites-available/artefacto-visualizer /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Verificar configuración de Nginx
nginx -t

# Reiniciar Nginx
systemctl restart nginx
systemctl enable nginx

# 11. Detener Gunicorn si está corriendo
print_step "Deteniendo Gunicorn anterior..."
pkill -f gunicorn || true
sleep 2

# 12. Iniciar Gunicorn
print_step "Iniciando Gunicorn..."
$VENV_DIR/bin/gunicorn \
  --config $PROJECT_DIR/deploy/gunicorn_config.py \
  visualizer.wsgi:application \
  --daemon

# Esperar a que Gunicorn inicie
sleep 3

# 13. Verificar que Gunicorn esté corriendo
if pgrep -f gunicorn > /dev/null; then
    print_step "Gunicorn iniciado correctamente"
else
    print_warning "Error: Gunicorn no se inició correctamente"
    echo "Revisa los logs en $LOG_DIR/error.log"
    exit 1
fi

# 14. Verificar que Nginx esté corriendo
if systemctl is-active --quiet nginx; then
    print_step "Nginx corriendo correctamente"
else
    print_warning "Error: Nginx no está corriendo"
    exit 1
fi

echo ""
echo "=================================================="
echo -e "${GREEN}✅ Despliegue completado exitosamente!${NC}"
echo "=================================================="
echo ""
echo "🌐 Accede a la aplicación en:"
echo "   https://releases.life/"
echo "   https://releases.life/dashboard/"
echo ""
echo "📝 Logs disponibles en:"
echo "   Gunicorn: $LOG_DIR/error.log"
echo "   Nginx: /var/log/nginx/artefacto-visualizer-error.log"
echo ""
echo "🔄 Para reiniciar Gunicorn:"
echo "   sudo pkill -f gunicorn"
echo "   cd $PROJECT_DIR"
echo "   $VENV_DIR/bin/gunicorn --config deploy/gunicorn_config.py visualizer.wsgi:application --daemon"
echo ""
