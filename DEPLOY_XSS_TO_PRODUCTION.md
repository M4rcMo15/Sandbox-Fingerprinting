# 🚀 Despliegue del Módulo XSS Audit a Producción

## Resumen de Cambios

Se ha implementado un módulo completo de auditoría XSS que permite detectar vulnerabilidades en sandboxes de análisis de malware.

### Archivos Nuevos

**Agente (Go):**
```
artefacto/
├── xss/
│   ├── payloads.go          # Definición de 10 payloads XSS variados
│   └── injector.go          # Lógica de inyección en diferentes vectores
├── test_xss_audit.bat       # Script de prueba local
└── .env.example             # Actualizado con XSS_AUDIT y CALLBACK_SERVER
```

**Visualizer (Django):**
```
visualizer/
├── xss_audit/               # Nueva app Django
│   ├── __init__.py
│   ├── models.py            # XSSPayload, XSSHit, SandboxVulnerability
│   ├── views.py             # Dashboard y endpoint de callback
│   ├── urls.py              # Rutas
│   ├── admin.py             # Admin
│   └── templates/
│       └── xss_audit/
│           └── dashboard.html  # Dashboard visual
└── deploy/
    └── deploy_xss_audit.sh  # Script de despliegue automático
```

**Documentación:**
```
├── XSS_AUDIT_README.md           # README del módulo
├── XSS_AUDIT_GUIDE.md            # Guía completa de uso
└── DEPLOY_XSS_TO_PRODUCTION.md   # Este archivo
```

### Archivos Modificados

**Agente:**
- `artefacto/models/payload.go` - Añadido `XSSPayloads []XSSPayloadMetadata`
- `artefacto/config/config.go` - Añadido `XSSAudit` y `CallbackServer`
- `artefacto/main.go` - Integración del modo XSS Audit

**Visualizer:**
- `visualizer/visualizer/settings.py` - Añadido `'xss_audit'` a `INSTALLED_APPS`
- `visualizer/visualizer/urls.py` - Añadido `include('xss_audit.urls')`
- `visualizer/collector/views.py` - Añadido guardado de payloads XSS
- `visualizer/collector/templates/collector/base.html` - Añadido enlace al dashboard XSS

**Documentación:**
- `README.md` - Añadida sección del módulo XSS Audit

---

## 📋 Checklist de Despliegue

### Paso 1: Backup

```bash
# Conectar al servidor
ssh ubuntu@54.37.226.179

# Crear backup de la base de datos
cd /opt/artefacto-visualizer
sudo cp db.sqlite3 ../db_backup_$(date +%Y%m%d_%H%M%S).sqlite3

# Crear backup del código
cd /opt
sudo tar -czf artefacto-visualizer-backup-$(date +%Y%m%d).tar.gz artefacto-visualizer/
```

### Paso 2: Subir Código Nuevo

**Opción A: Git (Recomendado)**
```bash
cd /opt/artefacto-visualizer
sudo git pull origin main
```

**Opción B: SCP Manual**
```bash
# Desde tu máquina local
scp -r visualizer/xss_audit ubuntu@54.37.226.179:/opt/artefacto-visualizer/
scp visualizer/visualizer/settings.py ubuntu@54.37.226.179:/opt/artefacto-visualizer/visualizer/
scp visualizer/visualizer/urls.py ubuntu@54.37.226.179:/opt/artefacto-visualizer/visualizer/
scp visualizer/collector/views.py ubuntu@54.37.226.179:/opt/artefacto-visualizer/collector/
scp visualizer/collector/templates/collector/base.html ubuntu@54.37.226.179:/opt/artefacto-visualizer/collector/templates/collector/
scp visualizer/deploy/deploy_xss_audit.sh ubuntu@54.37.226.179:/opt/artefacto-visualizer/deploy/
```

### Paso 3: Aplicar Migraciones

```bash
# Activar entorno virtual
source /opt/venv/bin/activate

# Ir al directorio del proyecto
cd /opt/artefacto-visualizer

# Crear migraciones
python manage.py makemigrations xss_audit

# Aplicar migraciones
python manage.py migrate

# Recolectar archivos estáticos
python manage.py collectstatic --noinput
```

### Paso 4: Verificar Configuración

```bash
# Verificar que xss_audit está en INSTALLED_APPS
grep -A 10 "INSTALLED_APPS" visualizer/settings.py | grep xss_audit

# Verificar URLs
grep "xss_audit" visualizer/urls.py
```

### Paso 5: Reiniciar Servicios

```bash
# Reiniciar Gunicorn
sudo systemctl restart gunicorn

# Verificar estado
sudo systemctl status gunicorn

# Reiniciar Nginx
sudo systemctl restart nginx

# Verificar estado
sudo systemctl status nginx
```

### Paso 6: Verificar Funcionamiento

```bash
# Probar endpoint de callback
curl -I "http://localhost/xss-callback?id=test123&v=test"
# Debe retornar HTTP 200

# Probar dashboard
curl -I "http://localhost/xss-audit/dashboard/"
# Debe retornar HTTP 200

# Ver logs
sudo journalctl -u gunicorn -n 50 | grep -i xss
```

### Paso 7: Probar desde Navegador

1. Abrir: `http://54.37.226.179/xss-audit/dashboard/`
2. Verificar que carga correctamente
3. Verificar que el enlace "🎯 XSS Audit" aparece en el menú

---

## 🔧 Script de Despliegue Automático

Para facilitar el despliegue, usa el script automático:

```bash
# Dar permisos de ejecución
chmod +x /opt/artefacto-visualizer/deploy/deploy_xss_audit.sh

# Ejecutar
cd /opt/artefacto-visualizer
./deploy/deploy_xss_audit.sh
```

El script hace:
1. ✅ Backup de la base de datos
2. ✅ Activar entorno virtual
3. ✅ Aplicar migraciones
4. ✅ Recolectar estáticos
5. ✅ Verificar configuración
6. ✅ Reiniciar servicios
7. ✅ Verificar endpoints
8. ✅ Mostrar resumen

---

## 🧪 Pruebas Post-Despliegue

### 1. Probar Callback Endpoint

```bash
# Desde el servidor
curl -v "http://localhost/xss-callback?id=test123&v=hostname"

# Desde fuera
curl -v "http://54.37.226.179/xss-callback?id=test123&v=hostname"
```

Debe retornar:
- HTTP 200
- Content-Type: image/gif
- Imagen GIF 1x1 transparente

### 2. Probar Dashboard

Abrir en navegador: `http://54.37.226.179/xss-audit/dashboard/`

Debe mostrar:
- Estadísticas (0 payloads inicialmente)
- Secciones vacías con mensajes "No hay datos aún"
- Diseño correcto con tema oscuro

### 3. Probar Flujo Completo

```bash
# En tu máquina local
cd artefacto

# Editar .env
nano .env
# Cambiar: XSS_AUDIT=true

# Compilar
go build -o conhost_test.exe

# Ejecutar
./conhost_test.exe
```

Verificar:
1. Salida muestra `[🎯] MODO XSS AUDIT ACTIVADO`
2. Muestra payloads inyectados
3. Dashboard muestra los payloads en estado "injected"

---

## 📊 Monitoreo

### Logs en Tiempo Real

```bash
# Logs de Gunicorn (filtrar XSS)
sudo journalctl -u gunicorn -f | grep -i xss

# Logs de Nginx (callbacks)
sudo tail -f /var/log/nginx/artefacto-visualizer-access.log | grep xss-callback

# Logs de errores
sudo tail -f /var/log/nginx/artefacto-visualizer-error.log
```

### Verificar Base de Datos

```bash
cd /opt/artefacto-visualizer
source /opt/venv/bin/activate
python manage.py shell
```

```python
from xss_audit.models import XSSPayload, XSSHit, SandboxVulnerability

# Ver payloads
print(f"Total payloads: {XSSPayload.objects.count()}")

# Ver hits
print(f"Total hits: {XSSHit.objects.count()}")

# Ver sandboxes
for s in SandboxVulnerability.objects.all():
    print(f"{s.sandbox_name}: {s.hit_count} hits")
```

---

## 🐛 Troubleshooting

### Error: "No module named 'xss_audit'"

```bash
# Verificar que está en INSTALLED_APPS
grep xss_audit visualizer/settings.py

# Si no está, añadirlo
nano visualizer/settings.py
# Añadir 'xss_audit' a INSTALLED_APPS

# Reiniciar
sudo systemctl restart gunicorn
```

### Error: "Table doesn't exist"

```bash
# Aplicar migraciones
python manage.py migrate xss_audit

# Verificar
python manage.py showmigrations xss_audit
```

### Error 404 en /xss-callback

```bash
# Verificar URLs
python manage.py show_urls | grep xss

# Si no aparece, verificar que está incluido
grep "xss_audit.urls" visualizer/urls.py

# Reiniciar
sudo systemctl restart gunicorn
```

### Callback no se registra

```bash
# Ver logs en tiempo real
sudo journalctl -u gunicorn -f

# Probar manualmente
curl -v "http://localhost/xss-callback?id=test&v=test"

# Verificar en base de datos
python manage.py shell
>>> from xss_audit.models import XSSPayload
>>> XSSPayload.objects.filter(payload_id='test')
```

---

## 🔄 Rollback (Si algo sale mal)

### Restaurar Base de Datos

```bash
# Listar backups
ls -lh /opt/db_backup_*.sqlite3

# Restaurar
cd /opt/artefacto-visualizer
sudo cp /opt/db_backup_YYYYMMDD_HHMMSS.sqlite3 db.sqlite3

# Reiniciar
sudo systemctl restart gunicorn
```

### Restaurar Código

```bash
# Listar backups
ls -lh /opt/artefacto-visualizer-backup-*.tar.gz

# Restaurar
cd /opt
sudo rm -rf artefacto-visualizer
sudo tar -xzf artefacto-visualizer-backup-YYYYMMDD.tar.gz

# Reiniciar
sudo systemctl restart gunicorn
sudo systemctl restart nginx
```

---

## ✅ Checklist Final

- [ ] Backup de base de datos creado
- [ ] Backup de código creado
- [ ] Código nuevo subido al servidor
- [ ] Migraciones aplicadas correctamente
- [ ] Archivos estáticos recolectados
- [ ] Gunicorn reiniciado sin errores
- [ ] Nginx reiniciado sin errores
- [ ] Endpoint /xss-callback responde HTTP 200
- [ ] Dashboard /xss-audit/dashboard/ carga correctamente
- [ ] Enlace "🎯 XSS Audit" aparece en menú
- [ ] Prueba completa con agente funciona
- [ ] Logs no muestran errores

---

## 📞 Soporte

Si encuentras problemas:

1. Revisar logs: `sudo journalctl -u gunicorn -n 100`
2. Verificar configuración: `python manage.py check`
3. Probar endpoints manualmente con curl
4. Revisar este documento de troubleshooting

---

## 🎉 ¡Listo!

El módulo XSS Audit está desplegado y funcionando. Ahora puedes:

1. Compilar agentes con `XSS_AUDIT=true`
2. Subirlos a sandboxes
3. Monitorear el dashboard para ver qué sandboxes son vulnerables
4. Contribuir a la seguridad del ecosistema con disclosure responsable

**Dashboard:** http://54.37.226.179/xss-audit/dashboard/

**Documentación:** [XSS_AUDIT_GUIDE.md](XSS_AUDIT_GUIDE.md)
