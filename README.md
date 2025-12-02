# 🔍 Artefacto - Sandbox Detection & Fingerprinting Tool

Sistema completo de detección de sandbox y fingerprinting para operaciones de Red Team, con visualización web de datos recopilados.

## �  Descripción

**Artefacto** es una herramienta de Red Team que consta de dos componentes principales:

1. **Agente (Go)**: Ejecutable que recopila información detallada del sistema objetivo
2. **Visualizer (Django)**: Servidor web para visualizar y analizar los datos recopilados

## 🎯 Características

### Agente de Recopilación

- ✅ Detección de sandbox y máquinas virtuales
- ✅ Detección de EDR/AV (Windows Defender, CrowdStrike, SentinelOne, etc.)
- ✅ Detección de herramientas de análisis (IDA Pro, x64dbg, Wireshark, etc.)
- ✅ Información completa del sistema (OS, CPU, RAM, procesos, conexiones)
- ✅ Detección de hooks en funciones críticas
- ✅ Crawler de archivos sensibles
- ✅ Screenshot del escritorio
- ✅ Geolocalización por IP
- ✅ Exfiltración automática de datos

### Visualizer Web

- ✅ Dashboard con lista de ejecuciones
- ✅ Vista detallada de cada ejecución
- ✅ Página de estadísticas con gráficos interactivos
- ✅ Análisis geográfico (países, ciudades)
- ✅ Estadísticas de sistemas operativos y arquitecturas
- ✅ Detección de VMs vs sistemas físicos
- ✅ Productos EDR/AV más comunes
- ✅ API REST para recepción de datos

## 🚀 Inicio Rápido

### Requisitos

**Agente:**
- Go 1.24 o superior
- Windows (target)

**Visualizer:**
- Python 3.8+
- Django 4.2+
- Nginx (producción)
- Gunicorn (producción)

### Instalación

#### 1. Clonar el Repositorio

```bash
git clone https://github.com/tu-usuario/artefacto.git
cd artefacto
```

#### 2. Configurar el Agente

```bash
cd artefacto

# Copiar configuración de ejemplo
cp .env.example .env

# Editar .env con tu servidor
nano .env
```

Configurar `SERVER_URL` con la URL de tu servidor visualizer:

```env
SERVER_URL=http://tu-servidor.com/api/collect
DEBUG=0
TIMEOUT=120s
```

#### 3. Compilar el Agente

```bash
# Linux/Mac
./build.sh

# Windows
go build -o agent.exe -ldflags="-s -w"
```

#### 4. Configurar el Visualizer

```bash
cd visualizer

# Instalar dependencias
pip install -r requirements.txt

# Configurar base de datos
python manage.py migrate

# Desarrollo local
python manage.py runserver 0.0.0.0:8080
```

## 📊 Uso

### Ejecutar el Agente

```bash
# Windows
.\agent.exe

# El agente recopilará datos y los enviará automáticamente al servidor
```

### Acceder al Visualizer

Abre tu navegador y ve a:

- **Página principal:** http://tu-servidor:8080
- **Estadísticas:** http://tu-servidor:8080/statistics/
- **API:** http://tu-servidor:8080/api/collect

## 🔧 Despliegue en Producción

### Opción 1: Despliegue Automático

```bash
cd visualizer
chmod +x deploy/quick_deploy.sh
sudo ./deploy/quick_deploy.sh
```

### Opción 2: Despliegue Manual

Ver documentación completa en `visualizer/README.md`

## 📁 Estructura del Proyecto

```
artefacto/
├── artefacto/              # Agente Go
│   ├── collectors/         # Módulos de recolección
│   ├── config/             # Configuración
│   ├── exfil/              # Exfiltración de datos
│   ├── models/             # Modelos de datos
│   ├── utils/              # Utilidades
│   ├── .env.example        # Ejemplo de configuración
│   ├── build.sh            # Script de compilación
│   └── main.go             # Punto de entrada
│
└── visualizer/             # Visualizer Django
    ├── collector/          # App principal
    ├── deploy/             # Scripts de despliegue
    ├── visualizer/         # Configuración Django
    ├── manage.py           # Django CLI
    └── requirements.txt    # Dependencias
```

## 🔐 Seguridad

### Configurar Autenticación

Para proteger el visualizer con autenticación HTTP Basic:

```bash
# En el servidor
sudo apt install apache2-utils
sudo mkdir -p /etc/nginx/auth
sudo htpasswd -c /etc/nginx/auth/.htpasswd tu_usuario

# Actualizar configuración de Nginx
# Ver visualizer/deploy/nginx.conf para ejemplo
```

### Configurar HTTPS

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d tu-dominio.com
```

## 📊 Datos Recopilados

El agente recopila:

- Información de sandbox (VM, indicadores)
- Sistema operativo, arquitectura, idioma
- Procesos en ejecución
- Conexiones de red
- Usuarios y grupos
- Servicios
- Variables de entorno
- Named pipes
- Aplicaciones instaladas
- Archivos recientes
- Screenshot del escritorio
- Funciones hooked
- DLLs sospechosas
- Archivos sensibles
- Productos EDR/AV
- Drivers de seguridad
- Geolocalización
- Herramientas de análisis

## 🛠️ Desarrollo

### Agente

```bash
cd artefacto

# Instalar dependencias
go mod download

# Compilar
go build -o agent.exe

# Ejecutar tests
go test ./...
```

### Visualizer

```bash
cd visualizer

# Crear entorno virtual
python -m venv venv
source venv/bin/activate  # Linux/Mac
venv\Scripts\activate     # Windows

# Instalar dependencias
pip install -r requirements.txt

# Ejecutar servidor de desarrollo
python manage.py runserver
```

## 📚 Documentación

- **Agente:** `artefacto/README.md`
- **Visualizer:** `visualizer/README.md`
- **Despliegue:** `visualizer/deploy/`

## ⚠️ Disclaimer

Esta herramienta está diseñada para uso legítimo en operaciones de Red Team y pruebas de penetración autorizadas. El uso no autorizado de esta herramienta puede ser ilegal. Los autores no se hacen responsables del mal uso de esta herramienta.

## 📝 Licencia

[Especificar licencia]

## 👥 Autores

- Marc Monfort

## 🤝 Contribuciones

Las contribuciones son bienvenidas. Por favor, abre un issue o pull request.

---

**Nota:** Esta herramienta debe usarse únicamente en entornos autorizados y con fines legítimos de seguridad.
