"""
Configuración de Gunicorn para producción
"""
import multiprocessing
import os

# Dirección y puerto
bind = f"0.0.0.0:{os.getenv('PORT', '8080')}"

# Workers
workers = int(os.getenv('WORKERS', multiprocessing.cpu_count() * 2 + 1))
worker_class = 'sync'
worker_connections = 1000
timeout = 300  # 5 minutos para payloads grandes con screenshots
keepalive = 5

# Logging
accesslog = '/var/log/artefacto-visualizer/access.log'
errorlog = '/var/log/artefacto-visualizer/error.log'
loglevel = 'info'
access_log_format = '%(h)s %(l)s %(u)s %(t)s "%(r)s" %(s)s %(b)s "%(f)s" "%(a)s"'

# Process naming
proc_name = 'artefacto-visualizer'

# Server mechanics
daemon = False
pidfile = '/var/run/artefacto-visualizer/gunicorn.pid'
user = 'www-data'
group = 'www-data'
tmp_upload_dir = None

# SSL (si se necesita)
# keyfile = '/path/to/keyfile'
# certfile = '/path/to/certfile'
