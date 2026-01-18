#!/bin/bash
# Script de despliegue automático para Ubuntu Server 24

set -e  # Salir si hay errores

echo "🚀 Iniciando despliegue de Artefacto Visualizer..."

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Variables
APP_DIR="/opt/artefacto-visualizer"
VENV_DIR="$APP_DIR/venv"
LOG_DIR="/var/log/artefacto-visualizer"
RUN_DIR="/var/run/artefacto-visualizer"

echo -e "${YELLOW}📦 Actualizando sistema...${NC}"
sudo apt update
sudo apt upgrade -y

echo -e "${YELLOW}📦 Instalando dependencias del sistema...${NC}"
sudo apt install -y \
    python3 \
    python3-pip \
    python3-venv \
    nginx \
    git \
    curl \
    build-essential \
    libpq-dev

echo -e "${YELLOW}📁 Creando directorios...${NC}"
sudo mkdir -p $APP_DIR
sudo mkdir -p $LOG_DIR
sudo mkdir -p $RUN_DIR
sudo mkdir -p $APP_DIR/staticfiles
sudo mkdir -p $APP_DIR/media

echo -e "${YELLOW}📋 Copiando archivos de la aplicación...${NC}"
# Asumiendo que estás en el directorio visualizer
sudo cp -r . $APP_DIR/
sudo chown -R www-data:www-data $APP_DIR
sudo chown -R www-data:www-data $LOG_DIR
sudo chown -R www-data:www-data $RUN_DIR

echo -e "${YELLOW}🐍 Creando entorno virtual...${NC}"
sudo -u www-data python3 -m venv $VENV_DIR

echo -e "${YELLOW}📦 Instalando dependencias de Python...${NC}"
sudo -u www-data $VENV_DIR/bin/pip install --upgrade pip
sudo -u www-data $VENV_DIR/bin/pip install -r $APP_DIR/deploy/requirements_production.txt

echo -e "${YELLOW}🔑 Generando SECRET_KEY...${NC}"
SECRET_KEY=$(python3 -c 'from django.core.management.utils import get_random_secret_key; print(get_random_secret_key())')

echo -e "${YELLOW}⚙️ Configurando variables de entorno...${NC}"
sudo cp $APP_DIR/deploy/.env.production $APP_DIR/.env
sudo sed -i "s/CHANGE_THIS_TO_A_RANDOM_SECRET_KEY_MINIMUM_50_CHARACTERS/$SECRET_KEY/" $APP_DIR/.env
sudo chown www-data:www-data $APP_DIR/.env
sudo chmod 600 $APP_DIR/.env

echo -e "${YELLOW}🗄️ Configurando base de datos...${NC}"
cd $APP_DIR
sudo -u www-data $VENV_DIR/bin/python manage.py migrate
sudo -u www-data $VENV_DIR/bin/python manage.py collectstatic --noinput

echo -e "${YELLOW}👤 Creando superusuario (opcional)...${NC}"
echo "¿Deseas crear un superusuario para el admin de Django? (s/n)"
read -r CREATE_SUPERUSER
if [ "$CREATE_SUPERUSER" = "s" ]; then
    sudo -u www-data $VENV_DIR/bin/python manage.py createsuperuser
fi

echo -e "${YELLOW}🔧 Configurando systemd service...${NC}"
sudo cp $APP_DIR/deploy/artefacto-visualizer.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable artefacto-visualizer
sudo systemctl start artefacto-visualizer

echo -e "${YELLOW}🌐 Configurando Nginx...${NC}"
sudo cp $APP_DIR/deploy/nginx.conf /etc/nginx/sites-available/artefacto-visualizer
sudo ln -sf /etc/nginx/sites-available/artefacto-visualizer /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl restart nginx

echo -e "${YELLOW}🔥 Configurando firewall...${NC}"
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw --force enable

echo -e "${GREEN}✅ Despliegue completado!${NC}"
echo ""
echo -e "${GREEN}🎉 El visualizer está corriendo en:${NC}"
echo -e "${GREEN}   https://releases.life${NC}"
echo ""
echo -e "${YELLOW}📊 Comandos útiles:${NC}"
echo "  - Ver logs: sudo journalctl -u artefacto-visualizer -f"
echo "  - Reiniciar: sudo systemctl restart artefacto-visualizer"
echo "  - Estado: sudo systemctl status artefacto-visualizer"
echo "  - Nginx logs: sudo tail -f /var/log/nginx/artefacto-error.log"
echo ""
echo -e "${YELLOW}🔐 Recuerda:${NC}"
echo "  - Configurar SSL/HTTPS con Let's Encrypt"
echo "  - Cambiar credenciales por defecto"
echo "  - Configurar backups de la base de datos"
