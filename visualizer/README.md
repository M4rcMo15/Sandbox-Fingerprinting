# Artefacto Visualizer - Servidor Analizador

Servidor Django que implementa la arquitectura **Agente Recolector + Servidor Analizador**. El servidor recibe datos en bruto del agente y realiza todo el análisis de forma centralizada.

## 🏗️ Arquitectura

```
Agente (Go)              →    Servidor (Django)
├─ Recopila datos en bruto    ├─ VMDetector (análisis de VM)
├─ Inyecta XSS payloads       ├─ EDRDetector (16 productos)
└─ Envía al servidor          ├─ ToolsDetector (25+ herramientas)
                              ├─ GeoLocator (geolocalización)
                              └─ Dashboard web profesional
```

## 🎨 Interfaz

**Nueva interfaz profesional con:**
- Bootstrap 5.3 - Framework CSS moderno
- DataTables - Tablas interactivas con búsqueda, filtrado y exportación
- Chart.js - Gráficos estadísticos
- Diseño limpio y minimalista en modo claro
- Responsive y mobile-friendly

## 📦 Instalación

### 1. Instalar dependencias

```bash
pip install -r requirements.txt
```

### 2. Configurar base de datos

```bash
python manage.py makemigrations
python manage.py migrate
```

### 3. Crear superusuario (opcional)

```bash
python manage.py createsuperuser
```

## 🚀 Uso

### Iniciar el servidor

**Desarrollo:**
```bash
python manage.py runserver 0.0.0.0:8000
```

**Producción:**
```bash
gunicorn visualizer.wsgi:application --bind 0.0.0.0:8000
```

### Acceder a la aplicación

- **Dashboard principal**: http://localhost:8000/
  - Vista de ejecuciones con DataTable interactiva
  - Búsqueda instantánea y filtros
  - Exportación a CSV
  
- **Vista detallada**: http://localhost:8000/execution/{guid}/
  - Tabs organizados: System, Detection, Network, Security, Raw Data
  - DataTables para procesos, conexiones, hooks
  - Summary cards con métricas clave
  
- **Estadísticas**: http://localhost:8000/statistics/
  - KPI cards con métricas principales
  - Gráficos interactivos (Chart.js)
  - Tablas detalladas con accordion
  
- **XSS Audit**: http://localhost:8000/dashboard/
  - Dashboard de auditoría XSS
  - Tracking de payloads
  - Identificación de sandboxes vulnerables
  
- **Admin panel**: http://localhost:8000/admin/
- **API endpoint**: http://localhost:8000/api/collect (POST)

## 🔧 Configurar el agente

Actualiza `artefacto/.env`:

```env
SERVER_URL=http://your-server.com:8000/api/collect
TIMEOUT=120s
```

## ✨ Características

### Análisis Automático
- ✅ **Detección de VM/Sandbox** - Múltiples indicadores (archivos, registro, CPU, disco)
- ✅ **Detección de EDR/AV** - 16 productos principales (Defender, CrowdStrike, SentinelOne, etc.)
- ✅ **Detección de herramientas** - 25+ herramientas de análisis en 5 categorías
- ✅ **Geolocalización** - Por IP pública (país, ciudad, ISP, coordenadas)

### XSS Audit
- ✅ **27 payloads únicos** por ejecución
- ✅ **11 vectores de ataque** (hostname, process, registry, DNS, HTTP, etc.)
- ✅ **Tracking automático** de callbacks
- ✅ **Identificación de sandboxes** vulnerables

### Visualización
- ✅ **Dashboard profesional** con Bootstrap 5
- ✅ **DataTables interactivas** con búsqueda, filtrado y exportación
- ✅ **Tabs organizados** para información detallada
- ✅ **Gráficos estadísticos** con Chart.js
- ✅ **Diseño limpio** y minimalista en modo claro
- ✅ **Responsive** y mobile-friendly

## 📊 Analizadores Implementados

### VMDetector
Analiza múltiples indicadores para determinar si es VM:
- Archivos de VM (VBoxMouse.sys, vmware.sys, etc.)
- Claves de registro de virtualización
- Identificadores de disco (VBOX HARDDISK, etc.)
- Temperatura de CPU (VMs = 0.0)
- Número de ventanas abiertas (VMs < 10)

**Precisión:** 95%+ con múltiples indicadores

### EDRDetector
Detecta 16 productos EDR/AV principales:
- Windows Defender, CrowdStrike Falcon, SentinelOne
- Carbon Black, Cylance, Symantec, McAfee
- Kaspersky, Trend Micro, ESET, Palo Alto
- FireEye, Sophos, Avast, AVG, Bitdefender, Norton

**Métodos:** Procesos + drivers específicos

### ToolsDetector
Identifica herramientas de análisis en 5 categorías:
- **Reversing:** IDA Pro, Ghidra, Binary Ninja, Radare2, Hopper
- **Debugging:** x64dbg, WinDbg, OllyDbg, Immunity Debugger, GDB
- **Monitoring:** Process Monitor, Wireshark, Fiddler, TCPView
- **Virtualization:** VMware, VirtualBox, Hyper-V, QEMU, Parallels
- **Analysis:** Cuckoo, CAPE, Joe Sandbox, Any.Run, Hybrid Analysis

### GeoLocator
Geolocalización automática por IP:
- País, región, ciudad
- Coordenadas (latitud, longitud)
- ISP y organización
- **API:** ip-api.com (gratuita)

## 🔄 Compatibilidad

El servidor es compatible con:
- ✅ **Agentes nuevos** (v2.x) - Envían `raw_data` para análisis
- ✅ **Agentes antiguos** (v1.x) - Envían datos ya procesados

## 📁 Estructura de datos

### Payload del agente (nuevo)

```json
{
  "timestamp": "2024-01-09T10:00:00Z",
  "hostname": "PC-VICTIM",
  "public_ip": "1.2.3.4",
  "raw_data": {
    "vm_files": ["C:\\Windows\\System32\\drivers\\VBoxMouse.sys"],
    "registry_keys": [{"path": "...", "exists": true}],
    "security_processes": ["MsMpEng.exe"],
    "drivers": ["WdFilter.sys", "VBoxMouse.sys"],
    "disk_info": {"identifier": "VBOX HARDDISK"},
    "cpu_info": {"temperature": 0.0},
    "window_count": 8
  },
  "system_info": {
    "os": "Windows 10 Pro",
    "processes": [...],
    "installed_apps": [...]
  },
  "xss_payloads": [
    {"id": "abc123", "type": "img-onerror", "vector": "hostname"}
  ]
}
```

### Respuesta del servidor

```json
{
  "status": "success",
  "execution_id": "12345-abcde-67890",
  "message": "Data processed and analyzed successfully"
}
```

## 🔒 Seguridad

### Autenticación HTTP Basic (recomendado)

```bash
# Nginx
sudo htpasswd -c /etc/nginx/auth/.htpasswd username

# Configurar en nginx.conf
auth_basic "Restricted";
auth_basic_user_file /etc/nginx/auth/.htpasswd;
```

### HTTPS con Let's Encrypt

```bash
sudo certbot --nginx -d your-domain.com
```

### Firewall

```bash
sudo ufw allow from YOUR_IP to any port 8000
sudo ufw enable
```

## 📈 Logs del servidor

El servidor muestra logs detallados:

```
[Analysis] VM Detection: True (5 indicators)
[Analysis] EDR Detection: 1 products found
[Analysis] Tools Detection: 3 tools found
[Analysis] Geolocation: Madrid, Spain
[XSS] Registered 27 payloads
[Server] Analysis completed for execution 12345-abcde
```

## 🐛 Troubleshooting

### Error de conexión
```bash
# Verificar que el servidor esté corriendo
netstat -an | findstr 8000

# Verificar firewall
sudo ufw status
```

### Error de base de datos
```bash
# Recrear migraciones
python manage.py makemigrations
python manage.py migrate
```

### Error de geolocalización
```bash
# Verificar conectividad
curl http://ip-api.com/json/8.8.8.8
```

## 📚 Documentación adicional

- **Arquitectura completa:** `docs/VISUALIZER_ANALYZER_UPDATE.md`
- **Refactorización del agente:** `docs/REFACTORING_AGENT.md`
- **Despliegue en producción:** `visualizer/deploy/`

## 🤝 Contribuciones

Las contribuciones son bienvenidas. Por favor:
1. Fork del repositorio
2. Crear feature branch
3. Commit cambios
4. Push al branch
5. Crear Pull Request

## 📄 Licencia

MIT License - Ver `LICENSE` para más detalles

