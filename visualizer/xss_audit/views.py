from django.shortcuts import render, get_object_or_404
from django.http import HttpResponse, JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods
from django.db.models import Count, Q
from django.utils import timezone
from datetime import timedelta
from .models import XSSPayload, XSSHit, SandboxVulnerability
from collector.models import AgentExecution


def get_client_ip(request):
    """Obtiene la IP real del cliente"""
    x_forwarded_for = request.META.get('HTTP_X_FORWARDED_FOR')
    if x_forwarded_for:
        ip = x_forwarded_for.split(',')[0]
    else:
        ip = request.META.get('REMOTE_ADDR')
    return ip


@csrf_exempt
@require_http_methods(["GET", "POST"])
def xss_callback(request):
    """
    Endpoint que recibe los callbacks cuando un XSS es triggerado.
    Este endpoint debe ser público (sin autenticación) para que funcione.
    """
    payload_id = request.GET.get('id')
    vector = request.GET.get('v', 'unknown')
    
    if not payload_id:
        # Retornar imagen 1x1 transparente para no romper nada
        return HttpResponse(
            b'\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\xff\xff\xff\x21\xf9\x04\x01\x00\x00\x00\x00\x2c\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3b',
            content_type='image/gif'
        )
    
    try:
        # Buscar el payload
        payload = XSSPayload.objects.get(payload_id=payload_id)
        
        # Obtener información del request
        source_ip = get_client_ip(request)
        user_agent = request.META.get('HTTP_USER_AGENT', '')
        referer = request.META.get('HTTP_REFERER', '')
        
        # Crear el hit
        hit = XSSHit.objects.create(
            payload=payload,
            source_ip=source_ip,
            user_agent=user_agent,
            referer=referer,
            request_headers={
                'user_agent': user_agent,
                'referer': referer,
                'accept': request.META.get('HTTP_ACCEPT', ''),
                'accept_language': request.META.get('HTTP_ACCEPT_LANGUAGE', ''),
                'accept_encoding': request.META.get('HTTP_ACCEPT_ENCODING', ''),
            }
        )
        
        # Actualizar estado del payload
        if payload.status == 'injected':
            payload.status = 'triggered'
            payload.save()
        
        # Intentar identificar el sandbox
        identify_sandbox(hit)
        
        print(f"[XSS HIT] Payload {payload_id[:8]} triggerado desde {source_ip}")
        
    except XSSPayload.DoesNotExist:
        print(f"[XSS] Payload {payload_id} no encontrado en la base de datos")
    except Exception as e:
        print(f"[XSS ERROR] {str(e)}")
    
    # Siempre retornar imagen 1x1 transparente
    return HttpResponse(
        b'\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\xff\xff\xff\x21\xf9\x04\x01\x00\x00\x00\x00\x2c\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3b',
        content_type='image/gif'
    )


def identify_sandbox(hit):
    """
    Intenta identificar el sandbox basándose en patrones de IP, user-agent, etc.
    """
    sandbox_patterns = {
        'any.run': ['any.run', 'anyrun'],
        'hybrid-analysis': ['hybrid-analysis', 'falcon'],
        'joe-sandbox': ['joe', 'joesan'],
        'tria.ge': ['triage', 'hatching'],
        'virustotal': ['virustotal', 'vt'],
        'cuckoo': ['cuckoo'],
        'cape': ['cape'],
    }
    
    user_agent_lower = hit.user_agent.lower()
    referer_lower = hit.referer.lower()
    
    identified_sandbox = None
    
    # Buscar patrones en user-agent y referer
    for sandbox_name, patterns in sandbox_patterns.items():
        for pattern in patterns:
            if pattern in user_agent_lower or pattern in referer_lower:
                identified_sandbox = sandbox_name
                break
        if identified_sandbox:
            break
    
    # Si no se identificó, usar "Unknown"
    if not identified_sandbox:
        identified_sandbox = f"Unknown ({hit.source_ip})"
    
    # Actualizar o crear registro de sandbox vulnerable
    sandbox, created = SandboxVulnerability.objects.get_or_create(
        sandbox_name=identified_sandbox,
        defaults={
            'identified_by': f"IP: {hit.source_ip}, UA: {hit.user_agent[:100]}",
            'vulnerable_vectors': [hit.payload.vector],
            'hit_count': 1,
        }
    )
    
    if not created:
        # Actualizar sandbox existente
        sandbox.hit_count += 1
        if hit.payload.vector not in sandbox.vulnerable_vectors:
            sandbox.vulnerable_vectors.append(hit.payload.vector)
        sandbox.save()


def xss_dashboard(request):
    """Dashboard principal de XSS Audit"""
    
    # Estadísticas generales
    total_payloads = XSSPayload.objects.count()
    triggered_payloads = XSSPayload.objects.filter(status='triggered').count()
    total_hits = XSSHit.objects.count()
    vulnerable_sandboxes = SandboxVulnerability.objects.count()
    
    # Calcular tasa de éxito
    success_rate = (triggered_payloads / total_payloads * 100) if total_payloads > 0 else 0
    
    # Últimos hits (últimos 20)
    recent_hits = XSSHit.objects.select_related('payload', 'payload__execution').order_by('-triggered_at')[:20]
    
    # Vectores más exitosos
    vector_stats_raw = XSSPayload.objects.filter(status='triggered').values('vector').annotate(
        count=Count('id')
    ).order_by('-count')[:10]
    
    # Calcular el ancho de la barra para cada vector (count * 10)
    vector_stats = []
    for stat in vector_stats_raw:
        vector_stats.append({
            'vector': stat['vector'],
            'count': stat['count'],
            'width': stat['count'] * 10  # Ancho en porcentaje
        })
    
    # Sandboxes identificados
    sandboxes = SandboxVulnerability.objects.all().order_by('-hit_count')[:10]
    
    # Hits en las últimas 24 horas
    last_24h = timezone.now() - timedelta(hours=24)
    hits_24h = XSSHit.objects.filter(triggered_at__gte=last_24h).count()
    
    # Hits en la última semana
    last_week = timezone.now() - timedelta(days=7)
    hits_week = XSSHit.objects.filter(triggered_at__gte=last_week).count()
    
    # Timeline de hits (últimos 7 días)
    timeline_data = []
    for i in range(7):
        day = timezone.now() - timedelta(days=i)
        day_start = day.replace(hour=0, minute=0, second=0, microsecond=0)
        day_end = day_start + timedelta(days=1)
        count = XSSHit.objects.filter(
            triggered_at__gte=day_start,
            triggered_at__lt=day_end
        ).count()
        timeline_data.append({
            'date': day_start.strftime('%Y-%m-%d'),
            'count': count
        })
    timeline_data.reverse()
    
    context = {
        'total_payloads': total_payloads,
        'triggered_payloads': triggered_payloads,
        'total_hits': total_hits,
        'vulnerable_sandboxes': vulnerable_sandboxes,
        'success_rate': round(success_rate, 2),
        'recent_hits': recent_hits,
        'vector_stats': vector_stats,
        'sandboxes': sandboxes,
        'hits_24h': hits_24h,
        'hits_week': hits_week,
        'timeline_data': timeline_data,
    }
    
    return render(request, 'xss_audit/dashboard.html', context)


def xss_hit_detail(request, hit_id):
    """Detalle de un hit específico"""
    hit = get_object_or_404(XSSHit, id=hit_id)
    
    context = {
        'hit': hit,
    }
    
    return render(request, 'xss_audit/hit_detail.html', context)


def xss_payloads_list(request):
    """Lista de todos los payloads inyectados"""
    payloads = XSSPayload.objects.select_related('execution').order_by('-created_at')
    
    # Filtros
    status_filter = request.GET.get('status')
    if status_filter:
        payloads = payloads.filter(status=status_filter)
    
    vector_filter = request.GET.get('vector')
    if vector_filter:
        payloads = payloads.filter(vector=vector_filter)
    
    context = {
        'payloads': payloads,
        'status_filter': status_filter,
        'vector_filter': vector_filter,
    }
    
    return render(request, 'xss_audit/payloads_list.html', context)


def sandbox_detail(request, sandbox_id):
    """Detalle de un sandbox vulnerable"""
    sandbox = get_object_or_404(SandboxVulnerability, id=sandbox_id)
    
    # Obtener todos los hits relacionados con este sandbox
    hits = XSSHit.objects.filter(
        payload__in=XSSPayload.objects.filter(
            hits__in=XSSHit.objects.filter(
                source_ip__in=XSSHit.objects.filter(
                    payload__hits__payload__in=sandbox.xsshit_set.values('payload')
                ).values('source_ip')
            )
        )
    ).distinct()[:50]
    
    context = {
        'sandbox': sandbox,
        'hits': hits,
    }
    
    return render(request, 'xss_audit/sandbox_detail.html', context)
