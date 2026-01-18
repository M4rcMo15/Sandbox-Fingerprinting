"""
ASGI config for visualizer project.
"""

import os

from django.core.asgi import get_asgi_application

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'visualizer.settings')

application = get_asgi_application()
