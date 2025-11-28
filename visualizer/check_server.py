#!/usr/bin/env python
"""
Script para verificar que el servidor Django está funcionando correctamente
"""

def check_server():
    print("="*70)
    print("  VERIFICACIÓN DEL SERVIDOR DJANGO")
    print("="*70)
    
    try:
        import requests
    except ImportError:
        print("\n⚠️  El módulo 'requests' no está instalado")
        print("Instala con: pip install requests")
        return
    
    base_url = "http://192.168.1.143:8080"
    
    # 1. Verificar página principal
    print("\n[1] Verificando página principal...")
    try:
        r = requests.get(f"{base_url}/", timeout=5)
        if r.status_code == 200:
            print(f"    ✅ OK - Status {r.status_code}")
        else:
            print(f"    ⚠️  Status {r.status_code}")
    except requests.exceptions.ConnectionError:
        print("    ❌ ERROR - No se puede conectar al servidor")
        print("    ¿Está el servidor corriendo?")
        print("    Ejecuta: python manage.py runserver 192.168.1.143:8080")
        return
    except Exception as e:
        print(f"    ❌ ERROR - {e}")
        return
    
    # 2. Verificar endpoint API (GET debería dar 405)
    print("\n[2] Verificando endpoint API (GET)...")
    try:
        r = requests.get(f"{base_url}/api/collect", timeout=5)
        if r.status_code == 405:
            print(f"    ✅ OK - Status {r.status_code} (Method Not Allowed es correcto)")
        else:
            print(f"    ⚠️  Status {r.status_code}")
    except Exception as e:
        print(f"    ❌ ERROR - {e}")
    
    # 3. Verificar endpoint API (POST con datos mínimos)
    print("\n[3] Verificando endpoint API (POST)...")
    test_data = {
        "timestamp": "2024-11-28T16:12:09Z",
        "hostname": "TEST-CHECK"
    }
    try:
        r = requests.post(
            f"{base_url}/api/collect",
            json=test_data,
            timeout=5
        )
        if r.status_code in [200, 201]:
            print(f"    ✅ OK - Status {r.status_code}")
            try:
                response_data = r.json()
                print(f"    GUID generado: {response_data.get('guid')}")
            except:
                pass
        else:
            print(f"    ⚠️  Status {r.status_code}")
            print(f"    Response: {r.text[:200]}")
    except Exception as e:
        print(f"    ❌ ERROR - {e}")
    
    # 4. Verificar con barra final
    print("\n[4] Verificando endpoint API con barra final (POST)...")
    try:
        r = requests.post(
            f"{base_url}/api/collect/",
            json=test_data,
            timeout=5
        )
        if r.status_code in [200, 201]:
            print(f"    ✅ OK - Status {r.status_code}")
        else:
            print(f"    ⚠️  Status {r.status_code}")
    except Exception as e:
        print(f"    ❌ ERROR - {e}")
    
    print("\n" + "="*70)
    print("  RESUMEN")
    print("="*70)
    print("\nSi todos los checks muestran ✅, el servidor está funcionando correctamente.")
    print("\nPuedes ejecutar el agente con:")
    print("  cd artefacto")
    print("  ./agent.exe")
    print("\nY ver los resultados en:")
    print(f"  {base_url}/")
    print("="*70)

if __name__ == '__main__':
    check_server()
