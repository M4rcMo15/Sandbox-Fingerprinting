#!/bin/bash

clear
cat banner.txt
echo ""

# Verificar si existe la base de datos
if [ ! -f db.sqlite3 ]; then
    echo "[*] Creando base de datos..."
    python manage.py makemigrations
    python manage.py migrate
    echo ""
fi

echo "[*] Iniciando servidor en 192.168.1.143:8080..."
echo "[*] Presiona Ctrl+C para detener el servidor"
echo ""
echo "================================================================"
echo "  Accede a: http://192.168.1.143:8080/"
echo "================================================================"
echo ""
python manage.py runserver 192.168.1.143:8080
