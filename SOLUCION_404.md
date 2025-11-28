# 🔧 Solución al Error 404

## ❌ Problema Actual

El agente muestra:
```
[!] Error enviando datos: servidor respondió con código: 404
```

## ✅ Solución (2 pasos)

### Paso 1: Reiniciar el Servidor Django

El problema está resuelto en el código, pero **debes reiniciar el servidor** para que los cambios surtan efecto.

1. **Detén el servidor actual:**
   - Ve a la terminal donde está corriendo el servidor
   - Presiona `Ctrl + C`

2. **Reinicia el servidor:**
   ```bash
   cd visualizer
   python manage.py runserver 192.168.1.143:8080
   ```

   O usa el script:
   ```bash
   cd visualizer
   start_server.bat
   ```

### Paso 2: Ejecutar el Agente Nuevamente

```bash
cd artefacto
./agent.exe
```

Ahora deberías ver:
```
[✓] Datos enviados correctamente al servidor
```

## 🧪 Verificar que Funciona

### Opción 1: Verificación Rápida (Recomendado)

Si tienes Python con requests instalado:

```bash
cd visualizer
python check_server.py
```

Este script verificará:
- ✅ Que el servidor esté corriendo
- ✅ Que la página principal funcione
- ✅ Que el endpoint API funcione
- ✅ Que acepte POST requests

### Opción 2: Verificación Manual con curl

```bash
curl -X POST http://192.168.1.143:8080/api/collect \
  -H "Content-Type: application/json" \
  -d "{\"timestamp\":\"2024-11-28T16:12:09Z\",\"hostname\":\"test\"}"
```

Deberías ver una respuesta como:
```json
{
  "status": "success",
  "guid": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Data received successfully"
}
```

### Opción 3: Verificar en el Navegador

1. Abre http://192.168.1.143:8080/
2. Deberías ver la página principal del visualizer
3. Si hay ejecuciones previas, aparecerán listadas

## 🔍 ¿Qué se Cambió?

Actualicé el archivo `visualizer/collector/urls.py` para aceptar la URL tanto con barra final como sin ella:

```python
urlpatterns = [
    path('', views.index, name='index'),
    path('api/collect/', views.collect_data, name='collect_data'),
    path('api/collect', views.collect_data, name='collect_data_no_slash'),  # ← NUEVO
    path('execution/<uuid:guid>/', views.execution_detail, name='execution_detail'),
]
```

Ahora el servidor acepta:
- `POST /api/collect` ✅
- `POST /api/collect/` ✅

## 📊 Verificar los Datos

Una vez que el agente envíe los datos correctamente:

1. **Abre el navegador:** http://192.168.1.143:8080/
2. **Verás la lista de ejecuciones** con el nuevo GUID
3. **Haz clic en la ejecución** para ver todos los detalles

## ⚠️ Si Aún No Funciona

### 1. Verifica que el servidor esté corriendo

```bash
curl http://192.168.1.143:8080/
```

Si no responde, el servidor no está corriendo.

### 2. Verifica los logs del servidor

En la terminal donde corre Django deberías ver:
```
[28/Nov/2025 16:12:09] "POST /api/collect HTTP/1.1" 201 123
```

Si ves `404`, el endpoint no está configurado correctamente.
Si ves `500`, hay un error en el código.

### 3. Verifica la configuración del agente

Archivo `artefacto/.env` debe contener:
```env
SERVER_URL=http://192.168.1.143:8080/api/collect
```

### 4. Verifica el firewall

```bash
# Windows
netsh advfirewall firewall show rule name="Django8080"

# Si no existe, créala:
netsh advfirewall firewall add rule name="Django8080" dir=in action=allow protocol=TCP localport=8080
```

### 5. Verifica las migraciones

```bash
cd visualizer
python manage.py migrate
```

## 📝 Checklist

- [ ] Servidor Django reiniciado después del cambio
- [ ] `check_server.py` muestra todos ✅
- [ ] Página principal carga en el navegador
- [ ] Agente ejecutado nuevamente
- [ ] Mensaje "Datos enviados correctamente" aparece
- [ ] Nueva ejecución visible en la web

## 🎉 Resultado Esperado

Después de seguir estos pasos:

1. **Terminal del agente:**
   ```
   [✓] Datos enviados correctamente al servidor
   ```

2. **Terminal del servidor:**
   ```
   [28/Nov/2025 16:12:09] "POST /api/collect HTTP/1.1" 201 123
   ```

3. **Navegador (http://192.168.1.143:8080/):**
   - Nueva ejecución listada con GUID único
   - Todos los datos visibles al hacer clic

---

**¿Necesitas más ayuda?** Consulta:
- `visualizer/DIAGNOSTICO_404.md` - Diagnóstico detallado
- `visualizer/TROUBLESHOOTING.md` - Solución de otros problemas
