#!/bin/bash
# Script para configurar autenticación en el visualizer

echo "🔐 Configurando autenticación para el visualizer..."

# Instalar apache2-utils para htpasswd
echo "📦 Instalando herramientas necesarias..."
sudo apt-get update
sudo apt-get install -y apache2-utils

# Crear directorio para archivos de autenticación
echo "📁 Creando directorio de autenticación..."
sudo mkdir -p /etc/nginx/auth

# Crear usuario marc.monfort
echo "👤 Creando usuario marc.monfort..."
echo "Introduce la contraseña para marc.monfort:"
sudo htpasswd -c /etc/nginx/auth/.htpasswd marc.monfort

# Configurar permisos
sudo chmod 640 /etc/nginx/auth/.htpasswd
sudo chown www-data:www-data /etc/nginx/auth/.htpasswd

echo ""
echo "✅ Autenticación configurada!"
echo ""
echo "📝 Ahora actualiza la configuración de Nginx:"
echo "   sudo nano /etc/nginx/sites-available/artefacto-visualizer"
echo ""
echo "Agrega estas líneas dentro del bloque 'location /':"
echo "   auth_basic \"Artefacto Visualizer - Acceso Restringido\";"
echo "   auth_basic_user_file /etc/nginx/auth/.htpasswd;"
echo ""
echo "O ejecuta: sudo ./deploy/update_nginx_auth.sh"
echo ""
echo "Luego reinicia Nginx:"
echo "   sudo nginx -t"
echo "   sudo systemctl restart nginx"
