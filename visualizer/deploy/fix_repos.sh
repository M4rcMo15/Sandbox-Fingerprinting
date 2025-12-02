#!/bin/bash
# Script para limpiar repositorios problemáticos

echo "🔧 Limpiando repositorios problemáticos..."

# Deshabilitar repositorio de MongoDB
if [ -f /etc/apt/sources.list.d/mongodb-org-6.0.list ]; then
    echo "📦 Deshabilitando repositorio de MongoDB..."
    sudo mv /etc/apt/sources.list.d/mongodb-org-6.0.list /etc/apt/sources.list.d/mongodb-org-6.0.list.disabled
fi

# Deshabilitar repositorio de Graylog
if [ -f /etc/apt/sources.list.d/graylog.list ]; then
    echo "📦 Deshabilitando repositorio de Graylog..."
    sudo mv /etc/apt/sources.list.d/graylog.list /etc/apt/sources.list.d/graylog.list.disabled
fi

# Actualizar lista de paquetes
echo "🔄 Actualizando lista de paquetes..."
sudo apt update

echo "✅ Repositorios limpiados!"
