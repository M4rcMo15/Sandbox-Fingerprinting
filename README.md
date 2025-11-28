# 🔍 Artefacto - Sistema de Fingerprinting y Visualización

Sistema completo de recopilación y visualización de información de sistemas para operaciones de Red Team.

## 📁 Estructura del Proyecto

```
.
├── artefacto/              # Agente de recopilación (Go)
│   ├── collectors/         # Módulos de detección
│   ├── config/            # Configuración
│   ├── exfil/             # Exfiltración de datos
│   ├── models/            # Estructuras de datos
│   ├── utils/             # Utilidades
│   ├── .env               # Configuración del agente
│   └── agent.exe          # Ejecutable compilado
│
├── visualizer/            # Servidor web de visualización (Django)
│   ├── collector/         # App Django principal
│   ├── visualizer/        # Configuración Django
│   ├── manage.py          # CLI de Django
│   ├── requirements.txt   # Dependencias Python
│   └── *.bat / *.sh       # Scripts de inicio
│
├── SETUP_VISUALIZER.md    # Guía de instalación completa
└── README.md              # Este archivo
```

## 🚀 Inicio Rápido

### 1. Configurar el Visualizer

```bash
cd visualizer
pip install -r requirements.txt
python manage.py makemigrations
python manage.py migrate
```

### 2. Iniciar el Servidor

**Windows:**
```cmd
cd visualizer
start_server.bat
```

**Linux/Mac:**
```bash
cd visualizer
chmod +x start_server.sh
./start_server.sh
```

El servidor estará disponible en: **http://192.168.1.143:8080/**

### 3. Ejecutar el Agente

```bash
cd artefacto
./agent.exe
```

## 📊 Características

### Agente (Go)
- ✅ Detección de máquinas virtuales y sandboxes
- ✅ Recopilación de información del sistema
- ✅ Detección de hooks en funciones
- ✅ Crawler de archivos
- ✅ Detección de EDR/AV
- ✅ Captura de screenshots
- ✅ Ejecución paralela de colectores
- ✅ Envío automático al servidor

### Visualizer (Django)
- ✅ Recepción automática de datos
- ✅ GUID único por cada ejecución
- ✅ Interfaz web moderna (tema oscuro)
- ✅ Desplegables para organizar información
- ✅ Tablas para procesos, conexiones, hooks
- ✅ Visualización de screenshots
- ✅ Panel de administración
- ✅ Base de datos SQLite3
- ✅ API REST para recepción de datos

## 🎯 Información Recopilada

### 🔍 Detección de Sandbox
- Indicadores de VM (VMware, VirtualBox, Hyper-V, etc.)
- Indicadores de registro
- Indicadores de disco
- Temperatura de CPU
- Conteo de ventanas
- Privilegios de debug

### 💻 Información del Sistema
- Sistema operativo y arquitectura
- CPU, RAM, disco
- Procesos en ejecución (PID, nombre, owner, path)
- Conexiones de red (protocolo, direcciones, estado)
- Usuarios y grupos
- Servicios
- Variables de entorno
- Named pipes
- Aplicaciones instaladas
- Archivos recientes
- Posición del mouse
- Screenshot

### 🪝 Detección de Hooks
- Funciones analizadas (módulo, función, estado)
- Primeros bytes de cada función
- DLLs sospechosas cargadas

### 📁 Crawler de Archivos
- Rutas escaneadas
- Archivos encontrados
- Patrones buscados

### 🛡️ Detección de EDR/AV
- Productos detectados (nombre, tipo, método)
- Procesos de seguridad en ejecución
- Drivers de seguridad instalados

## 📖 Documentación

### Documentación Principal
- **[SETUP_VISUALIZER.md](SETUP_VISUALIZER.md)** - Guía completa de instalación y configuración
- **[visualizer/INSTRUCCIONES.md](visualizer/INSTRUCCIONES.md)** - Instrucciones detalladas en español
- **[visualizer/ARQUITECTURA.md](visualizer/ARQUITECTURA.md)** - Arquitectura del sistema
- **[visualizer/TROUBLESHOOTING.md](visualizer/TROUBLESHOOTING.md)** - Solución de problemas

### Documentación del Agente
- **[artefacto/README.md](artefacto/README.md)** - Documentación del agente
- **[artefacto/TECHNICAL.md](artefacto/TECHNICAL.md)** - Detalles técnicos
- **[artefacto/USAGE.md](artefacto/USAGE.md)** - Guía de uso

## 🌐 URLs del Visualizer

Una vez iniciado el servidor:

- **Página principal**: http://192.168.1.143:8080/
  - Lista todas las ejecuciones del agente
  
- **Detalle de ejecución**: http://192.168.1.143:8080/execution/{GUID}/
  - Muestra toda la información de una ejecución
  
- **API Endpoint**: http://192.168.1.143:8080/api/collect
  - Recibe datos del agente (POST)
  
- **Admin Panel**: http://192.168.1.143:8080/admin/
  - Panel de administración (requiere superusuario)

## 🔧 Configuración

### Agente (artefacto/.env)
```env
SERVER_URL=http://192.168.1.143:8080/api/collect
DEBUG=0
TIMEOUT=30s
```

### Visualizer (visualizer/visualizer/settings.py)
```python
ALLOWED_HOSTS = ['192.168.1.143', 'localhost', '127.0.0.1']
DEBUG = True  # Cambiar a False en producción
```

## 🧪 Probar el Sistema

### Opción 1: Ejecutar el agente real
```bash
cd artefacto
./agent.exe
```

### Opción 2: Usar el script de prueba
```bash
cd visualizer
python test_api.py
```

## 🗄️ Base de Datos

El visualizer usa SQLite3 con los siguientes modelos:

- `AgentExecution` - Ejecución principal (GUID único)
- `SandboxInfo` - Detección de sandbox
- `SystemInfo` - Información del sistema
- `ProcessInfo` - Procesos individuales
- `NetworkConnection` - Conexiones de red
- `HookInfo` - Información de hooks
- `HookedFunction` - Funciones hooked
- `CrawlerInfo` - Resultados del crawler
- `EDRInfo` - Información de EDR
- `EDRProduct` - Productos EDR/AV detectados

## 🔐 Seguridad

⚠️ **Este sistema está configurado para desarrollo local**

Para uso en producción:
1. Cambiar `SECRET_KEY` en Django
2. Establecer `DEBUG = False`
3. Configurar HTTPS
4. Agregar autenticación
5. Configurar firewall
6. Usar PostgreSQL/MySQL en lugar de SQLite
7. Implementar rate limiting
8. Cifrar comunicaciones

## 📦 Dependencias

### Agente (Go)
Ver `artefacto/go.mod`

### Visualizer (Python)
```
Django>=4.2,<5.0
djangorestframework>=3.14.0
```

## 🛠️ Desarrollo

### Compilar el agente
```bash
cd artefacto
go build -o agent.exe main.go
```

### Ejecutar tests
```bash
# Agente
cd artefacto
go test ./...

# Visualizer
cd visualizer
python manage.py test
```

## 📝 Notas

- Cada ejecución del agente genera un **GUID único** automáticamente
- Los datos se almacenan permanentemente en la base de datos
- El servidor puede recibir múltiples ejecuciones simultáneas
- La interfaz web no requiere JavaScript
- Compatible con Windows, Linux y macOS

## 🆘 Soporte

Si encuentras problemas:

1. Consulta [TROUBLESHOOTING.md](visualizer/TROUBLESHOOTING.md)
2. Verifica los logs del servidor Django
3. Verifica los logs del agente
4. Prueba con `test_api.py`
5. Verifica la configuración de red

## 📄 Licencia

Este proyecto es para fines educativos y de investigación en seguridad.

## ⚠️ Disclaimer

Este software está diseñado para operaciones legítimas de Red Team y pruebas de seguridad. El uso indebido de esta herramienta puede ser ilegal. Úsala solo en sistemas donde tengas autorización explícita.

---

**Desarrollado para operaciones de Red Team y análisis de seguridad**
