#!/usr/bin/env python
"""
Script de prueba para verificar que el endpoint API funciona correctamente
"""
import requests
import json
from datetime import datetime

# Datos de prueba simulando el payload del agente
test_payload = {
    "timestamp": datetime.now().isoformat(),
    "hostname": "TEST-MACHINE",
    "sandbox_info": {
        "is_vm": True,
        "vm_indicators": ["VMware", "VirtualBox"],
        "registry_indicators": ["HKLM\\SOFTWARE\\VMware"],
        "disk_indicators": ["VMware Virtual Disk"],
        "cpu_temperature": 45.5,
        "window_count": 12,
        "has_debug_privilege": False
    },
    "system_info": {
        "os": "Windows 10 Pro",
        "architecture": "x64",
        "cpu_count": 4,
        "total_ram_mb": 8192,
        "total_disk_bytes": 500000000000,
        "bios": "American Megatrends Inc.",
        "processes": [
            {"pid": 1234, "name": "chrome.exe", "owner": "user", "path": "C:\\Program Files\\Google\\Chrome\\chrome.exe"},
            {"pid": 5678, "name": "explorer.exe", "owner": "user", "path": "C:\\Windows\\explorer.exe"}
        ],
        "users": ["Administrator", "user"],
        "groups": ["Administrators", "Users"],
        "network_connections": [
            {"protocol": "TCP", "local_addr": "192.168.1.100:443", "remote_addr": "8.8.8.8:443", "state": "ESTABLISHED"}
        ],
        "services": ["wuauserv", "BITS"],
        "environment_variables": {
            "PATH": "C:\\Windows\\System32",
            "TEMP": "C:\\Users\\user\\AppData\\Local\\Temp"
        },
        "pipes": ["\\\\.\\pipe\\chrome"],
        "screenshot_base64": "",
        "mouse_position": {"x": 100, "y": 200},
        "installed_apps": ["Google Chrome", "Microsoft Office"],
        "recent_files": ["C:\\Users\\user\\Documents\\test.txt"],
        "uptime_seconds": 86400
    },
    "hook_info": {
        "hooked_functions": [
            {"module": "ntdll.dll", "function": "NtCreateFile", "is_hooked": True, "first_bytes": "E9 XX XX XX XX"},
            {"module": "kernel32.dll", "function": "CreateFileW", "is_hooked": False, "first_bytes": "48 89 5C 24 08"}
        ],
        "suspicious_dlls": ["sbiedll.dll"]
    },
    "crawler_info": {
        "scanned_paths": ["C:\\Users\\user\\Documents"],
        "found_files": ["password.txt", "credentials.doc"],
        "total_files": 2
    },
    "edr_info": {
        "detected_products": [
            {"name": "Windows Defender", "type": "AV", "detected": True, "method": "process"},
            {"name": "CrowdStrike", "type": "EDR", "detected": True, "method": "driver"}
        ],
        "running_processes": ["MsMpEng.exe", "CSFalconService.exe"],
        "installed_drivers": ["WdFilter.sys", "csagent.sys"]
    }
}

def test_api():
    """Prueba el endpoint de la API"""
    url = "http://192.168.1.143:8080/api/collect"
    
    print("="*60)
    print("  PRUEBA DEL API - Artefacto Visualizer")
    print("="*60)
    print(f"\nEnviando datos de prueba a: {url}")
    print(f"Hostname: {test_payload['hostname']}")
    print(f"Timestamp: {test_payload['timestamp']}")
    
    try:
        response = requests.post(
            url,
            json=test_payload,
            headers={"Content-Type": "application/json"},
            timeout=10
        )
        
        print(f"\n✅ Respuesta recibida:")
        print(f"   Status Code: {response.status_code}")
        print(f"   Response: {response.json()}")
        
        if response.status_code in [200, 201]:
            guid = response.json().get('guid')
            print(f"\n🎉 ¡Éxito! Datos guardados con GUID: {guid}")
            print(f"\n📊 Ver detalles en:")
            print(f"   http://192.168.1.143:8080/execution/{guid}/")
        else:
            print(f"\n⚠️  Respuesta inesperada del servidor")
            
    except requests.exceptions.ConnectionError:
        print("\n❌ Error: No se pudo conectar al servidor")
        print("   Verifica que el servidor esté corriendo en 192.168.1.143:8080")
    except Exception as e:
        print(f"\n❌ Error: {e}")

if __name__ == '__main__':
    test_api()
