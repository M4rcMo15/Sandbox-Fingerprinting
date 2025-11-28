#!/usr/bin/env python
"""
Script de configuración inicial para Artefacto Visualizer
"""
import os
import sys
import subprocess

def run_command(command, description):
    """Ejecuta un comando y muestra el resultado"""
    print(f"\n{'='*60}")
    print(f"  {description}")
    print('='*60)
    try:
        result = subprocess.run(command, shell=True, check=True, 
                              capture_output=True, text=True)
        if result.stdout:
            print(result.stdout)
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error: {e}")
        if e.stderr:
            print(e.stderr)
        return False

def main():
    print("""
    ╔═══════════════════════════════════════════════════════════╗
    ║         ARTEFACTO VISUALIZER - Setup Inicial              ║
    ╚═══════════════════════════════════════════════════════════╝
    """)
    
    # Verificar Python
    print(f"Python version: {sys.version}")
    
    # Instalar dependencias
    if not run_command("pip install -r requirements.txt", 
                      "Instalando dependencias de Python"):
        print("\n⚠️  Error instalando dependencias. Verifica que pip esté instalado.")
        return
    
    # Crear migraciones
    if not run_command("python manage.py makemigrations", 
                      "Creando migraciones de base de datos"):
        print("\n⚠️  Error creando migraciones.")
        return
    
    # Aplicar migraciones
    if not run_command("python manage.py migrate", 
                      "Aplicando migraciones"):
        print("\n⚠️  Error aplicando migraciones.")
        return
    
    print("""
    ╔═══════════════════════════════════════════════════════════╗
    ║                  ✅ Setup Completado                      ║
    ╚═══════════════════════════════════════════════════════════╝
    
    Para iniciar el servidor:
    
    Windows:  start_server.bat
    Linux:    ./start_server.sh
    Manual:   python manage.py runserver 192.168.1.143:8080
    
    Accede a: http://192.168.1.143:8080/
    """)

if __name__ == '__main__':
    main()
