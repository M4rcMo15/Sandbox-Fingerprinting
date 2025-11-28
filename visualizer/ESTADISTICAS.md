# 📊 Página de Estadísticas - Artefacto Visualizer

## 🎯 Descripción

La página de estadísticas proporciona una vista analítica completa de todas las ejecuciones del agente, con gráficos interactivos y métricas clave para el equipo de red team.

## 🌐 Acceso

**URL:** http://192.168.1.143:8080/statistics/

Desde la interfaz web, usa el menú de navegación en la cabecera:
- **📋 Ejecuciones** - Lista de ejecuciones individuales
- **📊 Estadísticas** - Vista analítica con gráficos

## 📈 Métricas Disponibles

### 1. Tarjetas de Resumen
- **Total Ejecuciones** - Número total de ejecuciones registradas
- **Países Únicos** - Cantidad de países desde donde se ejecutó
- **IPs Únicas** - Número de direcciones IP diferentes
- **Con EDR/AV** - Ejecuciones con productos de seguridad detectados

### 2. Distribución Geográfica

#### 🌍 Distribución por País
- **Tipo:** Gráfico de pastel (pie chart)
- **Datos:** Top 10 países con más ejecuciones
- **Utilidad:** Identificar regiones objetivo principales

#### 🏙️ Top Ciudades
- **Tipo:** Gráfico de donut
- **Datos:** Top 10 ciudades con más ejecuciones
- **Utilidad:** Análisis geográfico detallado

#### 📍 Top IPs
- **Datos:** Las 10 IPs más activas
- **Utilidad:** Identificar objetivos recurrentes

### 3. Análisis de Sistemas

#### 💻 Sistemas Operativos
- **Tipo:** Gráfico de barras
- **Datos:** Distribución de sistemas operativos
- **Utilidad:** Conocer el entorno de ejecución más común
- **Ejemplos:** Windows 10, Windows 11, Windows Server

#### 🏗️ Arquitecturas
- **Tipo:** Gráfico de pastel
- **Datos:** Distribución de arquitecturas (x86, x64, ARM)
- **Utilidad:** Optimizar compilación del agente

### 4. Detección de Entornos

#### 🖥️ Entorno de Ejecución
- **Tipo:** Gráfico de donut
- **Datos:** Máquinas virtuales vs físicas
- **Colores:**
  - 🔴 Rojo - Máquina Virtual detectada
  - 🟢 Verde - Sistema físico
- **Utilidad:** Evaluar efectividad de evasión de sandbox

### 5. Análisis de Seguridad

#### 🛡️ Detección EDR/AV
- **Tipo:** Gráfico de donut
- **Datos:** Ejecuciones con/sin productos de seguridad
- **Colores:**
  - 🔴 Rojo - Con EDR/AV detectado
  - 🟢 Verde - Sin EDR/AV
- **Utilidad:** Evaluar exposición a productos de seguridad

#### 🔒 Productos EDR/AV Detectados
- **Tipo:** Gráfico de barras horizontal
- **Datos:** Top 10 productos más encontrados
- **Ejemplos:** Windows Defender, CrowdStrike, SentinelOne
- **Utilidad:** Priorizar técnicas de evasión

#### 🔧 Herramientas de Análisis
- **Tipo:** Gráfico de barras horizontal
- **Datos:** Top 15 herramientas de reversing/debugging
- **Ejemplos:** IDA Pro, x64dbg, Wireshark, Process Monitor
- **Utilidad:** Identificar entornos de análisis

### 6. Configuración Regional

#### 🌐 Idiomas del Sistema
- **Tipo:** Gráfico de pastel
- **Datos:** Top 10 idiomas configurados
- **Utilidad:** Adaptar payloads y técnicas de ingeniería social

#### 🕐 Zonas Horarias
- **Tipo:** Gráfico de barras
- **Datos:** Top 10 zonas horarias
- **Utilidad:** Planificar timing de operaciones

## 🎨 Características de los Gráficos

### Interactividad
- ✅ Hover para ver valores exactos
- ✅ Click en leyenda para ocultar/mostrar series
- ✅ Responsive - se adaptan al tamaño de pantalla
- ✅ Colores diferenciados por categoría

### Paleta de Colores
- **Azules** (#58a6ff, #1f6feb) - Información general
- **Rojos** (#f85149, #da3633) - Alertas y detecciones
- **Verdes** (#56d364, #3fb950) - Estados positivos
- **Amarillos** (#d29922, #bb8009) - Herramientas
- **Morados** (#bc8cff, #a371f7) - Configuración

### Tecnología
- **Librería:** Chart.js 4.4.0
- **CDN:** jsdelivr.net
- **Tema:** Adaptado al diseño oscuro de la aplicación

## 💡 Casos de Uso

### Para Red Team

1. **Planificación de Operaciones**
   - Identificar países/regiones objetivo
   - Conocer sistemas operativos más comunes
   - Adaptar payloads según arquitectura

2. **Evasión de Detección**
   - Ver qué EDR/AV son más comunes
   - Identificar herramientas de análisis presentes
   - Evaluar efectividad de técnicas anti-VM

3. **Análisis de Exposición**
   - Cuántas ejecuciones fueron detectadas
   - Qué productos de seguridad están activos
   - Entornos de análisis vs producción

4. **Optimización del Agente**
   - Priorizar soporte para sistemas más comunes
   - Mejorar evasión para EDR frecuentes
   - Adaptar técnicas según región

### Para Blue Team (Defensa)

1. **Análisis de Amenazas**
   - Identificar patrones de ataque
   - Regiones de origen de amenazas
   - Sistemas más vulnerables

2. **Mejora de Detección**
   - Ver qué ejecuciones pasaron desapercibidas
   - Evaluar efectividad de EDR/AV
   - Identificar gaps de cobertura

## 🔄 Actualización de Datos

Los gráficos se generan en tiempo real cada vez que se accede a la página:
- ✅ Datos siempre actualizados
- ✅ No requiere caché
- ✅ Refleja todas las ejecuciones en BD

## 🚀 Rendimiento

- **Optimizado** para hasta 1000+ ejecuciones
- **Queries eficientes** usando Django ORM
- **Carga rápida** con Chart.js CDN
- **Sin bloqueo** del servidor

## 📱 Responsive Design

La página se adapta a diferentes tamaños de pantalla:
- **Desktop:** Gráficos en 2 columnas
- **Tablet:** Gráficos en 2 columnas ajustadas
- **Mobile:** Gráficos en 1 columna

## 🔧 Personalización

### Agregar Nuevos Gráficos

Edita `visualizer/collector/views.py` en la función `statistics()`:

```python
# Agregar nueva métrica
new_metric = Counter()
for execution in executions:
    # Tu lógica aquí
    new_metric[valor] += 1

context['new_metric'] = dict(new_metric.most_common(10))
```

Edita `visualizer/collector/templates/collector/statistics.html`:

```html
<div class="chart-container">
    <div class="chart-title">📊 Tu Nuevo Gráfico</div>
    <canvas id="newChart"></canvas>
</div>

<script>
new Chart(document.getElementById('newChart'), {
    type: 'bar',  // o 'pie', 'doughnut', 'line'
    data: {
        labels: {{ new_metric.keys|safe }},
        datasets: [{
            data: {{ new_metric.values|safe }},
            backgroundColor: '#58a6ff'
        }]
    },
    options: chartOptions
});
</script>
```

### Cambiar Colores

Modifica el array `colors` en `statistics.html`:

```javascript
const colors = [
    '#58a6ff',  // Azul claro
    '#1f6feb',  // Azul
    // ... más colores
];
```

## 🐛 Troubleshooting

### Los gráficos no se muestran
- Verifica conexión a internet (Chart.js se carga desde CDN)
- Revisa la consola del navegador (F12)
- Asegúrate de que hay datos en la base de datos

### Datos incorrectos
- Verifica que las ejecuciones tengan los campos necesarios
- Revisa la función `statistics()` en `views.py`
- Comprueba que los modelos relacionados existen

### Rendimiento lento
- Considera agregar índices en la base de datos
- Implementa paginación para muchas ejecuciones
- Usa caché de Django para datos estáticos

## 📚 Referencias

- **Chart.js Docs:** https://www.chartjs.org/docs/latest/
- **Django Aggregation:** https://docs.djangoproject.com/en/stable/topics/db/aggregation/
- **Python Counter:** https://docs.python.org/3/library/collections.html#collections.Counter

---

**Última actualización:** 28 de noviembre de 2024
