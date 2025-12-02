#!/bin/bash
# Script rápido de despliegue (ejecutar en el servidor)

set -e

echo "🚀 Despliegue rápido de Artefacto Visualizer"
echo "============================================"

# Extraer archivos
echo "📦 Extrayendo archivos..."
cd /tmp
tar -xzf visualizer-deploy.tar.gz
cd visualizer-deploy

# Instalar dependencias básicas
echo "📦 Instalando dependencias..."
sudo apt update
sudo apt install -y python3 python3-pip python3-venv nginx

# Crear directorios
echo "📁 Creando estructura..."
sudo mkdir -p /opt/artefacto-visualizer
sudo mkdir -p /var/log/artefacto-visualizer
sudo mkdir -p /var/run/artefacto-visualizer

# Copiar archivos
echo "📋 Copiando aplicación..."
sudo cp -r * /opt/artefacto-visualizer/
sudo chown -R www-data:www-data /opt/artefacto-visualizer
sudo chown -R www-data:www-data /var/log/artefacto-visualizer

# Crear entorno virtual
echo "🐍 Configurando Python..."
cd /opt/artefacto-visualizer
sudo -u www-data python3 -m venv venv
sudo -u www-data venv/bin/pip install -r deploy/requirements_production.txt

# Generar SECRET_KEY
echo "🔑 Generando SECRET_KEY..."
SECRET_KEY=$(python3 -c 'from django.core.management.utils import get_random_secret_key; print(get_random_secret_key())')

# Configurar .env
echo "⚙️ Configurando variables..."
sudo cp deploy/.env.production .env
sudo sed -i "s/CHANGE_THIS_TO_A_RANDOM_SECRET_KEY_MINIMUM_50_CHARACTERS/$SECRET_KEY/" .env
sudo chown www-data:www-data .env

# Migraciones
echo "🗄️ Configurando base de datos..."
sudo -u www-data venv/bin/python manage.py migrate
sudo -u www-data venv/bin/python manage.py collectstatic --noinput

# Configurar systemd
echo "🔧 Configurando servicio..."
sudo cp deploy/artefacto-visualizer.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable artefacto-visualizer
sudo systemctl start artefacto-visualizer

# Configurar Nginx
echo "🌐 Configurando Nginx..."
sudo cp deploy/nginx.conf /etc/nginx/sites-available/artefacto-visualizer
sudo ln -sf /etc/nginx/sites-available/artefacto-visualizer /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t && sudo systemctl restart nginx

# Firewall
echo "🔥 Configurando firewall..."
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
echo "y" | sudo ufw enable

echo ""
echo "✅ ¡Despliegue completado!"
echo ""
echo "🌐 Accede a: http://54.37.226.179"
echo ""
echo "📊 Ver estado: sudo systemctl status artefacto-visualizer"
echo "📋 Ver logs: sudo journalctl -u artefacto-visualizer -f"
