"""
Configuración de Gunicorn para Artefacto Visualizer
"""
import multiprocessing

# Server socket
bind = "127.0.0.1:8000"
backlog = 2048

# Worker processes
workers = multiprocessing.cpu_count() * 2 + 1
worker_class = "sync"
worker_connections = 1000
max_requests = 1000
max_requests_jitter = 50
timeout = 300
keepalive = 5

# Logging
accesslog = "/var/log/gunicorn/access.log"
errorlog = "/var/log/gunicorn/error.log"
loglevel = "info"
access_log_format = '%(h)s %(l)s %(u)s %(t)s "%(r)s" %(s)s %(b)s "%(f)s" "%(a)s"'

# Process naming
proc_name = "artefacto-visualizer"

# Server mechanics
daemon = False
pidfile = "/var/run/gunicorn/artefacto-visualizer.pid"
user = "www-data"
group = "www-data"
umask = 0o007

# SSL (si se usa directamente, pero recomendamos Nginx)
# keyfile = None
# certfile = None
