#!/usr/bin/env python
"""
Script para probar la nueva página de estadísticas
"""
import requests
import sys

def test_statistics_page():
    """Prueba que la página de estadísticas esté accesible"""
    base_url = "http://127.0.0.1:8000"
    
    print("🧪 Probando página de estadísticas...")
    
    try:
        # Probar página principal
        response = requests.get(f"{base_url}/")
        if response.status_code == 200:
            print("✅ Página principal accesible")
        else:
            print(f"❌ Error en página principal: {response.status_code}")
            return False
        
        # Probar página de estadísticas
        response = requests.get(f"{base_url}/statistics/")
        if response.status_code == 200:
            print("✅ Página de estadísticas accesible")
            
            # Verificar que contiene elementos esperados
            content = response.text
            if "Estadísticas Generales" in content:
                print("✅ Contenido de estadísticas presente")
            if "Chart.js" in content or "chart.js" in content:
                print("✅ Librería de gráficos cargada")
            
            return True
        else:
            print(f"❌ Error en página de estadísticas: {response.status_code}")
            return False
            
    except requests.exceptions.ConnectionError:
        print("❌ No se puede conectar al servidor")
        print("   Asegúrate de que el servidor Django esté corriendo en http://127.0.0.1:8000")
        return False
    except Exception as e:
        print(f"❌ Error inesperado: {e}")
        return False

if __name__ == "__main__":
    success = test_statistics_page()
    sys.exit(0 if success else 1)
