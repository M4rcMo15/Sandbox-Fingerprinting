# 🎯 Setup Completo - Artefacto + Visualizer

## Estructura del Proyecto

```
.
├── artefacto/              # Agente de fingerprinting (Go)
│   ├── collectors/         # Módulos de recolección
│   ├── config/            # Configuración
│   ├── exfil/             # Envío de datos
│   ├── .env               # ✅ Actualizado con nueva URL
│   └── agent.exe          # Ejecutable compilado
│
└── visualizer/            # 🆕 Servidor web Django
    ├── collector/         # App principal
    │   ├── models.py      # Modelos de BD
    │   ├── views.py       # Vistas y API
    │   ├── templates/     # Plantillas HTML
    │   └── middleware.py  # Middleware CSRF
    ├── visualizer/        # Configuración Django
    ├── manage.py          # CLI de Django
    ├── requirements.txt   # Dependencias Python
    ├── start_server.bat   # Script Windows
    ├── start_server.sh    # Script Linux
    └── test_api.py        # Script de prueba
```

## 🚀 Inicio Rápido

### 1. Instalar y Configurar Visualizer

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

El agente enviará automáticamente los datos al visualizer.

## 📊 Funcionalidades del Visualizer

### Vista Principal (/)
- Lista de todas las ejecuciones del agente
- Cada ejecución tiene un **GUID único**
- Badges visuales: VM detectada, EDR/AV encontrados
- Ordenadas por fecha de recepción

### Vista de Detalle (/execution/{GUID}/)
Muestra toda la información recopilada:

#### 🔍 Detección de Sandbox
- Indicadores de VM (VMware, VirtualBox, etc.)
- Indicadores de registro
- Indicadores de disco
- Temperatura de CPU
- Ventanas abiertas
- Privilegios de debug

#### 💻 Información del Sistema
- Sistema operativo y arquitectura
- CPU, RAM, disco
- **Procesos** (tabla con PID, nombre, owner, path)
- **Conexiones de red** (protocolo, direcciones, estado)
- Usuarios y grupos
- Servicios
- Variables de entorno
- Pipes
- Aplicaciones instaladas
- Archivos recientes
- Screenshot (si disponible)

#### 🪝 Detección de Hooks
- **Funciones analizadas** (tabla con módulo, función, estado, bytes)
- DLLs sospechosas
- Indicador visual de funciones hooked

#### 📁 Crawler de Archivos
- Rutas escaneadas
- Archivos encontrados
- Total de archivos

#### 🛡️ Detección de EDR/AV
- **Productos detectados** (nombre, tipo, método)
- Procesos en ejecución
- Drivers instalados

## 🎨 Características de la Interfaz

- ✅ Tema oscuro moderno (estilo GitHub)
- ✅ Desplegables (`<details>`) para organizar información
- ✅ Tablas para datos estructurados
- ✅ Badges de colores para estados
- ✅ Responsive y fácil de navegar
- ✅ Sin dependencias de JavaScript
- ✅ Visualización de screenshots en base64

## 🔧 Configuración

### Agente (artefacto/.env)
```env
SERVER_URL=http://192.168.1.143:8080/api/collect
DEBUG=0
TIMEOUT=30s
```

### Servidor (visualizer/visualizer/settings.py)
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

## 📡 Endpoints API

### POST /api/collect
Recibe los datos del agente en formato JSON.

**Request:**
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "hostname": "MACHINE-NAME",
  "sandbox_info": { ... },
  "system_info": { ... },
  "hook_info": { ... },
  "crawler_info": { ... },
  "edr_info": { ... }
}
```

**Response:**
```json
{
  "status": "success",
  "guid": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Data received successfully"
}
```

## 🗄️ Base de Datos

SQLite3 local (`db.sqlite3`) con los siguientes modelos:

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

⚠️ **Configuración para desarrollo local**

Para producción:
1. Cambiar `SECRET_KEY` en `settings.py`
2. Establecer `DEBUG = False`
3. Configurar HTTPS
4. Agregar autenticación
5. Configurar firewall
6. Usar base de datos PostgreSQL/MySQL

## 📝 Notas

- Cada ejecución del agente genera un **GUID único** automáticamente
- Los datos se almacenan permanentemente en la base de datos
- El servidor puede recibir múltiples ejecuciones simultáneas
- La interfaz se actualiza automáticamente al recargar
- No se requiere JavaScript para la visualización

## 🆘 Solución de Problemas

### El agente no puede conectar
```bash
# Verificar que el servidor esté corriendo
curl http://192.168.1.143:8080/

# Verificar firewall
netsh advfirewall firewall add rule name="Django" dir=in action=allow protocol=TCP localport=8080
```

### Error de migraciones
```bash
cd visualizer
python manage.py makemigrations collector
python manage.py migrate
```

### Puerto en uso
```bash
# Cambiar el puerto en start_server.bat/sh
python manage.py runserver 192.168.1.143:8081
```

## 📚 Documentación Adicional

- `visualizer/README.md` - Documentación del visualizer
- `visualizer/INSTRUCCIONES.md` - Instrucciones detalladas en español
- `artefacto/README.md` - Documentación del agente
- `artefacto/TECHNICAL.md` - Detalles técnicos del agente
