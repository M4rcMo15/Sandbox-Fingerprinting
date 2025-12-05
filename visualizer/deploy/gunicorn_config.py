"""
Configuración de Gunicorn para producción
"""
import multiprocessing

# Dirección y puerto (solo localhost, Nginx hace de proxy)
bind = "127.0.0.1:8080"

# Workers (ajustar según CPU)
workers = multiprocessing.cpu_count() * 2 + 1
worker_class = 'sync'
worker_connections = 1000
timeout = 300  # 5 minutos para payloads grandes
keepalive = 2

# Logging
errorlog = "/var/log/gunicorn/error.log"
accesslog = "/var/log/gunicorn/access.log"
loglevel = "info"

# Process naming
proc_name = 'artefacto-visualizer'
