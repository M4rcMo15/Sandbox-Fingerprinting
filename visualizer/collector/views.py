from django.shortcuts import render, get_object_or_404
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.utils.dateparse import parse_datetime
import json
import uuid
from .models import (
    AgentExecution, SandboxInfo, SystemInfo, ProcessInfo,
    NetworkConnection, HookInfo, HookedFunction, CrawlerInfo,
    EDRInfo, EDRProduct
)


def index(request):
    """Lista todas las ejecuciones del agente"""
    executions = AgentExecution.objects.all()
    return render(request, 'collector/index.html', {'executions': executions})


def execution_detail(request, guid):
    """Muestra el detalle de una ejecución específica"""
    execution = get_object_or_404(AgentExecution, guid=guid)
    return render(request, 'collector/detail.html', {'execution': execution})


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
            hostname=data.get('hostname', 'unknown')
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
