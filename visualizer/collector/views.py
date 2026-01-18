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
from .analyzers import VMDetector, EDRDetector, ToolsDetector, GeoLocator, OCRAnalyzer
from django.core.serializers.json import DjangoJSONEncoder


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
    
    # Convertir a formato JSON para JavaScript
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
        # Versiones JSON para JavaScript
        'geo_stats_json': json.dumps(geo_stats, cls=DjangoJSONEncoder),
        'os_stats_json': json.dumps(dict(os_stats.most_common()), cls=DjangoJSONEncoder),
        'arch_stats_json': json.dumps(dict(arch_stats.most_common()), cls=DjangoJSONEncoder),
        'edr_products_json': json.dumps(dict(edr_products.most_common(10)), cls=DjangoJSONEncoder),
        'all_tools_json': json.dumps(dict(all_tools.most_common(15)), cls=DjangoJSONEncoder),
        'languages_json': json.dumps(dict(languages.most_common(10)), cls=DjangoJSONEncoder),
        'timezones_json': json.dumps(dict(timezones.most_common(10)), cls=DjangoJSONEncoder),
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
        timestamp_str = data.get('timestamp')
        timestamp = parse_datetime(timestamp_str)
        if timestamp is None:
            # Si parse_datetime falla, intentar con dateutil
            from dateutil import parser as dateutil_parser
            try:
                timestamp = dateutil_parser.parse(timestamp_str)
            except:
                # Si todo falla, usar la hora actual
                from django.utils import timezone
                timestamp = timezone.now()
        
        execution = AgentExecution.objects.create(
            timestamp=timestamp,
            hostname=data.get('hostname', 'unknown'),
            public_ip=data.get('public_ip', ''),
            binary_size_bytes=data.get('binary_size_bytes', 0),
            binary_hash=data.get('binary_hash', ''),
            target_sandbox=data.get('target_sandbox', '')
        )
        
        # === PROCESAR RAW_DATA (NUEVO) ===
        if 'raw_data' in data and data['raw_data']:
            raw_data = data['raw_data']
            system_info_dict = data.get('system_info', {})
            
            # 1. Analizar si es VM
            is_vm, score, vm_indicators = VMDetector.analyze(raw_data, system_info_dict)
            SandboxInfo.objects.create(
                execution=execution,
                is_vm=is_vm,
                score=score,
                vm_indicators=vm_indicators,
                registry_indicators=[k['path'] for k in raw_data.get('registry_keys', []) if k.get('exists')],
                disk_indicators=[raw_data.get('disk_info', {}).get('identifier', '')],
                cpu_temperature=raw_data.get('cpu_info', {}).get('temperature', 0.0),
                window_count=raw_data.get('window_count', 0),
                has_debug_privilege=False,
                timing_discrepancy=raw_data.get('timing_discrepancy', 0.0),
                cpuid_hypervisor=raw_data.get('cpuid_hypervisor_bit', False)
            )
            print(f"[Analysis] VM Detection: {is_vm} ({len(vm_indicators)} indicators)")
            
            # 2. Detectar EDR/AV
            detected_edrs = EDRDetector.analyze(raw_data)
            if detected_edrs:
                edr_info = EDRInfo.objects.create(
                    execution=execution,
                    running_processes=raw_data.get('security_processes', []),
                    installed_drivers=raw_data.get('drivers', [])
                )
                
                for edr in detected_edrs:
                    EDRProduct.objects.create(
                        edr_info=edr_info,
                        name=edr['name'],
                        type=edr['type'],
                        detected=edr['detected'],
                        method=edr['method']
                    )
                print(f"[Analysis] EDR Detection: {len(detected_edrs)} products found")
        
        # === COMPATIBILIDAD CON AGENTES ANTIGUOS ===
        # Si no viene raw_data pero viene sandbox_info directamente
        elif 'sandbox_info' in data and data['sandbox_info']:
            si = data['sandbox_info']
            SandboxInfo.objects.create(
                execution=execution,
                is_vm=si.get('is_vm', False),
                vm_indicators=si.get('vm_indicators', []),
                registry_indicators=si.get('registry_indicators', []),
                disk_indicators=si.get('disk_indicators', []),
                cpu_temperature=si.get('cpu_temperature', 0.0),
                window_count=si.get('window_count', 0),
                has_debug_privilege=si.get('has_debug_privilege', False)
            )
        
        # === PROCESAR SYSTEM_INFO ===
        if 'system_info' in data and data['system_info']:
            sysinfo_data = data['system_info']
            mouse_pos = sysinfo_data.get('mouse_position', {})
            
            # OCR y Resolución
            screenshot_b64 = sysinfo_data.get('screenshot_base64', '')
            ocr_text, resolution = OCRAnalyzer.analyze(screenshot_b64)
            
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
                ocr_extracted_text=ocr_text,
                screen_resolution=resolution,
                mouse_position_x=mouse_pos.get('x', 0),
                mouse_position_y=mouse_pos.get('y', 0),
                installed_apps=sysinfo_data.get('installed_apps', []),
                recent_files=sysinfo_data.get('recent_files', []),
                uptime_seconds=sysinfo_data.get('uptime_seconds', 0),
                mouse_history=raw_data.get('mouse_history', []),
                mac_oui=raw_data.get('mac_address_oui', ''),
                clipboard_preview=raw_data.get('clipboard_content_preview', '')
            )
            
            # Detectar herramientas de análisis
            detected_tools = ToolsDetector.analyze(sysinfo_data)
            ToolsInfo.objects.create(
                execution=execution,
                reversing_tools=detected_tools.get('reversing_tools', []),
                debugging_tools=detected_tools.get('debugging_tools', []),
                monitoring_tools=detected_tools.get('monitoring_tools', []),
                virtualization_tools=detected_tools.get('virtualization_tools', []),
                analysis_tools=detected_tools.get('analysis_tools', [])
            )
            print(f"[Analysis] Tools Detection: {sum(len(v) for v in detected_tools.values())} tools found")
            
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
        
        # === GEOLOCALIZACIÓN ===
        if execution.public_ip and not GeoLocation.objects.filter(execution=execution).exists():
            geo_data = GeoLocator.geolocate(execution.public_ip)
            if geo_data:
                GeoLocation.objects.create(
                    execution=execution,
                    country=geo_data.get('country', ''),
                    country_code=geo_data.get('country_code', ''),
                    region=geo_data.get('region', ''),
                    city=geo_data.get('city', ''),
                    latitude=geo_data.get('latitude', 0.0),
                    longitude=geo_data.get('longitude', 0.0),
                    isp=geo_data.get('isp', ''),
                    organization=geo_data.get('organization', ''),
                    is_datacenter=geo_data.get('is_datacenter', False)
                )
                print(f"[Analysis] Geolocation: {geo_data.get('city')}, {geo_data.get('country')}")
        
        # === PROCESAR HOOK_INFO ===
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
        
        # === PROCESAR CRAWLER_INFO ===
        if 'crawler_info' in data and data['crawler_info']:
            ci = data['crawler_info']
            CrawlerInfo.objects.create(
                execution=execution,
                scanned_paths=ci.get('scanned_paths', []),
                found_files=ci.get('found_files', []),
                total_files=ci.get('total_files', 0)
            )
        
        print(f"[Server] Analysis completed for execution {execution.guid}")
        
        return JsonResponse({
            'status': 'success',
            'execution_id': str(execution.guid),
            'message': 'Data processed and analyzed successfully'
        })
        
    except Exception as e:
        import traceback
        error_details = traceback.format_exc()
        print(f"[ERROR] Exception in collect_data: {str(e)}")
        print(f"[ERROR] Traceback: {error_details}")
        return JsonResponse({
            'status': 'error',
            'message': str(e),
            'details': error_details if request.GET.get('debug') else None
        }, status=400)
