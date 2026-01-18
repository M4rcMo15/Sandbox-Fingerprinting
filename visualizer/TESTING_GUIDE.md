# 🧪 Testing Guide - Nueva Interfaz

## Guía de Pruebas para la Nueva Interfaz del Visualizer

---

## 🚀 Inicio Rápido

### 1. Iniciar el Servidor

```bash
cd visualizer
python manage.py runserver 0.0.0.0:8000
```

### 2. Abrir en Navegador

```
http://localhost:8000/
```

---

## ✅ Checklist de Pruebas

### Base Template

- [ ] **Navbar**
  - [ ] Logo "Artefacto" visible
  - [ ] Links: Executions, Statistics, XSS Audit
  - [ ] Link activo resaltado en azul
  - [ ] Navbar responsive (colapsa en mobile)

- [ ] **Footer**
  - [ ] Visible en todas las páginas
  - [ ] Muestra "Artefacto Visualizer v2.0"
  - [ ] Muestra estadísticas (Total Executions)

- [ ] **Estilos**
  - [ ] Fondo blanco/gris claro
  - [ ] Texto negro/gris oscuro
  - [ ] Sin tema oscuro
  - [ ] Fuentes legibles

---

### Vista de Ejecuciones (/)

- [ ] **DataTable**
  - [ ] Tabla visible con todas las columnas
  - [ ] Headers: Time, Hostname, Location, IP, Type, EDR/AV, Tools, Size, Actions
  - [ ] Datos cargados correctamente

- [ ] **Búsqueda**
  - [ ] Campo "Search" visible arriba a la derecha
  - [ ] Búsqueda funciona en tiempo real
  - [ ] Filtra por cualquier columna
  - [ ] Muestra "No matching executions found" si no hay resultados

- [ ] **Ordenación**
  - [ ] Click en header ordena ascendente
  - [ ] Segundo click ordena descendente
  - [ ] Flecha indica dirección de ordenación
  - [ ] Funciona en todas las columnas

- [ ] **Paginación**
  - [ ] Selector "Show X entries" funciona
  - [ ] Opciones: 10, 25, 50, 100, All
  - [ ] Botones Previous/Next funcionan
  - [ ] Números de página clickeables
  - [ ] Info "Showing X to Y of Z" correcta

- [ ] **Export**
  - [ ] Botón "Export CSV" visible
  - [ ] Click descarga archivo CSV
  - [ ] CSV contiene todos los datos visibles
  - [ ] Botón "Copy" copia al clipboard

- [ ] **Datos**
  - [ ] Time muestra tiempo relativo ("2m ago")
  - [ ] Hostname muestra nombre + GUID
  - [ ] Location muestra bandera + ciudad, país
  - [ ] IP muestra dirección + ISP
  - [ ] Type muestra badge (VM rojo / Physical verde)
  - [ ] EDR/AV muestra count + nombres
  - [ ] Tools muestra count
  - [ ] Size muestra MB
  - [ ] Botón "View" lleva a detalle

- [ ] **Responsive**
  - [ ] Tabla scrolleable horizontalmente en mobile
  - [ ] Todas las columnas visibles
  - [ ] Búsqueda y paginación funcionan

---

### Vista Detallada (/execution/{guid}/)

- [ ] **Breadcrumb**
  - [ ] Muestra "Executions > Hostname"
  - [ ] Link "Executions" funciona
  - [ ] Hostname no es clickeable

- [ ] **Page Header**
  - [ ] Título muestra hostname
  - [ ] Subtítulo muestra timestamp + tiempo relativo

- [ ] **Summary Cards**
  - [ ] 4 cards visibles en fila
  - [ ] Location: bandera + código país + ciudad
  - [ ] System: OS + arquitectura + CPUs
  - [ ] Detection: badge VM/Physical + indicadores
  - [ ] Security: badge EDR + tools count
  - [ ] Cards tienen hover effect

- [ ] **Tabs**
  - [ ] 5 tabs visibles: System, Detection, Network, Security, Raw Data
  - [ ] Tab "System" activo por defecto
  - [ ] Click cambia de tab
  - [ ] Contenido cambia correctamente
  - [ ] Tab activo resaltado en azul

- [ ] **Tab: System**
  - [ ] Card "System Information" visible
  - [ ] Tabla con OS, Architecture, CPUs, RAM, Disk, etc.
  - [ ] Card "Processes" con DataTable
  - [ ] DataTable de procesos funciona (búsqueda, ordenación)
  - [ ] Card "Installed Applications" con DataTable
  - [ ] DataTable de apps funciona

- [ ] **Tab: Detection**
  - [ ] Card "VM/Sandbox Detection" visible
  - [ ] Muestra is_vm, CPU temp, window count, debug privilege
  - [ ] DataTable "VM Indicators" funciona
  - [ ] Card "EDR/AV Detection" visible
  - [ ] DataTable de EDR products funciona
  - [ ] Card "Analysis Tools Detected" visible
  - [ ] Listas de tools por categoría

- [ ] **Tab: Network**
  - [ ] Card "Network Connections" visible
  - [ ] DataTable de conexiones funciona
  - [ ] Card "Geolocation" visible
  - [ ] Tabla con país, región, ciudad, coordenadas, ISP

- [ ] **Tab: Security**
  - [ ] Card "Hooked Functions" visible
  - [ ] DataTable de hooks funciona
  - [ ] Badge HOOKED en rojo, OK en verde
  - [ ] Card "Suspicious DLLs" visible (si hay)
  - [ ] Card "File Crawler Results" visible (si hay)

- [ ] **Tab: Raw Data**
  - [ ] Card "Execution Metadata" visible
  - [ ] Tabla con GUID, timestamps, hostname, IP, size

- [ ] **Botón Back**
  - [ ] Botón "← Back to Executions" visible abajo
  - [ ] Click vuelve a lista de ejecuciones

- [ ] **Responsive**
  - [ ] Summary cards apiladas en mobile (1 columna)
  - [ ] Tabs scrolleables en mobile
  - [ ] DataTables scrolleables horizontalmente

---

### Vista de Estadísticas (/statistics/)

- [ ] **KPI Cards**
  - [ ] 4 cards visibles en fila
  - [ ] Total Executions con número grande
  - [ ] Unique Countries con número
  - [ ] EDR Detected con número
  - [ ] VM Detected con número
  - [ ] Números correctos

- [ ] **Gráficos**
  - [ ] 6 gráficos visibles
  - [ ] Geographic Distribution (pie chart)
  - [ ] Operating Systems (bar chart)
  - [ ] VM vs Physical (doughnut chart)
  - [ ] EDR/AV Detection (doughnut chart)
  - [ ] EDR Products (horizontal bar) - si hay datos
  - [ ] Analysis Tools (horizontal bar) - si hay datos
  - [ ] Gráficos interactivos (hover muestra valores)
  - [ ] Leyendas visibles

- [ ] **Accordion**
  - [ ] 4 secciones colapsables
  - [ ] Top Countries
  - [ ] Top Cities
  - [ ] EDR/AV Products
  - [ ] Analysis Tools
  - [ ] Click expande/colapsa
  - [ ] DataTables dentro funcionan

- [ ] **Responsive**
  - [ ] KPI cards apiladas en mobile
  - [ ] Gráficos apilados en mobile (1 columna)
  - [ ] Gráficos en 2 columnas en tablet/desktop

---

## 🎨 Pruebas Visuales

### Colores

- [ ] **Fondos**
  - [ ] Fondo principal: blanco (#ffffff)
  - [ ] Cards: blanco con borde gris
  - [ ] Navbar: blanco con sombra sutil
  - [ ] Footer: blanco con borde superior

- [ ] **Textos**
  - [ ] Títulos: negro (#212529)
  - [ ] Texto normal: negro suave
  - [ ] Texto secundario: gris (#6c757d)
  - [ ] Texto muted: gris claro (#adb5bd)

- [ ] **Badges**
  - [ ] VM: rojo (#dc3545)
  - [ ] Physical: verde (#198754)
  - [ ] EDR: naranja (#fd7e14)
  - [ ] Tools: cyan (#0dcaf0)
  - [ ] Unknown: gris (#6c757d)

- [ ] **Botones**
  - [ ] Primary: azul (#0d6efd)
  - [ ] Secondary: gris (#6c757d)
  - [ ] Hover: color más oscuro

### Tipografía

- [ ] **Fuentes**
  - [ ] System fonts (Segoe UI, etc.)
  - [ ] Legibles en todos los tamaños
  - [ ] Sin fuentes custom

- [ ] **Tamaños**
  - [ ] Page title: grande (1.75rem)
  - [ ] Card headers: medio (1rem)
  - [ ] Body text: normal (0.95rem)
  - [ ] Small text: pequeño (0.875rem)

### Espaciado

- [ ] **Padding**
  - [ ] Cards: 1.25rem
  - [ ] Navbar: 1rem
  - [ ] Tabs: 0.75rem

- [ ] **Margin**
  - [ ] Entre secciones: 2rem
  - [ ] Entre cards: 1rem
  - [ ] Entre elementos: 0.5rem

---

## 📱 Pruebas Responsive

### Mobile (< 768px)

- [ ] **Navbar**
  - [ ] Hamburger menu visible
  - [ ] Click expande menu
  - [ ] Links apilados verticalmente

- [ ] **Summary Cards**
  - [ ] Apiladas en 1 columna
  - [ ] Ancho completo
  - [ ] Legibles

- [ ] **DataTables**
  - [ ] Scroll horizontal funciona
  - [ ] Todas las columnas accesibles
  - [ ] Búsqueda y paginación funcionan

- [ ] **Gráficos**
  - [ ] Apilados en 1 columna
  - [ ] Tamaño adecuado
  - [ ] Interactivos

- [ ] **Tabs**
  - [ ] Scrolleables horizontalmente
  - [ ] Todas las tabs accesibles

### Tablet (768px - 1024px)

- [ ] **Summary Cards**
  - [ ] 2 columnas
  - [ ] Bien distribuidas

- [ ] **Gráficos**
  - [ ] 2 columnas
  - [ ] Tamaño adecuado

- [ ] **DataTables**
  - [ ] Ancho completo
  - [ ] Sin scroll horizontal

### Desktop (> 1024px)

- [ ] **Summary Cards**
  - [ ] 4 columnas
  - [ ] Bien distribuidas

- [ ] **Gráficos**
  - [ ] 2 columnas
  - [ ] Tamaño óptimo

- [ ] **DataTables**
  - [ ] Ancho completo
  - [ ] Todas las columnas visibles

---

## ⚡ Pruebas de Performance

### Tiempos de Carga

- [ ] **Index**
  - [ ] < 200ms con 100 ejecuciones
  - [ ] < 500ms con 1000 ejecuciones
  - [ ] Sin lag al scrollear

- [ ] **Detail**
  - [ ] < 300ms con todos los datos
  - [ ] Tabs cambian instantáneamente
  - [ ] DataTables cargan rápido

- [ ] **Statistics**
  - [ ] < 500ms con todos los gráficos
  - [ ] Gráficos renderizan suavemente
  - [ ] Sin lag al interactuar

### Búsqueda

- [ ] **DataTables**
  - [ ] Búsqueda instantánea (< 100ms)
  - [ ] Sin lag con 1000+ filas
  - [ ] Resultados correctos

### Interactividad

- [ ] **Clicks**
  - [ ] Respuesta inmediata
  - [ ] Sin delay perceptible
  - [ ] Feedback visual

- [ ] **Hover**
  - [ ] Efectos suaves
  - [ ] Sin parpadeos
  - [ ] Transiciones fluidas

---

## 🐛 Pruebas de Errores

### Sin Datos

- [ ] **Index vacío**
  - [ ] Muestra "No executions received yet"
  - [ ] DataTable no crashea
  - [ ] Botones de export deshabilitados

- [ ] **Detail sin datos**
  - [ ] Muestra "No information available"
  - [ ] Tabs vacíos no crashean
  - [ ] Summary cards muestran "Unknown"

- [ ] **Statistics sin datos**
  - [ ] KPI cards muestran 0
  - [ ] Gráficos vacíos no crashean
  - [ ] Accordion vacío no crashea

### Datos Incompletos

- [ ] **Sin geolocalización**
  - [ ] Muestra "Unknown"
  - [ ] Sin bandera
  - [ ] No crashea

- [ ] **Sin EDR**
  - [ ] Muestra "None"
  - [ ] Badge verde "None"
  - [ ] No crashea

- [ ] **Sin tools**
  - [ ] Muestra "None detected"
  - [ ] Listas vacías
  - [ ] No crashea

---

## 🔍 Pruebas de Funcionalidad

### Navegación

- [ ] **Links**
  - [ ] Todos los links funcionan
  - [ ] No hay 404
  - [ ] Breadcrumbs correctos

- [ ] **Botones**
  - [ ] View lleva a detalle
  - [ ] Back vuelve a lista
  - [ ] Export descarga archivo

### DataTables

- [ ] **Búsqueda**
  - [ ] Busca en todas las columnas
  - [ ] Case insensitive
  - [ ] Resultados correctos

- [ ] **Ordenación**
  - [ ] Ordena correctamente
  - [ ] Mantiene búsqueda
  - [ ] Mantiene paginación

- [ ] **Paginación**
  - [ ] Cambia de página
  - [ ] Mantiene búsqueda
  - [ ] Mantiene ordenación

- [ ] **Export**
  - [ ] CSV correcto
  - [ ] Incluye datos filtrados
  - [ ] Formato correcto

### Gráficos

- [ ] **Interactividad**
  - [ ] Hover muestra valores
  - [ ] Click en leyenda oculta/muestra
  - [ ] Responsive

- [ ] **Datos**
  - [ ] Valores correctos
  - [ ] Colores correctos
  - [ ] Leyendas correctas

---

## ✅ Checklist Final

### Funcionalidad
- [ ] Todas las páginas cargan
- [ ] Todos los links funcionan
- [ ] Todas las búsquedas funcionan
- [ ] Todas las ordenaciones funcionan
- [ ] Todos los exports funcionan
- [ ] Todos los gráficos renderizan

### Diseño
- [ ] Colores correctos (modo claro)
- [ ] Tipografía legible
- [ ] Espaciado uniforme
- [ ] Sin iconos excesivos
- [ ] Badges de colores correctos

### Responsive
- [ ] Mobile funciona
- [ ] Tablet funciona
- [ ] Desktop funciona
- [ ] Todos los breakpoints correctos

### Performance
- [ ] Carga rápida
- [ ] Búsqueda instantánea
- [ ] Sin lag
- [ ] Sin memory leaks

### Compatibilidad
- [ ] Chrome funciona
- [ ] Firefox funciona
- [ ] Safari funciona
- [ ] Edge funciona

---

## 🎉 Resultado Esperado

Si todas las pruebas pasan:

✅ **La nueva interfaz está lista para producción**

Si alguna prueba falla:

❌ **Revisar el template correspondiente**
❌ **Verificar que los datos existan**
❌ **Comprobar la consola del navegador**

---

## 📞 Soporte

Si encuentras algún problema:

1. Verifica la consola del navegador (F12)
2. Revisa los logs del servidor Django
3. Comprueba que los datos existan en la BD
4. Verifica que las librerías CDN carguen

---

**Fecha:** 2024-01-09  
**Versión:** 2.0  
**Estado:** ✅ Lista para testing
