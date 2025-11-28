# 🚀 Instrucciones Rápidas - Artefacto Visualizer

## Instalación y Configuración

### 1. Instalar Dependencias

```bash
cd visualizer
pip install -r requirements.txt
```

### 2. Configurar Base de Datos

```bash
python manage.py makemigrations
python manage.py migrate
```

### 3. Iniciar el Servidor

**Windows:**
```cmd
start_server.bat
```

**Linux/Mac:**
```bash
chmod +x start_server.sh
./start_server.sh
```

**Manual:**
```bash
python manage.py runserver 192.168.1.143:8080
```

## Configurar el Agente

Edita el archivo `artefacto/.env`:

```env
SERVER_URL=http://192.168.1.143:8080/api/collect
TIMEOUT=30s
```

## Uso

1. **Inicia el servidor Django** en `192.168.1.143:8080`
2. **Ejecuta el agente** desde la carpeta `artefacto/`
3. **Accede a la web** en http://192.168.1.143:8080/

## URLs Disponibles

- **Página principal**: http://192.168.1.143:8080/
  - Lista todas las ejecuciones del agente
  
- **Detalle de ejecución**: http://192.168.1.143:8080/execution/{GUID}/
  - Muestra toda la información de una ejecución específica
  
- **API Endpoint**: http://192.168.1.143:8080/api/collect
  - Recibe los datos del agente (POST)
  
- **Admin Panel**: http://192.168.1.143:8080/admin/
  - Panel de administración de Django (requiere superusuario)

## Crear Superusuario (Opcional)

Para acceder al panel de administración:

```bash
python manage.py createsuperuser
```

## Características

### 📊 Vista Principal
- Lista de todas las ejecuciones recibidas
- GUID único por ejecución
- Badges indicando si es VM, EDR detectado, etc.
- Ordenadas por fecha de recepción

### 🔍 Vista de Detalle
Cada ejecución muestra:

- **Información General**: GUID, hostname, timestamp
- **Detección de Sandbox**: Indicadores de VM, registro, disco
- **Sistema**: OS, CPU, RAM, procesos, conexiones de red
- **Hooks**: Funciones hooked, DLLs sospechosas
- **Crawler**: Archivos encontrados
- **EDR/AV**: Productos de seguridad detectados
- **Screenshot**: Si está disponible

### 🎨 Interfaz
- Tema oscuro moderno
- Desplegables para organizar información
- Tablas para datos estructurados
- Responsive y fácil de navegar

## Solución de Problemas

### El servidor no inicia
```bash
# Verifica que Django esté instalado
pip install Django

# Verifica la configuración
python manage.py check
```

### El agente no puede conectar
- Verifica que el servidor esté corriendo en `192.168.1.143:8080`
- Verifica que el firewall permita conexiones en el puerto 8080
- Verifica la URL en `artefacto/.env`

### No se muestran los datos
- Verifica que las migraciones estén aplicadas: `python manage.py migrate`
- Revisa los logs del servidor Django
- Verifica que el agente esté enviando a `/api/collect`

## Estructura de la Base de Datos

Cada ejecución genera:
- 1 registro en `AgentExecution` (con GUID único)
- 1 registro en `SandboxInfo`
- 1 registro en `SystemInfo`
- N registros en `ProcessInfo`
- N registros en `NetworkConnection`
- 1 registro en `HookInfo`
- N registros en `HookedFunction`
- 1 registro en `CrawlerInfo`
- 1 registro en `EDRInfo`
- N registros en `EDRProduct`

## Notas de Seguridad

⚠️ **Este servidor está configurado para desarrollo local**

- No usar en producción sin cambiar `SECRET_KEY`
- No exponer a Internet sin autenticación
- Cambiar `DEBUG = False` en producción
- Configurar HTTPS en producción
