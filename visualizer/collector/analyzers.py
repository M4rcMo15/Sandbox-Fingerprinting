"""
Analizadores del servidor - Procesan datos en bruto del agente
"""
import requests
import base64
import io
from typing import Dict, List, Tuple

# Intentar importar librerías de imagen y OCR
try:
    from PIL import Image
    import pytesseract
except ImportError:
    Image = None
    pytesseract = None


class VMDetector:
    """Detecta si el sistema es una máquina virtual"""
    
    @staticmethod
    def analyze(raw_data: dict, system_info: dict = None) -> Tuple[bool, int, List[str]]:
        """
        Analiza los datos en bruto para determinar si es VM y calcula un Score
        
        Returns:
            (is_vm, score, indicators): Tupla con booleano, score (0-100) y lista de indicadores
        """
        indicators = []
        score = 0
        
        # Analizar archivos de VM
        vm_files = raw_data.get('vm_files', [])
        if vm_files:
            indicators.extend(vm_files)
            score += len(vm_files) * 10  # 10 puntos por archivo
        
        # Analizar claves de registro
        registry_keys = raw_data.get('registry_keys', [])
        for key in registry_keys:
            if key.get('exists', False):
                indicators.append(f"Registry: {key.get('name', 'Unknown')}")
                score += 15
        
        # Analizar identificador de disco
        disk_info = raw_data.get('disk_info', {})
        disk_id = disk_info.get('identifier', '').upper()
        vm_disk_keywords = ['VBOX', 'VMWARE', 'QEMU', 'VIRTUAL', 'DADY HARDDISK']
        for keyword in vm_disk_keywords:
            if keyword in disk_id:
                indicators.append(f"Disk: {disk_id}")
                score += 40
                break
        
        # Analizar MAC OUI
        mac_oui = raw_data.get('mac_address_oui', '').upper()
        vm_mac_ouis = ['080027', '000C29', '005056', '001C14', '0003FF', '001C42'] # VBox, VMware, VirtualPC, Parallels
        if mac_oui and mac_oui in vm_mac_ouis:
            indicators.append(f"MAC OUI: {mac_oui} (VM Vendor)")
            score += 40

        # Analizar CPU
        cpu_info = raw_data.get('cpu_info', {})
        cpu_name = cpu_info.get('processor_name', '').upper()
        if 'VIRTUAL' in cpu_name or 'QEMU' in cpu_name:
            indicators.append(f"CPU: {cpu_name}")
            score += 40
        
        # Analizar CPUID Hypervisor Bit
        if raw_data.get('cpuid_hypervisor_bit', False):
            indicators.append("CPUID Hypervisor Bit Detected")
            score += 50

        # Analizar temperatura (VMs suelen tener 0.0)
        cpu_temp = cpu_info.get('temperature', 0.0)
        if cpu_temp == 0.0:
            indicators.append("CPU Temperature: 0.0 (VM indicator)")
            score += 10
        
        # Analizar número de ventanas (VMs suelen tener pocas)
        window_count = raw_data.get('window_count', 0)
        if window_count < 10:
            indicators.append(f"Low window count: {window_count}")
            score += 10

        # Analizar Timing Discrepancy
        timing = raw_data.get('timing_discrepancy', 0.0)
        if timing > 0.5: # Si la diferencia es mayor a 0.5s
            indicators.append(f"Timing Discrepancy: {timing:.4f}s")
            score += 30
        
        # Analizar hardware limitado (común en sandboxes)
        if system_info:
            ram = system_info.get('total_ram_mb', 0)
            if 0 < ram < 4096:  # Menos de 4GB
                indicators.append(f"Low RAM: {ram} MB")
                score += 15
            
            cpu_count = system_info.get('cpu_count', 0)
            if 0 < cpu_count < 2:  # 1 Core
                indicators.append(f"Low CPU count: {cpu_count}")
                score += 15
                
            uptime = system_info.get('uptime_seconds', 0)
            if 0 < uptime < 600:  # Menos de 10 minutos encendido
                indicators.append(f"Short uptime: {uptime}s")
                score += 10
                
            # Mouse estático
            mouse = system_info.get('mouse_position', {})
            if mouse.get('x') == 0 and mouse.get('y') == 0:
                indicators.append("Mouse at (0,0)")
                score += 10
            
            # Mouse History Entropy (Simple check)
            mouse_history = raw_data.get('mouse_history', [])
            if not mouse_history and system_info.get('uptime_seconds', 0) > 60:
                indicators.append("No mouse history recorded")
                score += 10
        
        # Determinar si es VM
        # Capar el score a 100
        score = min(score, 100)
        is_vm = score >= 50  # Umbral del 50%
        
        return is_vm, score, indicators


class OCRAnalyzer:
    """Realiza OCR en capturas de pantalla"""
    
    @staticmethod
    def analyze(base64_image: str) -> Tuple[str, str]:
        """
        Extrae texto y resolución de una imagen base64
        Returns: (texto_extraido, resolucion_str)
        """
        if not base64_image or not Image:
            return "", ""
            
        try:
            # Decodificar imagen
            image_data = base64.b64decode(base64_image)
            image = Image.open(io.BytesIO(image_data))
            
            # Obtener resolución
            width, height = image.size
            resolution = f"{width}x{height}"
            
            # Intentar OCR si tesseract está instalado
            text = ""
            if pytesseract:
                try:
                    text = pytesseract.image_to_string(image)
                except Exception:
                    text = "[OCR Error: Tesseract not found or error processing]"
            else:
                text = "[OCR Error: pytesseract library not installed]"
                
            return text.strip(), resolution
            
        except Exception as e:
            print(f"[OCR Error] {e}")
            return "", ""


class EDRDetector:
    """Detecta productos EDR/AV en el sistema"""
    
    # Base de datos de productos conocidos
    EDR_PRODUCTS = {
        'Windows Defender': {
            'processes': ['msmpeng.exe', 'nissrv.exe', 'securityhealthservice.exe', 'mssense.exe'],
            'drivers': ['wdfilter.sys', 'wdnisdrv.sys', 'wdboot.sys'],
        },
        'CrowdStrike Falcon': {
            'processes': ['csfalconservice.exe', 'csfalconcontainer.exe'],
            'drivers': ['csagent.sys', 'csdevicecontrol.sys', 'csboot.sys'],
        },
        'SentinelOne': {
            'processes': ['sentinelagent.exe', 'sentinelservicehost.exe', 'sentinelstaticengine.exe'],
            'drivers': ['sentinelmonitor.sys'],
        },
        'Carbon Black': {
            'processes': ['cb.exe', 'repmgr.exe', 'reputils.exe', 'repwsc.exe'],
            'drivers': ['cbk7.sys', 'parity.sys', 'cbstream.sys'],
        },
        'Cylance': {
            'processes': ['cylancesvc.exe', 'cylanceui.exe'],
            'drivers': ['cylancedrv.sys'],
        },
        'Symantec Endpoint Protection': {
            'processes': ['ccsvchst.exe', 'smc.exe', 'smcgui.exe'],
            'drivers': ['srtsp.sys', 'symefa.sys', 'symnets.sys'],
        },
        'McAfee Endpoint Security': {
            'processes': ['mfemms.exe', 'mfevtps.exe', 'mcshield.exe'],
            'drivers': ['mfehidk.sys', 'mfefirek.sys', 'mfeavfk.sys'],
        },
        'Kaspersky': {
            'processes': ['avp.exe', 'avpui.exe', 'kavfs.exe'],
            'drivers': ['klif.sys', 'kl1.sys', 'klbackupdisk.sys'],
        },
        'Trend Micro': {
            'processes': ['tmbmsrv.exe', 'tmccsf.exe', 'tmlisten.exe'],
            'drivers': ['tmcomm.sys', 'tmactmon.sys', 'tmevtmgr.sys'],
        },
        'ESET': {
            'processes': ['ekrn.exe', 'egui.exe', 'eguiproxy.exe'],
            'drivers': ['eamonm.sys', 'ehdrv.sys', 'epfw.sys'],
        },
        'Palo Alto Traps': {
            'processes': ['cyserver.exe', 'cyvera.exe', 'cyveraw.exe'],
            'drivers': ['tlaworker.sys', 'cyvrmtgn.sys'],
        },
        'FireEye': {
            'processes': ['xagt.exe', 'xagtnotif.exe'],
            'drivers': ['xagt.sys'],
        },
        'Sophos': {
            'processes': ['sophoshealth.exe', 'sophosfs.exe', 'sophosui.exe'],
            'drivers': ['sophosed.sys', 'sophossp.sys'],
        },
        'Avast': {
            'processes': ['avastui.exe', 'avastsvc.exe', 'avastbrowser.exe'],
            'drivers': ['aswsp.sys', 'aswvmm.sys', 'aswmon.sys'],
        },
        'AVG': {
            'processes': ['avgui.exe', 'avgsvc.exe'],
            'drivers': ['avgsp.sys', 'avgrkx64.sys'],
        },
        'Bitdefender': {
            'processes': ['bdagent.exe', 'bdservicehost.exe', 'updatesrv.exe'],
            'drivers': ['bdelam.sys', 'bdvedisk.sys', 'ignis.sys'],
        },
        'Norton': {
            'processes': ['norton.exe', 'nortonlifelock.exe'],
            'drivers': ['symds.sys', 'symefa.sys', 'symnets.sys'],
        },
    }
    
    @staticmethod
    def analyze(raw_data: dict) -> List[Dict]:
        """
        Detecta productos EDR/AV basándose en procesos y drivers
        
        Returns:
            Lista de productos detectados con metadata
        """
        detected_products = []
        
        security_processes = [p.lower() for p in raw_data.get('security_processes', [])]
        drivers = [d.lower() for d in raw_data.get('drivers', [])]
        
        for product_name, signatures in EDRDetector.EDR_PRODUCTS.items():
            detected = False
            method = None
            
            # Verificar procesos
            for proc in signatures['processes']:
                if proc.lower() in security_processes:
                    detected = True
                    method = 'process'
                    break
            
            # Verificar drivers
            if not detected:
                for driver in signatures['drivers']:
                    if driver.lower() in drivers:
                        detected = True
                        method = 'driver'
                        break
            
            if detected:
                detected_products.append({
                    'name': product_name,
                    'type': 'EDR/AV',
                    'detected': True,
                    'method': method
                })
        
        return detected_products


class ToolsDetector:
    """Detecta herramientas de análisis en el sistema"""
    
    TOOLS_DATABASE = {
        'reversing': {
            'IDA Pro': ['ida.exe', 'ida64.exe', 'idaq.exe', 'idaq64.exe'],
            'Ghidra': ['ghidra.exe', 'ghidrarun.exe'],
            'Binary Ninja': ['binaryninja.exe'],
            'Radare2': ['radare2.exe', 'r2.exe'],
            'Hopper': ['hopper.exe'],
        },
        'debugging': {
            'x64dbg': ['x64dbg.exe', 'x32dbg.exe'],
            'WinDbg': ['windbg.exe', 'windbgx.exe'],
            'OllyDbg': ['ollydbg.exe'],
            'Immunity Debugger': ['immunitydebugger.exe'],
            'GDB': ['gdb.exe'],
        },
        'monitoring': {
            'Process Monitor': ['procmon.exe', 'procmon64.exe'],
            'Process Explorer': ['procexp.exe', 'procexp64.exe'],
            'Process Hacker': ['processhacker.exe'],
            'API Monitor': ['apimonitor.exe'],
            'Wireshark': ['wireshark.exe'],
            'Fiddler': ['fiddler.exe'],
            'TCPView': ['tcpview.exe'],
        },
        'virtualization': {
            'VMware': ['vmware.exe', 'vmware-vmx.exe'],
            'VirtualBox': ['virtualbox.exe', 'vboxmanage.exe'],
            'Hyper-V': ['vmms.exe', 'vmwp.exe'],
            'QEMU': ['qemu.exe', 'qemu-system-x86_64.exe'],
            'Parallels': ['prl_client_app.exe'],
        },
        'analysis': {
            'Cuckoo Sandbox': ['cuckoo.exe', 'analyzer.exe'],
            'CAPE Sandbox': ['cape.exe'],
            'Joe Sandbox': ['joe.exe'],
            'Any.Run': ['anyrun.exe'],
            'Hybrid Analysis': ['hybrid.exe'],
        }
    }
    
    @staticmethod
    def analyze(system_info: dict) -> Dict[str, List[str]]:
        """
        Detecta herramientas de análisis basándose en procesos y aplicaciones
        
        Returns:
            Diccionario con categorías y herramientas detectadas
        """
        detected_tools = {
            'reversing_tools': [],
            'debugging_tools': [],
            'monitoring_tools': [],
            'virtualization_tools': [],
            'analysis_tools': []
        }
        
        # Obtener procesos y aplicaciones
        processes = system_info.get('processes', [])
        process_names = [p.get('name', '').lower() for p in processes]
        
        installed_apps = [app.lower() for app in system_info.get('installed_apps', [])]
        
        # Buscar en cada categoría
        for category, tools in ToolsDetector.TOOLS_DATABASE.items():
            category_key = f"{category}_tools"
            
            for tool_name, signatures in tools.items():
                found = False
                
                # Buscar en procesos
                for sig in signatures:
                    if sig.lower() in process_names:
                        found = True
                        break
                
                # Buscar en aplicaciones instaladas
                if not found:
                    for app in installed_apps:
                        if tool_name.lower() in app:
                            found = True
                            break
                
                if found:
                    detected_tools[category_key].append(tool_name)
        
        return detected_tools


class GeoLocator:
    """Obtiene geolocalización basada en IP pública"""
    
    @staticmethod
    def geolocate(public_ip: str) -> Dict:
        """
        Obtiene geolocalización de una IP pública
        
        Returns:
            Diccionario con datos de geolocalización
        """
        if not public_ip:
            return {}
        
        try:
            # Usar ip-api.com (gratis, sin API key)
            url = f"http://ip-api.com/json/{public_ip}"
            response = requests.get(url, timeout=10)
            
            if response.status_code == 200:
                data = response.json()
                
                if data.get('status') == 'success':
                    return {
                        'country': data.get('country', ''),
                        'country_code': data.get('countryCode', ''),
                        'region': data.get('regionName', ''),
                        'city': data.get('city', ''),
                        'latitude': data.get('lat', 0.0),
                        'longitude': data.get('lon', 0.0),
                        'isp': data.get('isp', ''),
                        'organization': data.get('org', ''),
                        'is_datacenter': GeoLocator.is_datacenter_ip(data.get('isp', ''), data.get('org', ''))
                    }
        except Exception as e:
            print(f"[Geolocation Error] {e}")
        
        return {}

    @staticmethod
    def is_datacenter_ip(isp: str, org: str) -> bool:
        """Detecta si la IP pertenece a un proveedor de hosting/cloud"""
        keywords = [
            'AMAZON', 'AWS', 'MICROSOFT', 'AZURE', 'GOOGLE', 'CLOUD', 
            'DIGITALOCEAN', 'HETZNER', 'OVH', 'ALIBABA', 'ORACLE', 
            'DATACENTER', 'HOSTING', 'VPS', 'LINODE'
        ]
        
        check_str = (isp + " " + org).upper()
        for keyword in keywords:
            if keyword in check_str:
                return True
        return False
