# 🚀 Despliegue Rápido desde Windows

## Guía Express para Desplegar en releases.life

---

## ⚡ Opción 1: Script Automático (Recomendado)

### Paso 1: Abrir PowerShell

```powershell
# Presiona Win+X → Windows PowerShell (Admin)

# Navegar al directorio visualizer
cd "C:\Users\Usuario\OneDrive\Escritorio\TFE-ENIGMA\Sandbox-Fingerprinting\visualizer"
```

### Paso 2: Ejecutar Script

```powershell
# Ejecutar script de despliegue
.\deploy_from_windows.ps1
```

El script te pedirá:
1. **Contraseña del servidor** (varias veces)
2. **¿Configurar HTTPS?** → Responde `S`
3. **Email para Let's Encrypt** → Tu email
4. **¿Abrir navegador?** → Responde `S`

**¡Listo!** La aplicación estará en https://releases.life/

---

## 📝 Opción 2: Manual (Paso a Paso)

### Paso 1: Comprimir Proyecto

```powershell
# En PowerShell, desde el directorio visualizer
cd "C:\Users\Usuario\OneDrive\Escritorio\TFE-ENIGMA\Sandbox-Fingerprinting\visualizer"

# Comprimir
Compress-Archive -Path * -DestinationPath artefacto-visualizer.zip -Force
```

### Paso 2: Transferir al Servidor

```powershell
# Transferir
scp artefacto-visualizer.zip root@79.137.76.235:/tmp/

# Ingresa la contraseña cuando se solicite
```

### Paso 3: Conectar al Servidor

```powershell
# Conectar por SSH
ssh root@79.137.76.235

# Ingresa la contraseña
```

### Paso 4: Ejecutar Script de Instalación

**En el servidor Ubuntu:**

```bash
# Descomprimir
mkdir -p /opt/artefacto-visualizer
cd /opt/artefacto-visualizer
unzip /tmp/artefacto-visualizer.zip

# Ejecutar script de despliegue rápido
bash deploy/quick_deploy.sh
```

### Paso 5: Configurar HTTPS

```bash
# Configurar SSL
certbot --nginx -d releases.life -d www.releases.life

# Responder:
# Email: tu@email.com
# Términos: Y
# Redirect: 2
```

**¡Listo!** Abre https://releases.life/

---

## 🔧 Configuración Post-Despliegue

### Crear Superusuario

```bash
# En el servidor
cd /opt/artefacto-visualizer
source venv/bin/activate
python manage.py createsuperuser

# Username: admin
# Email: tu@email.com
# Password: TuPasswordSeguro123!
```

### Configurar Agente

**Editar `artefacto/.env` en Windows:**

```env
SERVER_URL=https://releases.life/api/collect
TIMEOUT=120s
```

**Compilar y probar:**

```powershell
cd ..\artefacto
go build -ldflags="-s -w" -trimpath -o agent.exe
.\agent.exe
```

---

## ✅ Verificar Funcionamiento

### Desde el Navegador

Abre en tu navegador:
- https://releases.life/
- https://releases.life/admin/ (usuario: admin)
- https://releases.life/dashboard/

### Desde PowerShell

```powershell
# Probar API
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

Deberías recibir: `{"status": "success", ...}`

---

## 🔄 Comandos Útiles

### Conectar al Servidor

```powershell
ssh root@79.137.76.235
```

### Reiniciar Servicios

```bash
# En el servidor
systemctl restart artefacto-visualizer
systemctl restart nginx
```

### Ver Logs

```bash
# Logs de Gunicorn
tail -f /var/log/gunicorn/error.log

# Logs de Nginx
tail -f /var/log/nginx/artefacto-error.log
```

### Actualizar Código

**Desde Windows:**

```powershell
# Comprimir nuevo código
cd visualizer
Compress-Archive -Path * -DestinationPath artefacto-visualizer.zip -Force

# Transferir
scp artefacto-visualizer.zip root@79.137.76.235:/tmp/
```

**En el servidor:**

```bash
cd /opt/artefacto-visualizer
unzip -o /tmp/artefacto-visualizer.zip
source venv/bin/activate
pip install -r requirements.txt
python manage.py migrate
python manage.py collectstatic --noinput
systemctl restart artefacto-visualizer
```

---

## 🐛 Solución de Problemas

### No puedo conectar por SSH

```powershell
# Verificar conectividad
Test-NetConnection -ComputerName 79.137.76.235 -Port 22
```

### 502 Bad Gateway

```bash
# En el servidor
systemctl status artefacto-visualizer
systemctl restart artefacto-visualizer nginx
```

### Archivos estáticos no cargan

```bash
# En el servidor
cd /opt/artefacto-visualizer
source venv/bin/activate
python manage.py collectstatic --noinput
systemctl restart nginx
```

---

## 📚 Documentación Completa

Para más detalles, consulta:
- `visualizer/deploy/DEPLOYMENT_GUIDE_WINDOWS.md` - Guía completa
- `DEPLOYMENT_COMMANDS.md` - Todos los comandos
- `visualizer/deploy/quick_deploy.sh` - Script del servidor

---

## 🎉 ¡Listo!

Tu aplicación está en producción en:

**🌐 URL:** https://releases.life/  
**📊 Dashboard:** https://releases.life/dashboard/  
**🔐 Admin:** https://releases.life/admin/  
**🔌 API:** https://releases.life/api/collect

---

**Fecha:** 2024-01-09  
**Tiempo estimado:** 10-15 minutos  
**Dificultad:** Fácil  
**Estado:** ✅ Listo para usar
