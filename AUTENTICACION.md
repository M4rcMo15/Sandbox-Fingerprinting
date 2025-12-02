# 🔐 Autenticación del Visualizer

## 📊 Persistencia de Datos

### Almacenamiento

**Base de datos:** SQLite3
- **Ubicación:** `/opt/artefacto-visualizer/db.sqlite3`
- **Tipo:** Archivo único persistente
- **Tamaño actual:** Crece con cada ejecución
- **Persistencia:** ✅ Los datos NO se borran al reiniciar el servidor

### Estructura de Datos

```
db.sqlite3
├── collector_agentexecution      # Ejecuciones principales
├── collector_sandboxinfo          # Información de sandbox
├── collector_systeminfo           # Información del sistema
├── collector_processinfo          # Procesos
├── collector_networkconnection    # Conexiones de red
├── collector_hookinfo             # Hooks detectados
├── collector_hookedfunction       # Funciones hooked
├── collector_crawlerinfo          # Archivos encontrados
├── collector_edrinfo              # EDR/AV detectados
├── collector_edrproduct           # Productos EDR
├── collector_geolocation          # Geolocalización
└── collector_toolsinfo            # Herramientas detectadas
```

### Backups

**Automático:** No configurado por defecto

**Manual:**
```bash
# Conectar al servidor
ssh root@54.37.226.179

# Hacer backup
sudo cp /opt/artefacto-visualizer/db.sqlite3 \
       /opt/artefacto-visualizer/backups/db_$(date +%Y%m%d_%H%M%S).sqlite3
```

**Configurar backup automático (cron):**
```bash
# Editar crontab
sudo crontab -e

# Agregar backup diario a las 2 AM
0 2 * * * cp /opt/artefacto-visualizer/db.sqlite3 /opt/artefacto-visualizer/backups/db_$(date +\%Y\%m\%d).sqlite3
```

---

## 🔐 Configurar Autenticación

### Paso 1: Conectar al Servidor

```bash
ssh root@54.37.226.179
cd /opt/artefacto-visualizer
```

### Paso 2: Ejecutar Script de Configuración

```bash
# Dar permisos de ejecución
chmod +x deploy/setup_auth.sh
chmod +x deploy/update_nginx_auth.sh

# Ejecutar configuración
sudo ./deploy/setup_auth.sh
```

Te pedirá la contraseña para el usuario `marc.monfort`. Introdúcela dos veces.

### Paso 3: Aplicar Configuración de Nginx

```bash
sudo ./deploy/update_nginx_auth.sh
```

### Paso 4: Verificar

```bash
# Ver estado de Nginx
sudo systemctl status nginx

# Probar acceso (debería pedir usuario/contraseña)
curl http://54.37.226.179
```

---

## 🔓 Configuración de Autenticación

### Usuario Configurado

- **Usuario:** `marc.monfort`
- **Contraseña:** La que configures en el paso 2

### Rutas Protegidas

✅ **CON autenticación (requiere login):**
- `/` - Página principal
- `/statistics/` - Estadísticas
- `/execution/{guid}/` - Detalle de ejecuciones
- `/admin/` - Panel de administración Django

🔓 **SIN autenticación (acceso público):**
- `/api/collect` - Endpoint para el agente
- `/static/` - Archivos estáticos
- `/media/` - Archivos media

### ¿Por qué el API no tiene autenticación?

El endpoint `/api/collect` NO requiere autenticación para que el agente pueda enviar datos sin problemas. Solo las páginas de visualización están protegidas.

---

## 👥 Gestionar Usuarios

### Agregar Nuevo Usuario

```bash
# Agregar usuario adicional
sudo htpasswd /etc/nginx/auth/.htpasswd nombre.usuario

# Reiniciar Nginx
sudo systemctl restart nginx
```

### Cambiar Contraseña

```bash
# Cambiar contraseña de marc.monfort
sudo htpasswd /etc/nginx/auth/.htpasswd marc.monfort

# Reiniciar Nginx
sudo systemctl restart nginx
```

### Eliminar Usuario

```bash
# Eliminar usuario
sudo htpasswd -D /etc/nginx/auth/.htpasswd nombre.usuario

# Reiniciar Nginx
sudo systemctl restart nginx
```

### Ver Usuarios

```bash
# Listar usuarios configurados
sudo cat /etc/nginx/auth/.htpasswd
```

---

## 🌐 Acceder al Visualizer

### Desde el Navegador

1. Ir a: http://54.37.226.179
2. Aparecerá un popup pidiendo credenciales
3. Introducir:
   - **Usuario:** `marc.monfort`
   - **Contraseña:** [tu contraseña]
4. Click en "Iniciar sesión"

### Desde cURL

```bash
# Con autenticación
curl -u marc.monfort:tu_contraseña http://54.37.226.179

# Sin autenticación (solo API)
curl http://54.37.226.179/api/collect
```

---

## 🔒 Seguridad Adicional

### Configurar HTTPS (Recomendado)

```bash
# Instalar Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obtener certificado SSL
sudo certbot --nginx -d 54.37.226.179

# Renovación automática ya está configurada
```

Después de configurar HTTPS:
- URL: https://54.37.226.179
- Autenticación + Cifrado SSL

### Cambiar a PostgreSQL (Opcional)

Para producción pesada, considera cambiar de SQLite a PostgreSQL:

```bash
# Instalar PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# Configurar base de datos
sudo -u postgres createdb artefacto_visualizer
sudo -u postgres createuser artefacto_user -P

# Actualizar settings_production.py
# DATABASES = {
#     'default': {
#         'ENGINE': 'django.db.backends.postgresql',
#         'NAME': 'artefacto_visualizer',
#         'USER': 'artefacto_user',
#         'PASSWORD': 'tu_contraseña',
#         'HOST': 'localhost',
#         'PORT': '5432',
#     }
# }
```

---

## 🐛 Solución de Problemas

### No puedo acceder (pide contraseña pero no funciona)

```bash
# Verificar que existe el archivo de contraseñas
sudo ls -la /etc/nginx/auth/.htpasswd

# Verificar configuración de Nginx
sudo nginx -t

# Ver logs de Nginx
sudo tail -f /var/log/nginx/artefacto-visualizer-error.log
```

### El agente no puede enviar datos

```bash
# Verificar que /api/ NO tiene autenticación
sudo cat /etc/nginx/sites-available/artefacto-visualizer | grep -A 10 "location /api/"

# Debería mostrar que NO tiene auth_basic
```

### Olvidé la contraseña

```bash
# Resetear contraseña
sudo htpasswd /etc/nginx/auth/.htpasswd marc.monfort

# Reiniciar Nginx
sudo systemctl restart nginx
```

### Quitar autenticación

```bash
# Restaurar configuración sin autenticación
sudo cp /opt/artefacto-visualizer/deploy/nginx.conf \
       /etc/nginx/sites-available/artefacto-visualizer

# Reiniciar Nginx
sudo nginx -t
sudo systemctl restart nginx
```

---

## 📋 Resumen

### Persistencia
- ✅ Datos almacenados en SQLite3
- ✅ Ubicación: `/opt/artefacto-visualizer/db.sqlite3`
- ✅ Persistentes (no se borran al reiniciar)
- ⚠️ Configurar backups automáticos

### Autenticación
- ✅ Usuario: `marc.monfort`
- ✅ Contraseña: La que configures
- ✅ Páginas protegidas: /, /statistics/, /execution/
- 🔓 API sin protección: /api/collect

### Seguridad
- ✅ HTTP Basic Authentication
- ⚠️ Configurar HTTPS (recomendado)
- ⚠️ Backups regulares
- ⚠️ Actualizar contraseñas periódicamente

---

**Configuración:** 2 de diciembre de 2024  
**Servidor:** http://54.37.226.179  
**Usuario:** marc.monfort
