# 🔧 Troubleshooting - Artefacto Visualizer

## Problemas Comunes y Soluciones

### 1. Error: "No module named 'django'"

**Problema:** Django no está instalado

**Solución:**
```bash
pip install Django
# O instalar todas las dependencias
pip install -r requirements.txt
```

---

### 2. Error: "No such table: collector_agentexecution"

**Problema:** Las migraciones no se han aplicado

**Solución:**
```bash
python manage.py makemigrations
python manage.py migrate
```

---

### 3. Error: "That port is already in use"

**Problema:** El puerto 8080 ya está en uso

**Solución 1:** Cambiar el puerto
```bash
python manage.py runserver 192.168.1.143:8081
```

**Solución 2:** Encontrar y matar el proceso
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux
lsof -i :8080
kill -9 <PID>
```

---

### 4. Error: "Connection refused" desde el agente

**Problema:** El servidor no está corriendo o no es accesible

**Diagnóstico:**
```bash
# Verificar que el servidor esté corriendo
curl http://192.168.1.143:8080/

# Verificar conectividad
ping 192.168.1.143
```

**Soluciones:**

1. **Verificar que el servidor esté corriendo:**
```bash
cd visualizer
python manage.py runserver 192.168.1.143:8080
```

2. **Verificar firewall (Windows):**
```cmd
netsh advfirewall firewall add rule name="Django8080" dir=in action=allow protocol=TCP localport=8080
```

3. **Verificar firewall (Linux):**
```bash
sudo ufw allow 8080
```

4. **Verificar la IP:**
```bash
# Windows
ipconfig

# Linux
ip addr show
```

---

### 5. Error: "CSRF verification failed"

**Problema:** El middleware CSRF no está configurado correctamente

**Solución:** Verificar que el middleware esté en `settings.py`:
```python
MIDDLEWARE = [
    ...
    'collector.middleware.DisableCSRFForAPIMiddleware',
    'django.middleware.csrf.CsrfViewMiddleware',
    ...
]
```

---

### 6. Error: "Invalid JSON" en el endpoint

**Problema:** El payload enviado no es JSON válido

**Diagnóstico:**
```bash
# Probar el endpoint manualmente
curl -X POST http://192.168.1.143:8080/api/collect \
  -H "Content-Type: application/json" \
  -d '{"timestamp":"2024-01-01T12:00:00Z","hostname":"test"}'
```

**Solución:** Verificar que el agente esté enviando JSON correcto

---

### 7. No se muestran los datos en la interfaz

**Problema:** Los datos se guardan pero no se muestran

**Diagnóstico:**
```bash
# Verificar en la base de datos
python manage.py shell
>>> from collector.models import AgentExecution
>>> AgentExecution.objects.all()
>>> AgentExecution.objects.count()
```

**Soluciones:**

1. **Verificar que hay datos:**
```python
python manage.py shell
>>> from collector.models import AgentExecution
>>> for exe in AgentExecution.objects.all():
...     print(exe.guid, exe.hostname)
```

2. **Limpiar caché del navegador:**
- Ctrl + F5 (Windows)
- Cmd + Shift + R (Mac)

3. **Verificar templates:**
```bash
# Asegurarse de que existan
ls visualizer/collector/templates/collector/
```

---

### 8. Error: "OperationalError: database is locked"

**Problema:** SQLite está bloqueado por otro proceso

**Solución:**
```bash
# Cerrar todos los procesos que usen la BD
# Reiniciar el servidor
python manage.py runserver 192.168.1.143:8080
```

---

### 9. Screenshots no se muestran

**Problema:** El campo screenshot_base64 está vacío o mal formateado

**Diagnóstico:**
```python
python manage.py shell
>>> from collector.models import SystemInfo
>>> si = SystemInfo.objects.first()
>>> len(si.screenshot_base64)  # Debe ser > 0
>>> si.screenshot_base64[:50]  # Ver primeros caracteres
```

**Solución:** Verificar que el agente esté capturando screenshots correctamente

---

### 10. Error: "DisallowedHost at /"

**Problema:** La IP no está en ALLOWED_HOSTS

**Solución:** Editar `visualizer/visualizer/settings.py`:
```python
ALLOWED_HOSTS = ['192.168.1.143', 'localhost', '127.0.0.1', 'TU_IP_AQUI']
```

---

### 11. El agente envía datos pero retorna error 400

**Problema:** El formato del JSON no coincide con lo esperado

**Diagnóstico:** Ver logs del servidor Django

**Solución:** Verificar que el payload del agente coincida con el formato esperado:
```json
{
  "timestamp": "ISO8601 format",
  "hostname": "string",
  "sandbox_info": { ... },
  "system_info": { ... },
  "hook_info": { ... },
  "crawler_info": { ... },
  "edr_info": { ... }
}
```

---

### 12. Error: "ImportError: No module named 'rest_framework'"

**Problema:** Django REST Framework no está instalado

**Solución:**
```bash
pip install djangorestframework
```

---

### 13. Lentitud al cargar detalles de ejecución

**Problema:** Muchos datos relacionados (procesos, conexiones, etc.)

**Solución:** Ya está optimizado con `select_related` y `prefetch_related` en las vistas

Si persiste:
```python
# En views.py, agregar paginación
from django.core.paginator import Paginator
```

---

### 14. Error: "Permission denied" al ejecutar start_server.sh

**Problema:** El script no tiene permisos de ejecución

**Solución:**
```bash
chmod +x start_server.sh
./start_server.sh
```

---

### 15. No se puede acceder desde otra máquina

**Problema:** El servidor solo escucha en localhost

**Solución:** Asegurarse de usar la IP correcta:
```bash
# NO usar
python manage.py runserver

# SÍ usar
python manage.py runserver 192.168.1.143:8080
# O para todas las interfaces
python manage.py runserver 0.0.0.0:8080
```

---

## Comandos Útiles de Diagnóstico

### Verificar estado del servidor
```bash
# Ver procesos Python
ps aux | grep python

# Ver puertos en uso
netstat -tulpn | grep 8080
```

### Verificar base de datos
```bash
python manage.py dbshell
sqlite> .tables
sqlite> SELECT COUNT(*) FROM collector_agentexecution;
sqlite> .quit
```

### Ver logs detallados
```bash
# Iniciar con verbosidad
python manage.py runserver 192.168.1.143:8080 --verbosity 3
```

### Limpiar base de datos
```bash
# CUIDADO: Esto borra todos los datos
python manage.py flush
```

### Recrear base de datos desde cero
```bash
# CUIDADO: Esto borra todos los datos
rm db.sqlite3
python manage.py makemigrations
python manage.py migrate
```

### Verificar configuración
```bash
python manage.py check
python manage.py check --deploy  # Para producción
```

---

## Logs y Debugging

### Activar modo debug en el agente
Editar `artefacto/.env`:
```env
DEBUG=1
```

### Ver queries SQL en Django
Editar `visualizer/visualizer/settings.py`:
```python
LOGGING = {
    'version': 1,
    'handlers': {
        'console': {
            'class': 'logging.StreamHandler',
        },
    },
    'loggers': {
        'django.db.backends': {
            'handlers': ['console'],
            'level': 'DEBUG',
        },
    },
}
```

### Usar Django shell para debugging
```bash
python manage.py shell
>>> from collector.models import *
>>> from collector.views import *
>>> # Probar queries, funciones, etc.
```

---

## Contacto y Soporte

Si ninguna de estas soluciones funciona:

1. Verificar los logs del servidor Django
2. Verificar los logs del agente
3. Probar con el script `test_api.py`
4. Verificar la configuración de red
5. Revisar la documentación de Django: https://docs.djangoproject.com/

---

## Checklist de Verificación

Antes de reportar un problema, verificar:

- [ ] Python 3.8+ instalado
- [ ] Django instalado (`pip list | grep Django`)
- [ ] Migraciones aplicadas (`python manage.py showmigrations`)
- [ ] Servidor corriendo en la IP correcta
- [ ] Firewall permite conexiones en puerto 8080
- [ ] Agente configurado con URL correcta
- [ ] Conectividad de red (ping, curl)
- [ ] No hay errores en los logs del servidor
- [ ] Navegador actualizado y sin caché
