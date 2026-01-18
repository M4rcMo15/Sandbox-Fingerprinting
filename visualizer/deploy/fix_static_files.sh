#!/bin/bash
# Script para solucionar problemas con archivos estáticos

echo "🔧 Solucionando problemas con archivos estáticos..."
echo ""

cd /opt/artefacto-visualizer

# 1. Verificar permisos
echo "1. Verificando permisos..."
chown -R www-data:www-data /opt/artefacto-visualizer
chmod -R 755 /opt/artefacto-visualizer
chmod 664 /opt/artefacto-visualizer/db.sqlite3 2>/dev/null || true

# 2. Recolectar archivos estáticos
echo "2. Recolectando archivos estáticos..."
source venv/bin/activate
python manage.py collectstatic --noinput --clear

# 3. Verificar directorio staticfiles
echo "3. Verificando directorio staticfiles..."
if [ -d "staticfiles" ]; then
    echo "   ✓ Directorio staticfiles existe"
    echo "   Archivos: $(find staticfiles -type f | wc -l)"
    chown -R www-data:www-data staticfiles/
    chmod -R 755 staticfiles/
else
    echo "   ✗ Directorio staticfiles NO existe"
fi

# 4. Verificar configuración de Nginx
echo "4. Verificando configuración de Nginx..."
nginx -t

# 5. Reiniciar servicios
echo "5. Reiniciando servicios..."
systemctl restart artefacto-visualizer
systemctl restart nginx

echo ""
echo "✅ Proceso completado"
echo ""
echo "Verifica en el navegador:"
echo "  - Abre https://releases.life/"
echo "  - Presiona Ctrl+Shift+R (hard refresh)"
echo "  - Abre la consola (F12) y busca errores"
echo ""
