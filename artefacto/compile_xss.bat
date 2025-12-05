@echo off
echo ========================================
echo   COMPILANDO AGENTE CON XSS AUDIT
echo ========================================
echo.

REM Verificar que estamos en el directorio correcto
if not exist "main.go" (
    echo ERROR: No se encuentra main.go
    echo Ejecuta este script desde la carpeta artefacto
    pause
    exit /b 1
)

REM Verificar que existe .env
if not exist ".env" (
    echo ADVERTENCIA: No se encuentra .env
    echo Copiando .env.example a .env...
    copy .env.example .env
)

echo Verificando configuracion...
findstr /C:"XSS_AUDIT=true" .env >nul
if errorlevel 1 (
    echo.
    echo ADVERTENCIA: XSS_AUDIT no esta activado en .env
    echo Asegurate de que .env contenga: XSS_AUDIT=true
    echo.
)

echo.
echo Compilando con optimizaciones...
echo.

REM Compilar con optimizaciones y sin ventana
go build -ldflags="-s -w -H windowsgui" -trimpath -o conhost_xss.exe

if errorlevel 1 (
    echo.
    echo ERROR: La compilacion fallo
    pause
    exit /b 1
)

echo.
echo ========================================
echo   COMPILACION EXITOSA
echo ========================================
echo.
echo Binario generado: conhost_xss.exe
echo Tamano: 
for %%A in (conhost_xss.exe) do echo   %%~zA bytes

echo.
echo Para ejecutar:
echo   .\conhost_xss.exe
echo.
echo Para ver la configuracion:
echo   type .env
echo.

pause
