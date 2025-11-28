@echo off
title Artefacto Visualizer - Quick Start
color 0A

echo.
echo ========================================================
echo   ARTEFACTO VISUALIZER - Quick Start
echo ========================================================
echo.

:menu
echo Selecciona una opcion:
echo.
echo [1] Instalar dependencias y configurar
echo [2] Iniciar servidor
echo [3] Probar API
echo [4] Crear superusuario (admin)
echo [5] Salir
echo.
set /p option="Opcion: "

if "%option%"=="1" goto install
if "%option%"=="2" goto start
if "%option%"=="3" goto test
if "%option%"=="4" goto superuser
if "%option%"=="5" goto end
goto menu

:install
echo.
echo [*] Instalando dependencias...
pip install -r requirements.txt
echo.
echo [*] Creando base de datos...
python manage.py makemigrations
python manage.py migrate
echo.
echo [OK] Instalacion completada!
echo.
pause
goto menu

:start
echo.
echo [*] Iniciando servidor en 192.168.1.143:8080...
echo [*] Presiona Ctrl+C para detener
echo.
python manage.py runserver 192.168.1.143:8080
pause
goto menu

:test
echo.
echo [*] Probando API...
python test_api.py
echo.
pause
goto menu

:superuser
echo.
echo [*] Crear superusuario para el panel de administracion
echo.
python manage.py createsuperuser
echo.
pause
goto menu

:end
echo.
echo Hasta luego!
exit
