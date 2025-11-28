"""
Middleware personalizado para el collector
"""
from django.utils.deprecation import MiddlewareMixin

class DisableCSRFForAPIMiddleware(MiddlewareMixin):
    """
    Desactiva la verificación CSRF para el endpoint de API
    """
    def process_request(self, request):
        if request.path.startswith('/api/'):
            setattr(request, '_dont_enforce_csrf_checks', True)
