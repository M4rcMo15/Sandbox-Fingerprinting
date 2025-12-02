# 🔐 Configurar Autenticación - Guía Rápida

## 📊 Persistencia de Datos (Resumen)

**Base de datos:** SQLite3 en `/opt/artefacto-visualizer/db.sqlite3`
- ✅ Los datos **NO se borran** al reiniciar el servidor
- ✅ Persisten indefinidamente hasta que los elimines manualmente
- ✅ Crecen con cada ejecución del agente
- ⚠️ Recomendado: Configurar backups automáticos

---

## 🔐 Configurar Autenticación (3 Pasos)

### Paso 1: Copiar Archivos al Servidor

Desde tu PC Windows:

```powershell
# Comprimir visualizer actualizado
cd C:\Users\Usuario\OneDrive\Escritorio\TFE-ENIGMA\Sandbox-Fingerprinting
tar -czf visualizer-auth.tar.gz visualizer/deploy/

# Copiar al servidor
scp visualizer-auth.tar.gz root@54.37.226.179:/tmp/
```

### Paso 2: En el Servidor

```bash
# Conectar
ssh root@54.37.226.179

# Extraer archivos
cd /tmp
tar -xzf visualizer-auth.tar.gz

# Copiar scripts al directorio correcto
sudo cp -r visualizer/deploy/* /opt/artefacto-visualizer/deploy/

# Dar permisos de ejecución
cd /opt/artefacto-visualizer/deploy
sudo chmod +x setup_auth.sh update_nginx_auth.sh

# Ejecutar configuración
sudo ./setup_auth.sh
```

Te pedirá la contraseña para `marc.monfort`. Introdúcela dos veces.

### Paso 3: Aplicar Configuración

```bash
# Aplicar configuración de Nginx
sudo ./update_nginx_auth.sh

# Verificar
sudo systemctl status nginx
```

---

## ✅ Verificar

### Desde el Navegador

1. Ir a: http://54.37.226.179
2. Debería aparecer un popup pidiendo usuario/contraseña
3. Introducir:
   - **Usuario:** `marc.monfort`
   - **Contraseña:** [la que configuraste]

### Desde cURL

```bash
# Con autenticación (debería funcionar)
curl -u marc.monfort:tu_contraseña http://54.37.226.179

# Sin autenticación (debería pedir login)
curl http://54.37.226.179

# API sin autenticación (debería funcionar)
curl http://54.37.226.179/api/collect
```

---

## 🔓 Rutas Protegidas

### CON Autenticación (requiere login)
- ✅ `/` - Página principal
- ✅ `/statistics/` - Estadísticas
- ✅ `/execution/{guid}/` - Detalle de ejecuciones
- ✅ `/admin/` - Panel de administración

### SIN Autenticación (acceso público)
- 🔓 `/api/collect` - Endpoint para el agente
- 🔓 `/static/` - Archivos estáticos

**Importante:** El endpoint `/api/collect` NO tiene autenticación para que el agente pueda enviar datos sin problemas.

---

## 👥 Gestionar Usuarios

### Agregar Usuario Adicional

```bash
ssh root@54.37.226.179
sudo htpasswd /etc/nginx/auth/.htpasswd nuevo.usuario
sudo systemctl restart nginx
```

### Cambiar Contraseña

```bash
ssh root@54.37.226.179
sudo htpasswd /etc/nginx/auth/.htpasswd marc.monfort
sudo systemctl restart nginx
```

### Ver Usuarios

```bash
ssh root@54.37.226.179
sudo cat /etc/nginx/auth/.htpasswd
```

---

## 💾 Configurar Backups Automáticos

```bash
ssh root@54.37.226.179

# Crear directorio de backups
sudo mkdir -p /opt/artefacto-visualizer/backups

# Editar crontab
sudo crontab -e

# Agregar backup diario a las 2 AM
0 2 * * * cp /opt/artefacto-visualizer/db.sqlite3 /opt/artefacto-visualizer/backups/db_$(date +\%Y\%m\%d).sqlite3

# Agregar limpieza de backups antiguos (mantener últimos 30 días)
0 3 * * * find /opt/artefacto-visualizer/backups -name "db_*.sqlite3" -mtime +30 -delete
```

---

## 🔒 Configurar HTTPS (Recomendado)

```bash
ssh root@54.37.226.179

# Instalar Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obtener certificado SSL
sudo certbot --nginx -d 54.37.226.179

# Verificar renovación automática
sudo certbot renew --dry-run
```

Después de configurar HTTPS:
- URL: https://54.37.226.179
- Autenticación + Cifrado SSL ✅

---

## 🐛 Solución de Problemas

### No aparece el popup de autenticación

```bash
# Verificar que existe el archivo de contraseñas
ssh root@54.37.226.179
sudo ls -la /etc/nginx/auth/.htpasswd

# Verificar configuración de Nginx
sudo nginx -t

# Ver logs
sudo tail -f /var/log/nginx/artefacto-visualizer-error.log
```

### El agente no puede enviar datos

```bash
# Verificar que /api/ NO tiene autenticación
ssh root@54.37.226.179
sudo cat /etc/nginx/sites-available/artefacto-visualizer | grep -A 5 "location /api/"

# Debería mostrar que NO tiene auth_basic
```

### Olvidé la contraseña

```bash
ssh root@54.37.226.179
sudo htpasswd /etc/nginx/auth/.htpasswd marc.monfort
sudo systemctl restart nginx
```

---

## 📋 Resumen

### Persistencia
- ✅ SQLite3: `/opt/artefacto-visualizer/db.sqlite3`
- ✅ Datos permanentes (no se borran)
- ⚠️ Configurar backups automáticos

### Autenticación
- ✅ Usuario: `marc.monfort`
- ✅ Contraseña: La que configures
- ✅ Páginas protegidas
- 🔓 API sin protección

### Seguridad
- ✅ HTTP Basic Authentication
- ⚠️ Configurar HTTPS
- ⚠️ Backups regulares

---

**Servidor:** http://54.37.226.179  
**Usuario:** marc.monfort  
**Configuración:** 2 de diciembre de 2024
