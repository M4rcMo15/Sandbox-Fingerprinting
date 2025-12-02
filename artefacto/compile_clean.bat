@echo off
REM Script para compilar el agente limpio para producción
REM Sin referencias a IPs de desarrollo

echo ========================================
echo   Compilando Agente LIMPIO
echo   Servidor: 54.37.226.179
echo ========================================
echo.

REM Limpiar compilaciones anteriores
echo [1/4] Limpiando compilaciones anteriores...
if exist agent.exe del agent.exe
if exist agent_*.exe del agent_*.exe
if exist payload_*.json del payload_*.json

REM Verificar .env
echo.
echo [2/4] Verificando configuración...
if not exist ".env" (
    echo ADVERTENCIA: No existe .env, copiando desde .env.example
    copy .env.example .env
)

type .env | findstr "54.37.226.179"
if errorlevel 1 (
    echo.
    echo ERROR: El .env no apunta al servidor de producción
    echo Debe contener: SERVER_URL=http://54.37.226.179/api/collect
    echo.
    pause
    exit /b 1
)

REM Compilar
echo.
echo [3/4] Compilando agente...
go build -o agent.exe -ldflags="-s -w"
if errorlevel 1 (
    echo ERROR: Fallo la compilación
    pause
    exit /b 1
)

REM Verificar
echo.
echo [4/4] Verificación...
if exist "agent.exe" (
    echo OK - Agente compilado exitosamente
    for %%A in (agent.exe) do echo Tamaño: %%~zA bytes
    echo.
    echo ========================================
    echo   COMPILACION LIMPIA COMPLETADA
    echo ========================================
    echo.
    echo El agente está listo para:
    echo   - Testing en sandboxes
    echo   - Subir a VirusTotal
    echo   - Subir a Hybrid Analysis
    echo   - Subir a Any.Run
    echo.
    echo Servidor configurado:
    echo   http://54.37.226.179/api/collect
    echo.
    echo IMPORTANTE:
    echo   - Sin referencias a IPs de desarrollo
    echo   - Listo para análisis público
    echo   - Timeout: 120 segundos
    echo.
) else (
    echo ERROR: No se generó el ejecutable
    exit /b 1
)

pause
