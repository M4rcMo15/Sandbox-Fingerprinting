#!/usr/bin/env python
"""
Script rápido para probar el endpoint
"""
import json

# Datos mínimos de prueba
test_data = {
    "timestamp": "2024-11-28T16:12:09Z",
    "hostname": "TEST-MACHINE"
}

print("="*60)
print("  TEST ENDPOINT - Verificación Rápida")
print("="*60)
print("\nProbando URLs:")
print("1. http://192.168.1.143:8080/api/collect")
print("2. http://192.168.1.143:8080/api/collect/")
print("\nDatos de prueba:")
print(json.dumps(test_data, indent=2))
print("\n" + "="*60)

try:
    import requests
    
    # Probar sin barra
    print("\n[1] Probando sin barra final...")
    try:
        r1 = requests.post(
            "http://192.168.1.143:8080/api/collect",
            json=test_data,
            timeout=5
        )
        print(f"    Status: {r1.status_code}")
        print(f"    Response: {r1.text[:100]}")
    except Exception as e:
        print(f"    Error: {e}")
    
    # Probar con barra
    print("\n[2] Probando con barra final...")
    try:
        r2 = requests.post(
            "http://192.168.1.143:8080/api/collect/",
            json=test_data,
            timeout=5
        )
        print(f"    Status: {r2.status_code}")
        print(f"    Response: {r2.text[:100]}")
    except Exception as e:
        print(f"    Error: {e}")
    
    print("\n" + "="*60)
    print("✅ Si ves status 201 o 200, el endpoint funciona!")
    print("="*60)
    
except ImportError:
    print("\n⚠️  requests no está instalado")
    print("Instala con: pip install requests")
