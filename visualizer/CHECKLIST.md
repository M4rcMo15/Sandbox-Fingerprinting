# ✅ Checklist de Instalación y Verificación

## Pre-requisitos

- [ ] Python 3.8 o superior instalado
- [ ] pip instalado y actualizado
- [ ] Acceso a la red en 192.168.1.143
- [ ] Puerto 8080 disponible

## Instalación del Visualizer

### 1. Instalar Dependencias
```bash
cd visualizer
pip install -r requirements.txt
```
- [ ] Django instalado correctamente
- [ ] djangorestframework instalado correctamente
- [ ] Sin errores de instalación

### 2. Configurar Base de Datos
```bash
python manage.py makemigrations
python manage.py migrate
```
- [ ] Migraciones creadas sin errores
- [ ] Archivo `db.sqlite3` creado
- [ ] Todas las tablas creadas correctamente

### 3. Verificar Configuración
```bash
python manage.py check
```
- [ ] Sin errores de configuración
- [ ] Sin warnings críticos

## Iniciar el Servidor

### Windows
```cmd
start_server.bat
```

### Linux/Mac
```bash
chmod +x start_server.sh
./start_server.sh
```

- [ ] Servidor inicia sin errores
- [ ] Mensaje "Starting development server at http://192.168.1.143:8080/"
- [ ] Sin errores en la consola

## Verificar Acceso Web

### Desde el mismo equipo
- [ ] http://localhost:8080/ carga correctamente
- [ ] http://127.0.0.1:8080/ carga correctamente
- [ ] http://192.168.1.143:8080/ carga correctamente

### Desde otro equipo en la red
- [ ] http://192.168.1.143:8080/ es accesible
- [ ] Página principal se muestra correctamente
- [ ] Sin errores 404 o 500

## Verificar Firewall

### Windows
```cmd
netsh advfirewall firewall show rule name="Django8080"
```
Si no existe:
```cmd
netsh advfirewall firewall add rule name="Django8080" dir=in action=allow protocol=TCP localport=8080
```

### Linux
```bash
sudo ufw status
sudo ufw allow 8080
```

- [ ] Firewall permite conexiones en puerto 8080
- [ ] Regla creada correctamente

## Probar el API

### Opción 1: Script de prueba
```bash
python test_api.py
```
- [ ] Script ejecuta sin errores
- [ ] Respuesta 201 Created
- [ ] GUID generado correctamente
- [ ] Mensaje de éxito mostrado

### Opción 2: curl
```bash
curl -X POST http://192.168.1.143:8080/api/collect \
  -H "Content-Type: application/json" \
  -d '{"timestamp":"2024-01-01T12:00:00Z","hostname":"test"}'
```
- [ ] Respuesta JSON con status "success"
- [ ] GUID retornado
- [ ] Status code 201

## Verificar Visualización

### Página Principal
- [ ] Lista de ejecuciones se muestra
- [ ] Si hay datos, aparecen en la lista
- [ ] Badges de VM/EDR se muestran correctamente
- [ ] Links funcionan

### Página de Detalle
- [ ] Al hacer clic en una ejecución, carga el detalle
- [ ] GUID se muestra correctamente
- [ ] Todos los desplegables funcionan
- [ ] Tablas se renderizan correctamente
- [ ] Sin errores 404 o 500

## Configurar el Agente

### Archivo .env
```bash
cd artefacto
cat .env
```
Debe contener:
```env
SERVER_URL=http://192.168.1.143:8080/api/collect
DEBUG=0
TIMEOUT=30s
```

- [ ] Archivo `.env` existe
- [ ] URL apunta al visualizer correcto
- [ ] Formato correcto

## Probar Integración Completa

### 1. Ejecutar el Agente
```bash
cd artefacto
./agent.exe
```

- [ ] Agente inicia correctamente
- [ ] Colectores se ejecutan
- [ ] Mensaje "Datos enviados correctamente"
- [ ] Sin errores de conexión

### 2. Verificar Recepción
- [ ] Servidor Django muestra POST en logs
- [ ] Status code 201 en logs
- [ ] Nueva ejecución aparece en la web

### 3. Verificar Datos
- [ ] Hacer clic en la nueva ejecución
- [ ] Todos los datos se muestran
- [ ] Sandbox info presente
- [ ] System info presente
- [ ] Procesos listados
- [ ] Conexiones de red listadas
- [ ] Hook info presente
- [ ] Crawler info presente
- [ ] EDR info presente

## Verificar Base de Datos

```bash
python manage.py shell
```

```python
from collector.models import AgentExecution
print(f"Total ejecuciones: {AgentExecution.objects.count()}")
for exe in AgentExecution.objects.all():
    print(f"- {exe.guid}: {exe.hostname} ({exe.timestamp})")
```

- [ ] Shell de Django funciona
- [ ] Queries retornan datos
- [ ] Datos coinciden con la interfaz web

## Panel de Administración (Opcional)

### Crear Superusuario
```bash
python manage.py createsuperuser
```
- [ ] Superusuario creado

### Acceder al Admin
http://192.168.1.143:8080/admin/

- [ ] Login funciona
- [ ] Panel de admin carga
- [ ] Modelos visibles
- [ ] Datos editables

## Performance

- [ ] Página principal carga en < 2 segundos
- [ ] Página de detalle carga en < 3 segundos
- [ ] API responde en < 1 segundo
- [ ] Sin memory leaks visibles

## Logs

### Servidor Django
- [ ] Logs se muestran en consola
- [ ] Requests se registran
- [ ] Errores se muestran claramente

### Agente
- [ ] Output claro y legible
- [ ] Progreso visible
- [ ] Errores se reportan

## Seguridad (Desarrollo)

- [ ] CSRF desactivado solo para /api/*
- [ ] ALLOWED_HOSTS configurado
- [ ] DEBUG=True (solo desarrollo)
- [ ] SECRET_KEY presente

## Documentación

- [ ] README.md leído
- [ ] SETUP_VISUALIZER.md consultado
- [ ] INSTRUCCIONES.md revisado
- [ ] ARQUITECTURA.md entendido
- [ ] TROUBLESHOOTING.md disponible

## Limpieza y Mantenimiento

### Limpiar base de datos (si necesario)
```bash
python manage.py flush
```

### Recrear desde cero (si necesario)
```bash
rm db.sqlite3
python manage.py makemigrations
python manage.py migrate
```

### Actualizar dependencias
```bash
pip install --upgrade -r requirements.txt
```

## Checklist Final

- [ ] ✅ Visualizer instalado y funcionando
- [ ] ✅ Servidor accesible desde la red
- [ ] ✅ API endpoint funcional
- [ ] ✅ Agente configurado correctamente
- [ ] ✅ Integración completa probada
- [ ] ✅ Datos se visualizan correctamente
- [ ] ✅ Sin errores en logs
- [ ] ✅ Documentación revisada

## 🎉 ¡Sistema Listo!

Si todos los checks están marcados, el sistema está completamente funcional y listo para usar.

## Próximos Pasos

1. Ejecutar el agente en diferentes sistemas
2. Analizar los datos recopilados
3. Identificar patrones de sandbox/VM
4. Detectar EDR/AV en entornos objetivo
5. Documentar hallazgos

## Notas

- Mantén el servidor corriendo mientras ejecutas el agente
- Cada ejecución genera un GUID único
- Los datos se almacenan permanentemente
- Puedes ejecutar el agente múltiples veces
- La interfaz se actualiza automáticamente al recargar

---

**¿Algún problema?** Consulta [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
