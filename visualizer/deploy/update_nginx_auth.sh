#!/bin/bash
# Script para actualizar Nginx con autenticación

echo "🔐 Actualizando configuración de Nginx con autenticación..."

# Verificar que existe el archivo de contraseñas
if [ ! -f /etc/nginx/auth/.htpasswd ]; then
    echo "❌ Error: No existe /etc/nginx/auth/.htpasswd"
    echo "Ejecuta primero: sudo ./deploy/setup_auth.sh"
    exit 1
fi

# Backup de la configuración actual
echo "💾 Haciendo backup de la configuración actual..."
sudo cp /etc/nginx/sites-available/artefacto-visualizer /etc/nginx/sites-available/artefacto-visualizer.backup

# Copiar nueva configuración con autenticación
echo "📋 Aplicando nueva configuración..."
sudo cp /opt/artefacto-visualizer/deploy/nginx_auth.conf /etc/nginx/sites-available/artefacto-visualizer

# Verificar configuración
echo "🔍 Verificando configuración de Nginx..."
sudo nginx -t

if [ $? -eq 0 ]; then
    echo "✅ Configuración válida"
    
    # Reiniciar Nginx
    echo "🔄 Reiniciando Nginx..."
    sudo systemctl restart nginx
    
    echo ""
    echo "✅ Autenticación activada!"
    echo ""
    echo "📊 Acceso al visualizer:"
    echo "   URL: http://54.37.226.179"
    echo "   Usuario: marc.monfort"
    echo "   Contraseña: [la que configuraste]"
    echo ""
    echo "🔓 El endpoint API (/api/collect) NO requiere autenticación"
    echo "   para que el agente pueda enviar datos"
    echo ""
else
    echo "❌ Error en la configuración de Nginx"
    echo "Revirtiendo cambios..."
    sudo cp /etc/nginx/sites-available/artefacto-visualizer.backup /etc/nginx/sites-available/artefacto-visualizer
    exit 1
fi
