#!/bin/bash
# Script para verificar el estado del sistema

echo "🏥 Health Check - Artefacto Visualizer"
echo "======================================"
echo ""

# Verificar servicio
echo "📊 Estado del servicio:"
sudo systemctl is-active artefacto-visualizer && echo "✅ Servicio activo" || echo "❌ Servicio inactivo"
echo ""

# Verificar Nginx
echo "🌐 Estado de Nginx:"
sudo systemctl is-active nginx && echo "✅ Nginx activo" || echo "❌ Nginx inactivo"
echo ""

# Verificar puerto
echo "🔌 Puerto 8080:"
sudo netstat -tlnp | grep :8080 && echo "✅ Puerto en uso" || echo "❌ Puerto libre"
echo ""

# Verificar base de datos
echo "🗄️ Base de datos:"
if [ -f "/opt/artefacto-visualizer/db.sqlite3" ]; then
    SIZE=$(du -h /opt/artefacto-visualizer/db.sqlite3 | cut -f1)
    echo "✅ Base de datos existe ($SIZE)"
else
    echo "❌ Base de datos no encontrada"
fi
echo ""

# Verificar logs
echo "📋 Últimas 5 líneas de log:"
sudo journalctl -u artefacto-visualizer -n 5 --no-pager
echo ""

# Test HTTP
echo "🌐 Test HTTP:"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost)
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ HTTP responde correctamente (200)"
else
    echo "⚠️ HTTP responde con código: $HTTP_CODE"
fi
echo ""

# Espacio en disco
echo "💾 Espacio en disco:"
df -h /opt/artefacto-visualizer | tail -1
echo ""

# Memoria
echo "🧠 Uso de memoria:"
free -h | grep Mem
echo ""

echo "======================================"
echo "✅ Health check completado"
