# 🚀 Guía de Despliegue en Producción (Windows → Ubuntu Server)

## Despliegue de Artefacto Visualizer desde Windows

**Máquina Local:** Windows 10/11  
**Servidor:** Ubuntu Server 20.04/22.04  
**IP Servidor:** 79.137.76.235  
**Dominio:** releases.life  
**Usuario:** root

---

## 📋 Requisitos Previos

### En tu máquina Windows:
- PowerShell o CMD
- Cliente SSH (OpenSSH viene con Windows 10/11)
- WinSCP o FileZilla (opcional, para transferencia gráfica)
- 7-Zip o WinRAR (para comprimir)

### Verificar SSH en Windows:
```powershell
ssh -V
```

Si no está instalado:
1. Configuración → Aplicaciones → Características opcionales
2. Agregar "Cliente OpenSSH"

---

## 🔧 PARTE 1: Preparar el Proyecto en Windows

### Paso 1.1: Abrir PowerShell

```powershell
# Abrir PowerShell como Administrador
# Presiona Win+X → Windows PowerShell (Admin)

# Navegar al directorio del proyecto
cd "C:\Users\Usuario\OneDrive\Escritorio\TFE-ENIGMA\Sandbox-Fingerprinting\visualizer"
```

### Paso 1.2: Comprimir el Proyecto

**Opción A: Con PowerShell (Recomendado)**

```powershell
# Crear archivo ZIP excluyendo archivos innecesarios
$exclude = @('*.pyc', '__pycache__', '*.sqlite3', 'staticfiles', 'media', '.git', 'venv')

# Comprimir
Compress-Archive -Path * -DestinationPath artefacto-visualizer.zip -Force

# Verificar que se creó
ls artefacto-visualizer.zip
```

**Opción B: Con 7-Zip (si lo tienes instalado)**

```powershell
# Usando 7-Zip desde línea de comandos
& "C:\Program Files\7-Zip\7z.exe" a -tzip artefacto-visualizer.zip * -xr!*.pyc -xr!__pycache__ -xr!*.sqlite3 -xr!staticfiles -xr!media -xr!.git -xr!venv
```

**Opción C: Manualmente con el Explorador**

1. Selecciona todos los archivos EXCEPTO:
   - `__pycache__` (carpetas)
   - `*.pyc` (archivos)
   - `db.sqlite3`
   - `staticfiles/`
   - `media/`
   - `venv/`
2. Click derecho → Enviar a → Carpeta comprimida (zip)
3. Renombrar a `artefacto-visualizer.zip`

---

## 🌐 PARTE 2: Transferir al Servidor

### Paso 2.1: Transferir con SCP (PowerShell)

```powershell
# Transferir el archivo ZIP al servidor
scp artefacto-visualizer.zip root@79.137.76.235:/tmp/

# Te pedirá la contraseña del servidor
# Escribe la contraseña y presiona Enter
```

**Nota:** Si es la primera vez, te preguntará si confías en el servidor. Escribe `yes` y Enter.

### Paso 2.2: Alternativa - Transferir con WinSCP (GUI)

Si prefieres una interfaz gráfica:

1. **Descargar WinSCP:** https://winscp.net/
2. **Abrir WinSCP**
3. **Nueva Sesión:**
   - Protocolo: SFTP
   - Host: 79.137.76.235
   - Puerto: 22
   - Usuario: root
   - Contraseña: [tu contraseña]
4. **Conectar**
5. **Arrastrar** `artefacto-visualizer.zip` a `/tmp/`

---

## 🖥️ PARTE 3: Configurar el Servidor Ubuntu

### Paso 3.1: Conectar por SSH

```powershell
# Desde PowerShell
ssh root@79.137.76.235

# Ingresa la contraseña cuando se solicite
```

**A partir de aquí, todos los comandos se ejecutan en el servidor Ubuntu.**

---

### Paso 3.2: Actualizar el Sistema

```bash
apt update && apt upgrade -y
```

### Paso 3.3: Instalar Dependencias

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
  ufw \
  unzip
```

### Paso 3.4: Configurar Firewall

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

Deberías ver:
```
Status: active

To                         Action      From
--                         ------      ----
22/tcp                     ALLOW       Anywhere
80/tcp                     ALLOW       Anywhere
443/tcp                     ALLOW       Anywhere
```

---

## 📦 PARTE 4: Descomprimir y Configurar Proyecto

### Paso 4.1: Descomprimir

```bash
# Crear directorio del proyecto
mkdir -p /opt/artefacto-visualizer

# Descomprimir
cd /opt/artefacto-visualizer
unzip /tmp/artefacto-visualizer.zip

# Limpiar
rm /tmp/artefacto-visualizer.zip

# Verificar archivos
ls -la
```

Deberías ver: `manage.py`, `requirements.txt`, `collector/`, `xss_audit/`, etc.

### Paso 4.2: Crear Entorno Virtual

```bash
cd /opt/artefacto-visualizer
python3 -m venv venv
```

### Paso 4.3: Instalar Dependencias Python

```bash
source venv/bin/activate
pip install --upgrade pip
pip install -r requirements.txt
```

Esto instalará Django, gunicorn, requests, etc.

---

## ⚙️ PARTE 5: Configurar Django para Producción

### Paso 5.1: Editar settings.py

```bash
nano visualizer/settings.py
```

**Buscar y modificar estas líneas:**

```python
# Línea ~6: Cambiar DEBUG
DEBUG = False

# Línea ~8: Actualizar ALLOWED_HOSTS
ALLOWED_HOSTS = ['releases.life', 'www.releases.life', '79.137.76.235', 'localhost']

# Línea ~4: Generar nueva SECRET_KEY
SECRET_KEY = 'PEGAR-AQUI-LA-CLAVE-GENERADA'

# Línea ~90: Actualizar CSRF_TRUSTED_ORIGINS
CSRF_TRUSTED_ORIGINS = [
    'https://releases.life',
    'https://www.releases.life',
    'http://releases.life',
    'http://www.releases.life'
]
```

**Guardar:** `Ctrl+O`, `Enter`, `Ctrl+X`

### Paso 5.2: Generar SECRET_KEY

```bash
python3 -c "from django.core.management.utils import get_random_secret_key; print(get_random_secret_key())"
```

Copia el resultado y pégalo en `SECRET_KEY` del paso anterior.

---

## 🗄️ PARTE 6: Configurar Base de Datos

### Paso 6.1: Aplicar Migraciones

```bash
cd /opt/artefacto-visualizer
source venv/bin/activate

python manage.py makemigrations
python manage.py makemigrations collector
python manage.py makemigrations xss_audit
python manage.py migrate
```

### Paso 6.2: Crear Superusuario

```bash
python manage.py createsuperuser
```

Ingresa:
- **Username:** `admin`
- **Email:** `tu@email.com`
- **Password:** `TuPasswordSeguro123!`
- **Password (again):** `TuPasswordSeguro123!`

### Paso 6.3: Recolectar Archivos Estáticos

```bash
python manage.py collectstatic --noinput
```

### Paso 6.4: Configurar Permisos

```bash
chown -R www-data:www-data /opt/artefacto-visualizer
chmod 664 /opt/artefacto-visualizer/db.sqlite3
```

---

## 🔧 PARTE 7: Configurar Gunicorn

### Paso 7.1: Crear Directorios de Logs

```bash
mkdir -p /var/log/gunicorn
mkdir -p /var/run/gunicorn
chown -R www-data:www-data /var/log/gunicorn
chown -R www-data:www-data /var/run/gunicorn
```

### Paso 7.2: Crear Servicio Systemd

```bash
nano /etc/systemd/system/artefacto-visualizer.service
```

**Pegar este contenido:**

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

### Paso 7.3: Habilitar e Iniciar Gunicorn

```bash
systemctl daemon-reload
systemctl enable artefacto-visualizer
systemctl start artefacto-visualizer
systemctl status artefacto-visualizer
```

Deberías ver: `Active: active (running)` en verde.

**Si hay error:**
```bash
journalctl -u artefacto-visualizer -n 50
```

---

## 🌐 PARTE 8: Configurar Nginx

### Paso 8.1: Copiar Configuración

```bash
cp /opt/artefacto-visualizer/deploy/nginx.conf /etc/nginx/sites-available/artefacto-visualizer
```

### Paso 8.2: Habilitar Sitio

```bash
# Crear symlink
ln -s /etc/nginx/sites-available/artefacto-visualizer /etc/nginx/sites-enabled/

# Eliminar sitio default
rm -f /etc/nginx/sites-enabled/default

# Verificar configuración
nginx -t
```

Deberías ver:
```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

### Paso 8.3: Reiniciar Nginx

```bash
systemctl restart nginx
systemctl enable nginx
systemctl status nginx
```

Deberías ver: `Active: active (running)` en verde.

---

## 🔒 PARTE 9: Configurar HTTPS con Let's Encrypt

### Paso 9.1: Obtener Certificado SSL

```bash
certbot --nginx -d releases.life -d www.releases.life
```

**Responder las preguntas:**

1. **Email:** `tu@email.com` (para notificaciones de renovación)
2. **Términos de servicio:** `Y` (Sí, acepto)
3. **Compartir email con EFF:** `N` (No, opcional)
4. **Redirect HTTP to HTTPS:** `2` (Sí, redirigir automáticamente)

Deberías ver:
```
Congratulations! You have successfully enabled HTTPS on https://releases.life and https://www.releases.life
```

### Paso 9.2: Verificar Renovación Automática

```bash
certbot renew --dry-run
```

Deberías ver: `Congratulations, all simulated renewals succeeded`

---

## ✅ PARTE 10: Verificar el Despliegue

### Paso 10.1: Verificar Servicios en el Servidor

```bash
# Estado de Gunicorn
systemctl status artefacto-visualizer

# Estado de Nginx
systemctl status nginx

# Ver logs de Gunicorn
tail -f /var/log/gunicorn/error.log

# Ver logs de Nginx
tail -f /var/log/nginx/artefacto-error.log
```

### Paso 10.2: Probar desde Windows

**Abrir navegador en Windows:**

- **HTTP:** http://releases.life/
- **HTTPS:** https://releases.life/
- **Dashboard:** https://releases.life/dashboard/
- **Admin:** https://releases.life/admin/
- **Estadísticas:** https://releases.life/statistics/

### Paso 10.3: Probar API desde PowerShell (Windows)

```powershell
# Probar el endpoint API
$body = @{
    timestamp = "2024-01-09T10:00:00Z"
    hostname = "test-windows"
    public_ip = "1.2.3.4"
    binary_size_bytes = 6500000
} | ConvertTo-Json

Invoke-RestMethod -Uri "https://releases.life/api/collect" `
                  -Method Post `
                  -Body $body `
                  -ContentType "application/json"
```

Deberías recibir:
```json
{
  "status": "success",
  "execution_id": "...",
  "message": "Data processed and analyzed successfully"
}
```

---

## 🔄 PARTE 11: Comandos de Mantenimiento

### Desde Windows - Conectar al Servidor

```powershell
# Conectar por SSH
ssh root@79.137.76.235
```

### En el Servidor - Reiniciar Servicios

```bash
# Reiniciar Gunicorn
systemctl restart artefacto-visualizer

# Reiniciar Nginx
systemctl restart nginx

# Reiniciar ambos
systemctl restart artefacto-visualizer nginx
```

### Ver Logs en Tiempo Real

```bash
# Logs de Gunicorn
tail -f /var/log/gunicorn/error.log

# Logs de Nginx
tail -f /var/log/nginx/artefacto-error.log

# Logs de systemd
journalctl -u artefacto-visualizer -f
```

### Actualizar el Código

**Desde Windows:**

```powershell
# 1. Comprimir nuevo código
cd "C:\Users\Usuario\OneDrive\Escritorio\TFE-ENIGMA\Sandbox-Fingerprinting\visualizer"
Compress-Archive -Path * -DestinationPath artefacto-visualizer.zip -Force

# 2. Transferir al servidor
scp artefacto-visualizer.zip root@79.137.76.235:/tmp/
```

**En el servidor:**

```bash
# 3. Descomprimir
cd /opt/artefacto-visualizer
unzip -o /tmp/artefacto-visualizer.zip
rm /tmp/artefacto-visualizer.zip

# 4. Actualizar dependencias
source venv/bin/activate
pip install -r requirements.txt

# 5. Aplicar migraciones
python manage.py migrate

# 6. Recolectar estáticos
python manage.py collectstatic --noinput

# 7. Reiniciar
systemctl restart artefacto-visualizer
```

### Backup de Base de Datos

```bash
# Crear backup
cp /opt/artefacto-visualizer/db.sqlite3 \
   /opt/artefacto-visualizer/db.sqlite3.backup.$(date +%Y%m%d_%H%M%S)

# Listar backups
ls -lh /opt/artefacto-visualizer/*.backup.*

# Restaurar backup
cp /opt/artefacto-visualizer/db.sqlite3.backup.YYYYMMDD_HHMMSS \
   /opt/artefacto-visualizer/db.sqlite3
systemctl restart artefacto-visualizer
```

### Descargar Backup a Windows

```powershell
# Desde PowerShell en Windows
scp root@79.137.76.235:/opt/artefacto-visualizer/db.sqlite3 ./db.sqlite3.backup
```

---

## 🐛 Troubleshooting

### Problema: No puedo conectar por SSH

**Desde PowerShell:**

```powershell
# Verificar conectividad
Test-NetConnection -ComputerName 79.137.76.235 -Port 22

# Probar SSH con verbose
ssh -v root@79.137.76.235
```

### Problema: Gunicorn no inicia

```bash
# Ver logs detallados
journalctl -u artefacto-visualizer -n 100

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

# Verificar puerto
netstat -tlnp | grep 8000

# Reiniciar servicios
systemctl restart artefacto-visualizer nginx
```

### Problema: Archivos estáticos no cargan

```bash
# Recolectar estáticos
cd /opt/artefacto-visualizer
source venv/bin/activate
python manage.py collectstatic --noinput

# Verificar permisos
chown -R www-data:www-data /opt/artefacto-visualizer/staticfiles/

# Reiniciar Nginx
systemctl restart nginx
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

### Ver Estadísticas del Servidor

```bash
# Uso de memoria
free -h

# Uso de disco
df -h

# Procesos de Gunicorn
ps aux | grep gunicorn

# Conexiones activas
netstat -an | grep :80 | wc -l
netstat -an | grep :443 | wc -l
```

### Instalar Herramientas de Monitoreo

```bash
# Instalar htop
apt install htop

# Ejecutar
htop
```

---

## 🔐 Seguridad Adicional (Recomendado)

### Cambiar Puerto SSH (Opcional)

```bash
# Editar configuración SSH
nano /etc/ssh/sshd_config

# Cambiar línea:
# Port 22
# Por:
Port 2222

# Guardar y reiniciar
systemctl restart sshd

# Actualizar firewall
ufw allow 2222/tcp
ufw delete allow 22/tcp
```

**Conectar desde Windows:**
```powershell
ssh -p 2222 root@79.137.76.235
```

### Deshabilitar Login Root (Recomendado)

```bash
# Crear usuario normal primero
adduser deploy
usermod -aG sudo deploy

# Editar SSH config
nano /etc/ssh/sshd_config

# Cambiar:
PermitRootLogin no

# Reiniciar SSH
systemctl restart sshd
```

### Configurar Fail2Ban

```bash
# Instalar
apt install fail2ban

# Configurar
cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
nano /etc/fail2ban/jail.local

# Habilitar
systemctl enable fail2ban
systemctl start fail2ban
```

---

## ✅ Checklist Final

- [ ] Proyecto comprimido en Windows
- [ ] Archivo transferido al servidor
- [ ] Dependencias del sistema instaladas
- [ ] Proyecto descomprimido en `/opt/artefacto-visualizer`
- [ ] Entorno virtual creado
- [ ] Dependencias Python instaladas
- [ ] `settings.py` configurado (DEBUG=False, ALLOWED_HOSTS, SECRET_KEY)
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

Tu aplicación está ahora en producción:

**🌐 URLs:**
- Principal: https://releases.life/
- Dashboard: https://releases.life/dashboard/
- Estadísticas: https://releases.life/statistics/
- Admin: https://releases.life/admin/
- API: https://releases.life/api/collect

### Configurar el Agente

**Editar `artefacto/.env`:**

```env
SERVER_URL=https://releases.life/api/collect
TIMEOUT=120s
```

**Compilar y probar:**

```powershell
cd artefacto
go build -ldflags="-s -w" -trimpath -o agent.exe
.\agent.exe
```

---

## 📞 Soporte

### Comandos Rápidos desde Windows

```powershell
# Conectar al servidor
ssh root@79.137.76.235

# Transferir archivo
scp archivo.zip root@79.137.76.235:/tmp/

# Descargar archivo
scp root@79.137.76.235:/path/to/file ./

# Ver logs remotos
ssh root@79.137.76.235 "tail -f /var/log/gunicorn/error.log"

# Reiniciar servicios remotos
ssh root@79.137.76.235 "systemctl restart artefacto-visualizer"
```

---

**Fecha:** 2024-01-09  
**Versión:** 2.0  
**Plataforma Local:** Windows 10/11  
**Servidor:** Ubuntu Server 20.04/22.04  
**IP:** 79.137.76.235  
**Dominio:** releases.life  
**Estado:** ✅ Listo para desplegar
