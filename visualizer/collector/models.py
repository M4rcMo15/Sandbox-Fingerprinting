from django.db import models
import uuid


class AgentExecution(models.Model):
    """Representa una ejecución única del agente"""
    guid = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    timestamp = models.DateTimeField()
    hostname = models.CharField(max_length=255)
    public_ip = models.CharField(max_length=50, blank=True)
    binary_size_bytes = models.BigIntegerField(default=0)
    binary_hash = models.CharField(max_length=64, blank=True)
    target_sandbox = models.CharField(max_length=50, blank=True, null=True, default='')
    received_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        ordering = ['-received_at']
    
    def __str__(self):
        return f"{self.hostname} - {self.timestamp}"
    
    def binary_size_mb(self):
        return round(self.binary_size_bytes / (1024 * 1024), 2)


class SandboxInfo(models.Model):
    execution = models.OneToOneField(AgentExecution, on_delete=models.CASCADE, related_name='sandbox_info')
    is_vm = models.BooleanField(default=False)
    score = models.IntegerField(default=0)  # 0-100% Probability
    vm_indicators = models.JSONField(default=list)
    registry_indicators = models.JSONField(default=list)
    disk_indicators = models.JSONField(default=list)
    cpu_temperature = models.FloatField(null=True, blank=True)
    window_count = models.IntegerField(default=0)
    has_debug_privilege = models.BooleanField(default=False)
    timing_discrepancy = models.FloatField(null=True, blank=True)
    cpuid_hypervisor = models.BooleanField(default=False)


class SystemInfo(models.Model):
    execution = models.OneToOneField(AgentExecution, on_delete=models.CASCADE, related_name='system_info')
    os = models.CharField(max_length=255)
    architecture = models.CharField(max_length=50)
    language = models.CharField(max_length=100, blank=True)
    timezone = models.CharField(max_length=100, blank=True)
    cpu_count = models.IntegerField()
    total_ram_mb = models.BigIntegerField()
    total_disk_bytes = models.BigIntegerField()
    bios = models.TextField(blank=True)
    users = models.JSONField(default=list)
    groups = models.JSONField(default=list)
    services = models.JSONField(default=list)
    environment_variables = models.JSONField(default=dict)
    pipes = models.JSONField(default=list)
    screenshot_base64 = models.TextField(blank=True)
    mouse_position_x = models.IntegerField(default=0)
    mouse_position_y = models.IntegerField(default=0)
    ocr_extracted_text = models.TextField(blank=True)
    screen_resolution = models.CharField(max_length=50, blank=True)
    installed_apps = models.JSONField(default=list)
    recent_files = models.JSONField(default=list)
    uptime_seconds = models.BigIntegerField(default=0)
    mouse_history = models.JSONField(default=list)
    mac_oui = models.CharField(max_length=20, blank=True)
    clipboard_preview = models.TextField(blank=True)


class ProcessInfo(models.Model):
    system_info = models.ForeignKey(SystemInfo, on_delete=models.CASCADE, related_name='processes')
    pid = models.IntegerField()
    name = models.CharField(max_length=255)
    owner = models.CharField(max_length=255, blank=True)
    path = models.TextField(blank=True)


class NetworkConnection(models.Model):
    system_info = models.ForeignKey(SystemInfo, on_delete=models.CASCADE, related_name='network_connections')
    protocol = models.CharField(max_length=10)
    local_addr = models.CharField(max_length=100)
    remote_addr = models.CharField(max_length=100)
    state = models.CharField(max_length=50)


class HookInfo(models.Model):
    execution = models.OneToOneField(AgentExecution, on_delete=models.CASCADE, related_name='hook_info')
    suspicious_dlls = models.JSONField(default=list)


class HookedFunction(models.Model):
    hook_info = models.ForeignKey(HookInfo, on_delete=models.CASCADE, related_name='hooked_functions')
    module = models.CharField(max_length=255)
    function = models.CharField(max_length=255)
    is_hooked = models.BooleanField(default=False)
    first_bytes = models.CharField(max_length=255, blank=True)


class CrawlerInfo(models.Model):
    execution = models.OneToOneField(AgentExecution, on_delete=models.CASCADE, related_name='crawler_info')
    scanned_paths = models.JSONField(default=list)
    found_files = models.JSONField(default=list)
    total_files = models.IntegerField(default=0)


class EDRInfo(models.Model):
    execution = models.OneToOneField(AgentExecution, on_delete=models.CASCADE, related_name='edr_info')
    running_processes = models.JSONField(default=list)
    installed_drivers = models.JSONField(default=list)


class EDRProduct(models.Model):
    edr_info = models.ForeignKey(EDRInfo, on_delete=models.CASCADE, related_name='detected_products')
    name = models.CharField(max_length=255)
    type = models.CharField(max_length=50)  # EDR, AV, Sandbox
    detected = models.BooleanField(default=False)
    method = models.CharField(max_length=50)  # process, driver, registry



class GeoLocation(models.Model):
    execution = models.OneToOneField(AgentExecution, on_delete=models.CASCADE, related_name='geo_location')
    country = models.CharField(max_length=100, blank=True)
    country_code = models.CharField(max_length=10, blank=True)
    region = models.CharField(max_length=100, blank=True)
    city = models.CharField(max_length=100, blank=True)
    latitude = models.FloatField(default=0.0)
    longitude = models.FloatField(default=0.0)
    isp = models.CharField(max_length=255, blank=True)
    organization = models.CharField(max_length=255, blank=True)
    is_datacenter = models.BooleanField(default=False)


class ToolsInfo(models.Model):
    execution = models.OneToOneField(AgentExecution, on_delete=models.CASCADE, related_name='tools_info')
    reversing_tools = models.JSONField(default=list)
    debugging_tools = models.JSONField(default=list)
    monitoring_tools = models.JSONField(default=list)
    virtualization_tools = models.JSONField(default=list)
    analysis_tools = models.JSONField(default=list)
