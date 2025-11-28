from django.shortcuts import render, get_object_or_404
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.utils.dateparse import parse_datetime
from django.db.models import Count, Q
import json
import uuid
from collections import Counter
from .models import (
    AgentExecution, SandboxInfo, SystemInfo, ProcessInfo,
    NetworkConnection, HookInfo, HookedFunction, CrawlerInfo,
    EDRInfo, EDRProduct, GeoLocation, ToolsInfo
)


def index(request):
    """Lista todas las ejecuciones del agente"""
    executions = AgentExecution.objects.all()
    return render(request, 'collector/index.html', {'executions': executions})


def execution_detail(request, guid):
    """Muestra el detalle de una ejecución específica"""
    execution = get_object_or_404(AgentExecution, guid=guid)
    return render(request, 'collector/detail.html', {'execution': execution})


def statistics(request):
    """Vista de estadísticas con gráficos"""
    executions = AgentExecution.objects.all()
    total_executions = executions.count()
    
    # Estadísticas de geolocalización
    geo_stats = {}
    countries = Counter()
    cities = Counter()
    ips = Counter()
    
    for execution in executions:
        if hasattr(execution, 'geo_location'):
            geo = execution.geo_location
            if geo.country:
                countries[geo.country] += 1
            if geo.city:
                cities[f"{geo.city}, {geo.country}"] += 1
        if execution.public_ip:
            ips[execution.public_ip] += 1
    
    geo_stats['countries'] = dict(countries.most_common(10))
    geo_stats['cities'] = dict(cities.most_common(10))
    geo_stats['top_ips'] = dict(ips.most_common(10))
    
    # Estadísticas de sistemas operativos
    os_stats = Counter()
    arch_stats = Counter()
    
    for execution in executions:
        if hasattr(execution, 'system_info'):
            sysinfo = execution.system_info
            if sysinfo.os:
                os_stats[sysinfo.os] += 1
            if sysinfo.architecture:
                arch_stats[sysinfo.architecture] += 1
    
    # Estadísticas de sandbox/VM
    vm_count = 0
    physical_count = 0
    
    for execution in executions:
        if hasattr(execution, 'sandbox_info'):
            if execution.sandbox_info.is_vm:
                vm_count += 1
            else:
                physical_count += 1
    
    # Estadísticas de EDR/AV
    edr_products = Counter()
    edr_types = Counter()
    executions_with_edr = 0
    
    for execution in executions:
        if hasattr(execution, 'edr_info'):
            products = execution.edr_info.detected_products.all()
            if products.exists():
                executions_with_edr += 1
            for product in products:
                if product.detected:
                    edr_products[product.name] += 1
                    edr_types[product.type] += 1
    
    # Estadísticas de herramientas de análisis
    all_tools = Counter()
    
    for execution in executions:
        if hasattr(execution, 'tools_info'):
            tools = execution.tools_info
            for tool in tools.reversing_tools + tools.debugging_tools + tools.monitoring_tools + tools.analysis_tools:
                all_tools[tool] += 1
    
    # Estadísticas de idioma y zona horaria
    languages = Counter()
    timezones = Counter()
    
    for execution in executions:
        if hasattr(execution, 'system_info'):
            sysinfo = execution.system_info
            if sysinfo.language:
                languages[sysinfo.language] += 1
            if sysinfo.timezone:
                timezones[sysinfo.timezone] += 1
    
    context = {
        'total_executions': total_executions,
        'geo_stats': geo_stats,
        'os_stats': dict(os_stats.most_common()),
        'arch_stats': dict(arch_stats.most_common()),
        'vm_count': vm_count,
        'physical_count': physical_count,
        'edr_products': dict(edr_products.most_common(10)),
        'edr_types': dict(edr_types.most_common()),
        'executions_with_edr': executions_with_edr,
        'executions_without_edr': total_executions - executions_with_edr,
        'all_tools': dict(all_tools.most_common(15)),
        'languages': dict(languages.most_common(10)),
        'timezones': dict(timezones.most_common(10)),
    }
    
    return render(request, 'collector/statistics.html', context)


@csrf_exempt
def collect_data(request):
    """Endpoint para recibir datos del agente"""
    if request.method != 'POST':
        return JsonResponse({'error': 'Only POST allowed'}, status=405)
    
    try:
        data = json.loads(request.body)
        
        # Crear ejecución con GUID único
        execution = AgentExecution.objects.create(
            timestamp=parse_datetime(data.get('timestamp')),
            hostname=data.get('hostname', 'unknown'),
            public_ip=data.get('public_ip', ''),
            binary_size_bytes=data.get('binary_size_bytes', 0)
        )
        
        # Guardar SandboxInfo
        if 'sandbox_info' in data and data['sandbox_info']:
            si = data['sandbox_info']
            SandboxInfo.objects.create(
                execution=execution,
                is_vm=si.get('is_vm', False),
                vm_indicators=si.get('vm_indicators', []),
                registry_indicators=si.get('registry_indicators', []),
                disk_indicators=si.get('disk_indicators', []),
                cpu_temperature=si.get('cpu_temperature'),
                window_count=si.get('window_count', 0),
                has_debug_privilege=si.get('has_debug_privilege', False)
            )
        
        # Guardar SystemInfo
        if 'system_info' in data and data['system_info']:
            sysinfo_data = data['system_info']
            mouse_pos = sysinfo_data.get('mouse_position', {})
            
            sysinfo = SystemInfo.objects.create(
                execution=execution,
                os=sysinfo_data.get('os', ''),
                architecture=sysinfo_data.get('architecture', ''),
                language=sysinfo_data.get('language', ''),
                timezone=sysinfo_data.get('timezone', ''),
                cpu_count=sysinfo_data.get('cpu_count', 0),
                total_ram_mb=sysinfo_data.get('total_ram_mb', 0),
                total_disk_bytes=sysinfo_data.get('total_disk_bytes', 0),
                bios=sysinfo_data.get('bios', ''),
                users=sysinfo_data.get('users', []),
                groups=sysinfo_data.get('groups', []),
                services=sysinfo_data.get('services', []),
                environment_variables=sysinfo_data.get('environment_variables', {}),
                pipes=sysinfo_data.get('pipes', []),
                screenshot_base64=sysinfo_data.get('screenshot_base64', ''),
                mouse_position_x=mouse_pos.get('x', 0),
                mouse_position_y=mouse_pos.get('y', 0),
                installed_apps=sysinfo_data.get('installed_apps', []),
                recent_files=sysinfo_data.get('recent_files', []),
                uptime_seconds=sysinfo_data.get('uptime_seconds', 0)
            )
            
            # Guardar procesos
            for proc in sysinfo_data.get('processes', []):
                ProcessInfo.objects.create(
                    system_info=sysinfo,
                    pid=proc.get('pid', 0),
                    name=proc.get('name', ''),
                    owner=proc.get('owner', ''),
                    path=proc.get('path', '')
                )
            
            # Guardar conexiones de red
            for conn in sysinfo_data.get('network_connections', []):
                NetworkConnection.objects.create(
                    system_info=sysinfo,
                    protocol=conn.get('protocol', ''),
                    local_addr=conn.get('local_addr', ''),
                    remote_addr=conn.get('remote_addr', ''),
                    state=conn.get('state', '')
                )
        
        # Guardar HookInfo
        if 'hook_info' in data and data['hook_info']:
            hi = data['hook_info']
            hook_info = HookInfo.objects.create(
                execution=execution,
                suspicious_dlls=hi.get('suspicious_dlls', [])
            )
            
            for hf in hi.get('hooked_functions', []):
                HookedFunction.objects.create(
                    hook_info=hook_info,
                    module=hf.get('module', ''),
                    function=hf.get('function', ''),
                    is_hooked=hf.get('is_hooked', False),
                    first_bytes=hf.get('first_bytes', '')
                )
        
        # Guardar CrawlerInfo
        if 'crawler_info' in data and data['crawler_info']:
            ci = data['crawler_info']
            CrawlerInfo.objects.create(
                execution=execution,
                scanned_paths=ci.get('scanned_paths', []),
                found_files=ci.get('found_files', []),
                total_files=ci.get('total_files', 0)
            )
        
        # Guardar EDRInfo
        if 'edr_info' in data and data['edr_info']:
            ei = data['edr_info']
            edr_info = EDRInfo.objects.create(
                execution=execution,
                running_processes=ei.get('running_processes', []),
                installed_drivers=ei.get('installed_drivers', [])
            )
            
            for prod in ei.get('detected_products', []):
                EDRProduct.objects.create(
                    edr_info=edr_info,
                    name=prod.get('name', ''),
                    type=prod.get('type', ''),
                    detected=prod.get('detected', False),
                    method=prod.get('method', '')
                )
        
        # Guardar GeoLocation
        if 'geo_location' in data and data['geo_location']:
            geo = data['geo_location']
            from .models import GeoLocation
            GeoLocation.objects.create(
                execution=execution,
                country=geo.get('country', ''),
                country_code=geo.get('country_code', ''),
                region=geo.get('region', ''),
                city=geo.get('city', ''),
                latitude=geo.get('latitude', 0.0),
                longitude=geo.get('longitude', 0.0),
                isp=geo.get('isp', ''),
                organization=geo.get('organization', '')
            )
        
        # Guardar ToolsInfo
        if 'tools_info' in data and data['tools_info']:
            tools = data['tools_info']
            from .models import ToolsInfo
            ToolsInfo.objects.create(
                execution=execution,
                reversing_tools=tools.get('reversing_tools', []),
                debugging_tools=tools.get('debugging_tools', []),
                monitoring_tools=tools.get('monitoring_tools', []),
                virtualization_tools=tools.get('virtualization_tools', []),
                analysis_tools=tools.get('analysis_tools', [])
            )
        
        return JsonResponse({
            'status': 'success',
            'guid': str(execution.guid),
            'message': 'Data received successfully'
        }, status=201)
        
    except Exception as e:
        return JsonResponse({
            'status': 'error',
            'message': str(e)
        }, status=400)
