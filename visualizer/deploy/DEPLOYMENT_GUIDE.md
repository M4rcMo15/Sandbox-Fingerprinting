# 🚀 Guía de Despliegue en Producción

## Despliegue de Artefacto Visualizer en VPS

**Servidor:** 79.137.76.235  
**Dominio:** releases.life  
**Sistema:** Ubuntu Server 20.04/22.04  
**Usuario:** root

---

## 📋 Requisitos Previos

### En tu máquina local:
- Acceso SSH al servidor: `ssh root@79.137.76.235`
- Dominio `releases.life` apuntando a `79.137.76.235` (DNS configurado)

### En el servidor:
- Ubuntu Server 20.04 o 22.04
- Acceso root
- Puerto 80 y 443 abiertos

---

## 🔧 Paso 1: Preparar el Proyecto Localmente

### 1.1 Comprimir el proyecto

En tu máquina local, desde el directorio del proyecto:

```bash
# Ir al directorio visualizer
cd visualizer

# Crear archivo tar.gz excluyendo archivos innecesarios
tar -czf artefacto-visualizer.tar.gz \
  --exclude='*.pyc' \
  --exclude='__pycache__' \
  --exclude='*.sqlite3' \
  --exclude='staticfiles' \
  --exclude='media' \
  --exclude='.git' \
  --exclude='venv' \
  .
```

Esto creará `artefacto-visualizer.tar.gz` con todo el código necesario.

---

## 🌐 Paso 2: Conectar al Servidor y Preparar el Entorno

### 2.1 Conectar por SSH

```bash
ssh root@79.137.76.235
```

### 2.2 Actualizar el sistema

```bash
apt update && apt upgrade -y
```

### 2.3 Instalar dependencias del sistema

```bash
apt install -y \
  python3 \
  python3-pip \
  python3-venv \
  nginx \
  certbot \
  python3-certbot-nginx \
  git \
  curl \
  ufw
```

### 2.4 Configurar firewall

```bash
# Permitir SSH, HTTP y HTTPS
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp

# Activar firewall
ufw --force enable

# Verificar estado
ufw status
```

---

## 📦 Paso 3: Transferir el Proyecto al Servidor

### 3.1 Desde tu máquina local (nueva terminal)

```bash
# Transferir el archivo comprimido
scp artefacto-visualizer.tar.gz root@79.137.76.235:/tmp/
```

### 3.2 En el servidor, descomprimir

```bash
# Crear directorio del proyecto
mkdir -p /opt/artefacto-visualizer

# Descomprimir
cd /opt/artefacto-visualizer
tar -xzf /tmp/artefacto-visualizer.tar.gz

# Limpiar
rm /tmp/artefacto-visualizer.tar.gz

# Verificar que los archivos estén ahí
ls -la
# Deberías ver: manage.py, requirements.txt, collector/, xss_audit/, etc.
```

---

## 🐍 Paso 4: Configurar Python y Dependencias

### 4.1 Crear entorno virtual

```bash
cd /opt/artefacto-visualizer
python3 -m venv venv
```

### 4.2 Activar entorno virtual e instalar dependencias

```bash
source venv/bin/activate
pip install --upgrade pip
pip install -r requirements.txt
```

Esto instalará:
- Django
- djangorestframework
- gunicorn
- python-dateutil
- requests

---

## ⚙️ Paso 5: Configurar Django para Producción

### 5.1 Editar settings.py

```bash
nano visualizer/settings.py
```

Modificar las siguientes líneas:

```python
# Cambiar DEBUG a False
DEBUG = False

# Actualizar ALLOWED_HOSTS
ALLOWED_HOSTS = ['releases.life', 'www.releases.life', '79.137.76.235', 'localhost']

# Actualizar SECRET_KEY (generar una nueva)
SECRET_KEY = 'tu-clave-secreta-super-segura-aqui-cambiar'

# Actualizar CSRF_TRUSTED_ORIGINS
CSRF_TRUSTED_ORIGINS = [
    'https://releases.life',
    'https://www.releases.life',
    'http://releases.life',
    'http://www.releases.life'
]
```

**Guardar:** `Ctrl+O`, `Enter`, `Ctrl+X`

### 5.2 Generar SECRET_KEY segura (opcional pero recomendado)

```bash
python3 -c "from django.core.management.utils import get_random_secret_key; print(get_random_secret_key())"
```

Copia el resultado y úsalo en `SECRET_KEY`.

---

## 🗄️ Paso 6: Configurar Base de Datos

### 6.1 Crear migraciones y aplicarlas

```bash
cd /opt/artefacto-visualizer
source venv/bin/activate

python manage.py makemigrations
python manage.py makemigrations collector
python manage.py makemigrations xss_audit
python manage.py migrate
```

### 6.2 Crear superusuario (para acceder al admin)

```bash
python manage.py createsuperuser
```

Ingresa:
- Username: `admin`
- Email: `tu@email.com`
- Password: `tu-password-seguro`

### 6.3 Recolectar archivos estáticos

```bash
python manage.py collectstatic --noinput
```

Esto creará el directorio `staticfiles/` con todos los CSS, JS, etc.

---

## 🔧 Paso 7: Configurar Gunicorn

### 7.1 Crear archivo de configuración de Gunicorn

```bash
nano /opt/artefacto-visualizer/deploy/gunicorn_config.py
```

Contenido:

```python
import multiprocessing

# Directorios
bind = "127.0.0.1:8000"
workers = multiprocessing.cpu_count() * 2 + 1
worker_class = "sync"
worker_connections = 1000
max_requests = 1000
max_requests_jitter = 50
timeout = 300
keepalive = 5

# Logs
accesslog = "/var/log/gunicorn/access.log"
errorlog = "/var/log/gunicorn/error.log"
loglevel = "info"
access_log_format = '%(h)s %(l)s %(u)s %(t)s "%(r)s" %(s)s %(b)s "%(f)s" "%(a)s"'

# Process naming
proc_name = "artefacto-visualizer"

# Server mechanics
daemon = False
pidfile = "/var/run/gunicorn/artefacto-visualizer.pid"
user = "www-data"
group = "www-data"
```

**Guardar:** `Ctrl+O`, `Enter`, `Ctrl+X`

### 7.2 Crear directorios de logs y PID

```bash
mkdir -p /var/log/gunicorn
mkdir -p /var/run/gunicorn
chown -R www-data:www-data /var/log/gunicorn
chown -R www-data:www-data /var/run/gunicorn
```

### 7.3 Crear servicio systemd para Gunicorn

```bash
nano /etc/systemd/system/artefacto-visualizer.service
```

Contenido:

```ini
[Unit]
Description=Artefacto Visualizer Gunicorn daemon
After=network.target

[Service]
Type=notify
User=www-data
Group=www-data
RuntimeDirectory=gunicorn
WorkingDirectory=/opt/artefacto-visualizer
Environment="PATH=/opt/artefacto-visualizer/venv/bin"
ExecStart=/opt/artefacto-visualizer/venv/bin/gunicorn \
          --config /opt/artefacto-visualizer/deploy/gunicorn_config.py \
          visualizer.wsgi:application
ExecReload=/bin/kill -s HUP $MAINPID
KillMode=mixed
TimeoutStopSec=5
PrivateTmp=true
Restart=always

[Install]
WantedBy=multi-user.target
```

**Guardar:** `Ctrl+O`, `Enter`, `Ctrl+X`

### 7.4 Configurar permisos

```bash
chown -R www-data:www-data /opt/artefacto-visualizer
chmod 664 /opt/artefacto-visualizer/db.sqlite3
```

### 7.5 Habilitar e iniciar el servicio

```bash
systemctl daemon-reload
systemctl enable artefacto-visualizer
systemctl start artefacto-visualizer
systemctl status artefacto-visualizer
```

Deberías ver: `Active: active (running)`

---

## 🌐 Paso 8: Configurar Nginx

### 8.1 Crear configuración de Nginx

```bash
nano /etc/nginx/sites-available/artefacto-visualizer
```

Contenido:

```nginx
upstream artefacto_visualizer {
    server 127.0.0.1:8000 fail_timeout=0;
}

server {
    listen 80;
    server_name releases.life www.releases.life;
    
    client_max_body_size 100M;
    
    # Logs
    access_log /var/log/nginx/artefacto-access.log;
    error_log /var/log/nginx/artefacto-error.log;
    
    # Static files
    location /static/ {
        alias /opt/artefacto-visualizer/staticfiles/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
    
    # Media files
    location /media/ {
        alias /opt/artefacto-visualizer/media/;
        expires 30d;
    }
    
    # Proxy to Gunicorn
    location / {
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Host $http_host;
        proxy_redirect off;
        proxy_buffering off;
        
        # Timeouts para payloads grandes
        proxy_connect_timeout 300s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;
        
        proxy_pass http://artefacto_visualizer;
    }
}
```

**Guardar:** `Ctrl+O`, `Enter`, `Ctrl+X`

### 8.2 Habilitar el sitio

```bash
# Crear symlink
ln -s /etc/nginx/sites-available/artefacto-visualizer /etc/nginx/sites-enabled/

# Eliminar sitio default
rm -f /etc/nginx/sites-enabled/default

# Verificar configuración
nginx -t
```

Deberías ver: `syntax is ok` y `test is successful`

### 8.3 Reiniciar Nginx

```bash
systemctl restart nginx
systemctl enable nginx
systemctl status nginx
```

---

## 🔒 Paso 9: Configurar HTTPS con Let's Encrypt

### 9.1 Obtener certificado SSL

```bash
certbot --nginx -d releases.life -d www.releases.life
```

Responde:
- Email: `tu@email.com`
- Términos: `Y`
- Compartir email: `N` (opcional)
- Redirect HTTP to HTTPS: `2` (Sí, recomendado)

Certbot configurará automáticamente Nginx para HTTPS.

### 9.2 Verificar renovación automática

```bash
certbot renew --dry-run
```

Si todo está bien, verás: `Congratulations, all simulated renewals succeeded`

### 9.3 Verificar configuración final de Nginx

```bash
cat /etc/nginx/sites-available/artefacto-visualizer
```

Deberías ver bloques adicionales para el puerto 443 con SSL.

---

## ✅ Paso 10: Verificar el Despliegue

### 10.1 Verificar servicios

```bash
# Gunicorn
systemctl status artefacto-visualizer

# Nginx
systemctl status nginx

# Ver logs de Gunicorn
tail -f /var/log/gunicorn/error.log

# Ver logs de Nginx
tail -f /var/log/nginx/artefacto-error.log
```

### 10.2 Probar la aplicación

Abre en tu navegador:

- **HTTP:** http://releases.life/
- **HTTPS:** https://releases.life/
- **Dashboard:** https://releases.life/dashboard/
- **Admin:** https://releases.life/admin/

### 10.3 Probar el API endpoint

Desde tu máquina local:

```bash
curl -X POST https://releases.life/api/collect \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2024-01-09T10:00:00Z",
    "hostname": "test-machine",
    "public_ip": "1.2.3.4",
    "binary_size_bytes": 6500000
  }'
```

Deberías recibir: `{"status": "success", ...}`

---

## 🔄 Paso 11: Comandos de Mantenimiento

### Reiniciar servicios

```bash
# Reiniciar Gunicorn
systemctl restart artefacto-visualizer

# Reiniciar Nginx
systemctl restart nginx

# Reiniciar ambos
systemctl restart artefacto-visualizer nginx
```

### Ver logs en tiempo real

```bash
# Logs de Gunicorn
tail -f /var/log/gunicorn/error.log

# Logs de Nginx
tail -f /var/log/nginx/artefacto-error.log

# Logs de Django (si DEBUG=True)
tail -f /opt/artefacto-visualizer/debug.log
```

### Actualizar el código

```bash
# 1. Transferir nuevo código
scp artefacto-visualizer.tar.gz root@79.137.76.235:/tmp/

# 2. En el servidor
cd /opt/artefacto-visualizer
tar -xzf /tmp/artefacto-visualizer.tar.gz

# 3. Activar venv y actualizar dependencias
source venv/bin/activate
pip install -r requirements.txt

# 4. Aplicar migraciones
python manage.py migrate

# 5. Recolectar estáticos
python manage.py collectstatic --noinput

# 6. Reiniciar Gunicorn
systemctl restart artefacto-visualizer
```

### Backup de la base de datos

```bash
# Crear backup
cp /opt/artefacto-visualizer/db.sqlite3 \
   /opt/artefacto-visualizer/db.sqlite3.backup.$(date +%Y%m%d_%H%M%S)

# Restaurar backup
cp /opt/artefacto-visualizer/db.sqlite3.backup.YYYYMMDD_HHMMSS \
   /opt/artefacto-visualizer/db.sqlite3
systemctl restart artefacto-visualizer
```

---

## 🐛 Troubleshooting

### Problema: Gunicorn no inicia

```bash
# Ver logs detallados
journalctl -u artefacto-visualizer -n 50

# Verificar permisos
ls -la /opt/artefacto-visualizer/db.sqlite3
chown www-data:www-data /opt/artefacto-visualizer/db.sqlite3

# Probar manualmente
cd /opt/artefacto-visualizer
source venv/bin/activate
gunicorn --bind 127.0.0.1:8000 visualizer.wsgi:application
```

### Problema: 502 Bad Gateway

```bash
# Verificar que Gunicorn esté corriendo
systemctl status artefacto-visualizer

# Verificar que esté escuchando en el puerto correcto
netstat -tlnp | grep 8000

# Reiniciar servicios
systemctl restart artefacto-visualizer nginx
```

### Problema: Archivos estáticos no cargan

```bash
# Recolectar estáticos nuevamente
cd /opt/artefacto-visualizer
source venv/bin/activate
python manage.py collectstatic --noinput

# Verificar permisos
chown -R www-data:www-data /opt/artefacto-visualizer/staticfiles/

# Verificar configuración de Nginx
nginx -t
```

### Problema: Error de permisos en base de datos

```bash
# Arreglar permisos
chown www-data:www-data /opt/artefacto-visualizer/db.sqlite3
chmod 664 /opt/artefacto-visualizer/db.sqlite3
chown www-data:www-data /opt/artefacto-visualizer/
```

---

## 📊 Monitoreo

### Ver estadísticas de uso

```bash
# Conexiones activas
netstat -an | grep :80 | wc -l

# Uso de memoria
free -h

# Uso de disco
df -h

# Procesos de Gunicorn
ps aux | grep gunicorn
```

### Configurar monitoreo automático (opcional)

```bash
# Instalar htop para monitoreo interactivo
apt install htop

# Ejecutar
htop
```

---

## 🔐 Seguridad Adicional (Recomendado)

### Configurar autenticación HTTP Basic para el admin

```bash
# Instalar apache2-utils
apt install apache2-utils

# Crear archivo de contraseñas
htpasswd -c /etc/nginx/.htpasswd admin

# Editar configuración de Nginx
nano /etc/nginx/sites-available/artefacto-visualizer
```

Añadir dentro del bloque `location /admin/`:

```nginx
location /admin/ {
    auth_basic "Restricted Access";
    auth_basic_user_file /etc/nginx/.htpasswd;
    
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header Host $http_host;
    proxy_redirect off;
    proxy_pass http://artefacto_visualizer;
}
```

Reiniciar Nginx:

```bash
nginx -t
systemctl restart nginx
```

---

## ✅ Checklist Final

- [ ] Servidor actualizado
- [ ] Dependencias instaladas
- [ ] Proyecto transferido y descomprimido
- [ ] Entorno virtual creado
- [ ] Dependencias Python instaladas
- [ ] settings.py configurado (DEBUG=False, ALLOWED_HOSTS, SECRET_KEY)
- [ ] Migraciones aplicadas
- [ ] Superusuario creado
- [ ] Archivos estáticos recolectados
- [ ] Gunicorn configurado y corriendo
- [ ] Nginx configurado y corriendo
- [ ] SSL/HTTPS configurado con Let's Encrypt
- [ ] Firewall configurado
- [ ] Aplicación accesible en https://releases.life/
- [ ] API endpoint funcional
- [ ] Admin panel accesible

---

## 🎉 ¡Despliegue Completado!

Tu aplicación Artefacto Visualizer está ahora en producción en:

**URL Principal:** https://releases.life/  
**Dashboard XSS:** https://releases.life/dashboard/  
**Admin Panel:** https://releases.life/admin/  
**API Endpoint:** https://releases.life/api/collect

### Próximos pasos:

1. Configurar el agente para enviar datos a `https://releases.life/api/collect`
2. Monitorear logs regularmente
3. Configurar backups automáticos de la base de datos
4. Considerar usar PostgreSQL en lugar de SQLite para producción

---

**Fecha:** 2024-01-09  
**Versión:** 2.0  
**Servidor:** 79.137.76.235  
**Dominio:** releases.life
