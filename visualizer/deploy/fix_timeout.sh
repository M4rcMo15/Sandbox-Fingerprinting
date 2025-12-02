#!/bin/bash
# Script para aplicar los cambios de timeout en el servidor

echo "🔧 Aplicando cambios de timeout..."

# Copiar configuraciones actualizadas
echo "📋 Actualizando configuraciones..."
sudo cp /opt/artefacto-visualizer/deploy/gunicorn_config.py /opt/artefacto-visualizer/deploy/gunicorn_config.py.bak
sudo cp /opt/artefacto-visualizer/deploy/nginx.conf /etc/nginx/sites-available/artefacto-visualizer

# Verificar Nginx
echo "🔍 Verificando configuración de Nginx..."
sudo nginx -t

if [ $? -eq 0 ]; then
    echo "✅ Configuración de Nginx válida"
    
    # Reiniciar servicios
    echo "🔄 Reiniciando servicios..."
    sudo systemctl restart artefacto-visualizer
    sudo systemctl restart nginx
    
    echo ""
    echo "✅ Cambios aplicados exitosamente!"
    echo ""
    echo "📊 Timeouts configurados:"
    echo "  - Gunicorn: 300 segundos (5 minutos)"
    echo "  - Nginx: 300 segundos (5 minutos)"
    echo ""
    echo "🧪 Prueba ahora el agente:"
    echo "  .\agent.exe"
    echo ""
    echo "📋 Ver logs:"
    echo "  sudo journalctl -u artefacto-visualizer -f"
else
    echo "❌ Error en configuración de Nginx"
    echo "Revirtiendo cambios..."
    sudo cp /opt/artefacto-visualizer/deploy/gunicorn_config.py.bak /opt/artefacto-visualizer/deploy/gunicorn_config.py
    exit 1
fi
