"""
Middleware personalizado para el collector
"""
from django.utils.deprecation import MiddlewareMixin
from django.conf import settings
from django.shortcuts import redirect

class DisableCSRFForAPIMiddleware(MiddlewareMixin):
    """
    Desactiva la verificación CSRF para el endpoint de API
    """
    def process_request(self, request):
        if request.path.startswith('/api/'):
            setattr(request, '_dont_enforce_csrf_checks', True)


class GlobalLoginRequiredMiddleware:
    """
    Middleware para forzar login en toda la aplicación excepto en rutas públicas
    """
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        # Si el usuario ya está autenticado, permitir acceso
        if request.user.is_authenticated:
            return self.get_response(request)

        path = request.path_info

        # Rutas públicas que NO requieren login
        public_paths = [
            '/admin/',           # Django Admin (tiene su propio login)
            '/api/',             # API para que el agente envíe datos
            '/static/',          # Archivos estáticos (CSS, JS)
            '/media/',           # Archivos media
        ]

        # Verificar si la ruta actual es pública
        for p in public_paths:
            if path.startswith(p):
                return self.get_response(request)
        
        # Si no es pública y no estamos logueados, redirigir al login
        # settings.LOGIN_URL es '/admin/login/'
        if path == settings.LOGIN_URL:
            return self.get_response(request)
            
        return redirect(f"{settings.LOGIN_URL}?next={request.path}")
