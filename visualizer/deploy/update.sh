#!/bin/bash
# Script para actualizar la aplicación en producción

set -e

echo "🔄 Actualizando Artefacto Visualizer..."

APP_DIR="/opt/artefacto-visualizer"
VENV_DIR="$APP_DIR/venv"

# Detener el servicio
echo "⏸️ Deteniendo servicio..."
sudo systemctl stop artefacto-visualizer

# Backup de la base de datos
echo "💾 Haciendo backup de la base de datos..."
BACKUP_DIR="$APP_DIR/backups"
sudo mkdir -p $BACKUP_DIR
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
sudo cp $APP_DIR/db.sqlite3 $BACKUP_DIR/db_backup_$TIMESTAMP.sqlite3

# Actualizar código (si usas git)
# cd $APP_DIR
# sudo -u www-data git pull origin main

# Actualizar dependencias
echo "📦 Actualizando dependencias..."
sudo -u www-data $VENV_DIR/bin/pip install -r $APP_DIR/deploy/requirements_production.txt

# Migraciones
echo "🗄️ Aplicando migraciones..."
cd $APP_DIR
sudo -u www-data $VENV_DIR/bin/python manage.py migrate

# Recolectar archivos estáticos
echo "📁 Recolectando archivos estáticos..."
sudo -u www-data $VENV_DIR/bin/python manage.py collectstatic --noinput

# Reiniciar servicio
echo "▶️ Reiniciando servicio..."
sudo systemctl start artefacto-visualizer
sudo systemctl restart nginx

echo "✅ Actualización completada!"
echo "📊 Estado del servicio:"
sudo systemctl status artefacto-visualizer --no-pager
