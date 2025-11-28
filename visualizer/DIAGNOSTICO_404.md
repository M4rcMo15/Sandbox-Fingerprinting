# 🔧 Diagnóstico Error 404

## Problema
El agente recibe error 404 al intentar enviar datos al servidor.

## Causa
Django requiere que las URLs terminen con `/` por defecto, pero el agente está enviando a `/api/collect` sin la barra final.

## Solución Aplicada

He actualizado `collector/urls.py` para aceptar ambas versiones:
- `/api/collect` (sin barra)
- `/api/collect/` (con barra)

## Verificación

### 1. Reiniciar el Servidor Django

**IMPORTANTE:** Debes reiniciar el servidor para que los cambios surtan efecto.

```bash
# Detener el servidor actual (Ctrl+C)
# Luego reiniciar:
python manage.py runserver 192.168.1.143:8080
```

### 2. Verificar que el Endpoint Funciona

Desde otra terminal:

```bash
curl -X POST http://192.168.1.143:8080/api/collect \
  -H "Content-Type: application/json" \
  -d '{"timestamp":"2024-11-28T16:12:09Z","hostname":"test"}'
```

Deberías ver una respuesta JSON con status "success" y un GUID.

### 3. Ejecutar el Agente Nuevamente

```bash
cd artefacto
./agent.exe
```

## Verificación Alternativa

Si tienes Python con requests instalado:

```bash
cd visualizer
python test_endpoint.py
```

## URLs Configuradas

Ahora el servidor acepta:
- `POST http://192.168.1.143:8080/api/collect` ✅
- `POST http://192.168.1.143:8080/api/collect/` ✅

## Si Persiste el Error 404

1. **Verifica que el servidor esté corriendo:**
   ```bash
   curl http://192.168.1.143:8080/
   ```
   Deberías ver la página principal HTML.

2. **Verifica las URLs disponibles:**
   ```bash
   python manage.py show_urls
   ```
   O revisa manualmente `collector/urls.py`

3. **Verifica los logs del servidor Django:**
   En la terminal donde corre el servidor deberías ver:
   ```
   [28/Nov/2025 16:12:09] "POST /api/collect HTTP/1.1" 201 123
   ```

4. **Verifica que no haya errores de migración:**
   ```bash
   python manage.py migrate
   ```

## Checklist

- [ ] Servidor Django reiniciado después del cambio
- [ ] URL en artefacto/.env es correcta: `http://192.168.1.143:8080/api/collect`
- [ ] Servidor responde en http://192.168.1.143:8080/
- [ ] No hay errores en los logs del servidor
- [ ] Firewall permite conexiones en puerto 8080

## Próximos Pasos

Una vez reiniciado el servidor, ejecuta el agente nuevamente:

```bash
cd artefacto
./agent.exe
```

Deberías ver:
```
[✓] Datos enviados correctamente al servidor
```

Y en la web (http://192.168.1.143:8080/) aparecerá la nueva ejecución.
