# 🚀 Inicio Rápido - 3 Pasos

## Paso 1: Instalar el Visualizer (2 minutos)

```bash
cd visualizer
pip install -r requirements.txt
python manage.py makemigrations
python manage.py migrate
```

## Paso 2: Iniciar el Servidor (1 comando)

### Windows
```cmd
cd visualizer
start_server.bat
```

### Linux/Mac
```bash
cd visualizer
chmod +x start_server.sh
./start_server.sh
```

**✅ Servidor corriendo en: http://192.168.1.143:8080/**

## Paso 3: Ejecutar el Agente

```bash
cd artefacto
./agent.exe
```

**✅ Los datos aparecerán automáticamente en la web!**

---

## 🎯 Acceder a la Web

Abre tu navegador y ve a:

**http://192.168.1.143:8080/**

Verás:
- Lista de todas las ejecuciones del agente
- Cada ejecución con su GUID único
- Haz clic en cualquier ejecución para ver todos los detalles

---

## 🧪 Probar sin el Agente

Si quieres probar el sistema sin ejecutar el agente:

```bash
cd visualizer
python test_api.py
```

Esto enviará datos de prueba al servidor.

---

## 📚 Más Información

- **Guía completa**: [SETUP_VISUALIZER.md](SETUP_VISUALIZER.md)
- **Instrucciones detalladas**: [visualizer/INSTRUCCIONES.md](visualizer/INSTRUCCIONES.md)
- **Solución de problemas**: [visualizer/TROUBLESHOOTING.md](visualizer/TROUBLESHOOTING.md)
- **Arquitectura**: [visualizer/ARQUITECTURA.md](visualizer/ARQUITECTURA.md)

---

## ⚠️ Problemas Comunes

### "No module named 'django'"
```bash
pip install -r requirements.txt
```

### "Port already in use"
```bash
# Cambiar el puerto
python manage.py runserver 192.168.1.143:8081
```

### "Connection refused" desde el agente
1. Verifica que el servidor esté corriendo
2. Verifica el firewall
3. Verifica la IP en `artefacto/.env`

---

## 🎉 ¡Listo!

Tu sistema de visualización está funcionando. Cada vez que ejecutes el agente, los datos aparecerán automáticamente en la web con un GUID único.
