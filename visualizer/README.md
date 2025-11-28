# Artefacto Visualizer

Aplicación Django para visualizar los datos de fingerprinting recopilados por el agente Artefacto.

## Instalación

1. Instalar dependencias:
```bash
pip install -r requirements.txt
```

2. Crear la base de datos:
```bash
python manage.py makemigrations
python manage.py migrate
```

3. (Opcional) Crear superusuario para el admin:
```bash
python manage.py createsuperuser
```

## Uso

### Iniciar el servidor en 192.168.1.143:8080

```bash
python manage.py runserver 192.168.1.143:8080
```

### Acceder a la aplicación

- **Interfaz web**: http://192.168.1.143:8080/
- **Panel de administración**: http://192.168.1.143:8080/admin/
- **Endpoint API**: http://192.168.1.143:8080/api/collect (POST)

## Configurar el agente

Actualiza la configuración del agente en `artefacto/.env`:

```
SERVER_URL=http://192.168.1.143:8080/api/collect
TIMEOUT=30s
```

## Características

- ✅ Recepción automática de datos del agente
- ✅ GUID único por cada ejecución
- ✅ Visualización completa de todos los datos recopilados
- ✅ Desplegables para organizar la información
- ✅ Interfaz oscura y moderna
- ✅ Tablas para procesos, conexiones de red, hooks, etc.
- ✅ Visualización de screenshots en base64
- ✅ Detección de VM, EDR/AV, hooks y más

## Estructura de datos

Cada ejecución del agente genera:
- **GUID único**: Identificador de la ejecución
- **Sandbox Info**: Detección de VM y sandbox
- **System Info**: Información del sistema, procesos, red, etc.
- **Hook Info**: Funciones hooked y DLLs sospechosas
- **Crawler Info**: Archivos encontrados
- **EDR Info**: Productos de seguridad detectados
