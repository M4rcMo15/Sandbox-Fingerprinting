# 🏗️ Arquitectura del Sistema Artefacto + Visualizer

## Flujo de Datos

```
┌─────────────────────────────────────────────────────────────────┐
│                         AGENTE (Go)                             │
│                      artefacto/agent.exe                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   Sandbox    │  │   System     │  │    Hook      │         │
│  │   Detector   │  │   Info       │  │   Detector   │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐                           │
│  │   Crawler    │  │     EDR      │                           │
│  │              │  │   Detector   │                           │
│  └──────────────┘  └──────────────┘                           │
│                                                                 │
│                    ↓ Recopila datos                            │
│                                                                 │
│              ┌──────────────────┐                              │
│              │  Payload (JSON)  │                              │
│              └──────────────────┘                              │
└─────────────────────────────────────────────────────────────────┘
                           │
                           │ HTTP POST
                           │ http://192.168.1.143:8080/api/collect
                           ↓
┌─────────────────────────────────────────────────────────────────┐
│                    VISUALIZER (Django)                          │
│                   192.168.1.143:8080                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────────────────────────────────────────────┐      │
│  │              API Endpoint (/api/collect)             │      │
│  │  - Recibe JSON                                       │      │
│  │  - Genera GUID único                                 │      │
│  │  - Guarda en base de datos                           │      │
│  └──────────────────────────────────────────────────────┘      │
│                           │                                     │
│                           ↓                                     │
│  ┌──────────────────────────────────────────────────────┐      │
│  │              Base de Datos (SQLite3)                 │      │
│  │                                                      │      │
│  │  • AgentExecution (GUID único)                      │      │
│  │  • SandboxInfo                                      │      │
│  │  • SystemInfo                                       │      │
│  │    ├─ ProcessInfo (N)                               │      │
│  │    └─ NetworkConnection (N)                         │      │
│  │  • HookInfo                                         │      │
│  │    └─ HookedFunction (N)                            │      │
│  │  • CrawlerInfo                                      │      │
│  │  • EDRInfo                                          │      │
│  │    └─ EDRProduct (N)                                │      │
│  └──────────────────────────────────────────────────────┘      │
│                           │                                     │
│                           ↓                                     │
│  ┌──────────────────────────────────────────────────────┐      │
│  │           Interfaz Web (Templates HTML)              │      │
│  │                                                      │      │
│  │  / ────────────► Lista de ejecuciones               │      │
│  │  /execution/{GUID}/ ──► Detalle completo            │      │
│  │  /admin/ ──────────► Panel de administración        │      │
│  └──────────────────────────────────────────────────────┘      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                           │
                           │ HTTP
                           ↓
                    ┌──────────────┐
                    │   USUARIO    │
                    │   Navegador  │
                    └──────────────┘
```

## Componentes del Sistema

### 1. Agente (Go)

**Ubicación:** `artefacto/`

**Función:** Recopilar información del sistema y enviarla al servidor

**Módulos:**
- `collectors/sandbox.go` - Detecta virtualización y sandboxes
- `collectors/sysinfo.go` - Información del sistema
- `collectors/hooks.go` - Detecta hooks en funciones
- `collectors/crawler.go` - Busca archivos específicos
- `collectors/edr.go` - Detecta EDR/AV
- `exfil/sender.go` - Envía datos al servidor

**Configuración:** `.env`
```env
SERVER_URL=http://192.168.1.143:8080/api/collect
TIMEOUT=30s
```

### 2. Visualizer (Django)

**Ubicación:** `visualizer/`

**Función:** Recibir, almacenar y visualizar los datos del agente

**Estructura:**
```
visualizer/
├── collector/              # App principal
│   ├── models.py          # 10 modelos de BD
│   ├── views.py           # 3 vistas + API endpoint
│   ├── urls.py            # Rutas
│   ├── admin.py           # Configuración admin
│   ├── middleware.py      # CSRF exempt para API
│   └── templates/         # Plantillas HTML
│       └── collector/
│           ├── base.html      # Template base
│           ├── index.html     # Lista de ejecuciones
│           └── detail.html    # Detalle completo
│
├── visualizer/            # Configuración Django
│   ├── settings.py        # Configuración principal
│   ├── urls.py            # URLs raíz
│   └── wsgi.py            # WSGI
│
├── manage.py              # CLI Django
├── requirements.txt       # Dependencias
├── start_server.bat       # Inicio Windows
├── start_server.sh        # Inicio Linux
└── test_api.py           # Script de prueba
```

## Modelos de Base de Datos

### Relaciones

```
AgentExecution (1) ──┬── (1) SandboxInfo
                     │
                     ├── (1) SystemInfo ──┬── (N) ProcessInfo
                     │                    └── (N) NetworkConnection
                     │
                     ├── (1) HookInfo ──── (N) HookedFunction
                     │
                     ├── (1) CrawlerInfo
                     │
                     └── (1) EDRInfo ────── (N) EDRProduct
```

### Campos Principales

**AgentExecution**
- `guid` (UUID, PK) - Identificador único de ejecución
- `timestamp` - Momento de ejecución
- `hostname` - Nombre del equipo
- `received_at` - Momento de recepción

**SandboxInfo**
- `is_vm` - ¿Es máquina virtual?
- `vm_indicators` - Lista de indicadores
- `cpu_temperature` - Temperatura CPU
- `window_count` - Ventanas abiertas

**SystemInfo**
- `os`, `architecture`, `cpu_count`, `total_ram_mb`
- `users`, `groups`, `services`
- `environment_variables` (JSON)
- `screenshot_base64` - Screenshot en base64

**ProcessInfo**
- `pid`, `name`, `owner`, `path`

**NetworkConnection**
- `protocol`, `local_addr`, `remote_addr`, `state`

**HookedFunction**
- `module`, `function`, `is_hooked`, `first_bytes`

**EDRProduct**
- `name`, `type`, `detected`, `method`

## Flujo de Ejecución

### 1. Ejecución del Agente

```
1. Agente inicia
2. Lee configuración (.env)
3. Ejecuta colectores en paralelo (goroutines)
4. Construye payload JSON
5. Envía POST a /api/collect
6. Recibe confirmación con GUID
```

### 2. Recepción en Visualizer

```
1. Django recibe POST en /api/collect
2. Middleware desactiva CSRF para /api/*
3. View parsea JSON
4. Genera GUID único (UUID4)
5. Crea AgentExecution
6. Crea registros relacionados (SandboxInfo, SystemInfo, etc.)
7. Guarda en SQLite3
8. Retorna JSON con GUID
```

### 3. Visualización

```
1. Usuario accede a http://192.168.1.143:8080/
2. Vista index lista todas las ejecuciones
3. Usuario hace clic en una ejecución
4. Vista detail carga todos los datos relacionados
5. Template renderiza con desplegables
6. Usuario navega por la información
```

## Características de Seguridad

### Agente
- ✅ Timeout configurable
- ✅ Manejo de errores
- ✅ User-Agent personalizado
- ⚠️ Sin cifrado (agregar HTTPS en producción)

### Visualizer
- ✅ CSRF protection (excepto API)
- ✅ SQL injection protection (ORM Django)
- ✅ XSS protection (auto-escape templates)
- ⚠️ DEBUG=True (cambiar en producción)
- ⚠️ SECRET_KEY por defecto (cambiar en producción)

## Escalabilidad

### Actual (Desarrollo)
- SQLite3 local
- Servidor single-threaded
- Sin autenticación
- Sin rate limiting

### Producción (Recomendado)
- PostgreSQL/MySQL
- Gunicorn + Nginx
- Autenticación JWT
- Rate limiting
- HTTPS
- Logging centralizado
- Backup automático

## Extensiones Futuras

### Agente
- [ ] Cifrado de payload
- [ ] Compresión de datos
- [ ] Modo stealth
- [ ] Persistencia local si falla envío
- [ ] Múltiples servidores C2

### Visualizer
- [ ] Dashboard con estadísticas
- [ ] Búsqueda y filtros avanzados
- [ ] Exportación a PDF/CSV
- [ ] Comparación entre ejecuciones
- [ ] Alertas automáticas
- [ ] API REST completa
- [ ] WebSockets para updates en tiempo real
- [ ] Autenticación multi-usuario
- [ ] Roles y permisos

## Performance

### Agente
- Ejecución paralela de colectores
- Timeout por colector
- Memoria: ~50MB
- Tiempo ejecución: ~5-10 segundos

### Visualizer
- Queries optimizadas con select_related/prefetch_related
- Índices en campos clave (GUID, timestamp)
- Paginación para listas grandes
- Lazy loading de screenshots

## Monitoreo

### Logs del Agente
```bash
# Salida estándar
[*] Iniciando agente...
[+] Ejecutando CheckSandbox...
[✓] CheckSandbox completado
...
[✓] Datos enviados correctamente
```

### Logs del Visualizer
```bash
# Django development server
[28/Nov/2025 10:30:45] "POST /api/collect HTTP/1.1" 201 123
[28/Nov/2025 10:31:12] "GET / HTTP/1.1" 200 4567
[28/Nov/2025 10:31:15] "GET /execution/550e8400-e29b-41d4-a716-446655440000/ HTTP/1.1" 200 12345
```

## Testing

### Test del Agente
```bash
cd artefacto
go test ./...
```

### Test del Visualizer
```bash
cd visualizer
python manage.py test
python test_api.py  # Test manual del endpoint
```

## Deployment

### Desarrollo (Actual)
```bash
# Terminal 1: Visualizer
cd visualizer
python manage.py runserver 192.168.1.143:8080

# Terminal 2: Agente
cd artefacto
./agent.exe
```

### Producción
```bash
# Visualizer con Gunicorn
gunicorn visualizer.wsgi:application --bind 192.168.1.143:8080 --workers 4

# Nginx como reverse proxy
# HTTPS con Let's Encrypt
# Systemd service para auto-start
```
