from django.db import models
from collector.models import AgentExecution


class XSSPayload(models.Model):
    """Payload XSS inyectado en una ejecución"""
    payload_id = models.CharField(max_length=100, primary_key=True)
    execution = models.ForeignKey(AgentExecution, on_delete=models.CASCADE, related_name='xss_payloads')
    payload_type = models.CharField(max_length=50)  # 'img-onerror', 'script-direct', etc
    vector = models.CharField(max_length=100)  # 'hostname', 'filename', 'process', etc
    status = models.CharField(max_length=20, default='injected')  # 'injected', 'triggered'
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        ordering = ['-created_at']
    
    def __str__(self):
        return f"{self.payload_id} - {self.vector} ({self.status})"


class XSSHit(models.Model):
    """Registro de un XSS que fue triggerado"""
    payload = models.ForeignKey(XSSPayload, on_delete=models.CASCADE, related_name='hits')
    triggered_at = models.DateTimeField(auto_now_add=True)
    source_ip = models.CharField(max_length=50)
    user_agent = models.TextField(blank=True)
    referer = models.TextField(blank=True)
    
    # Datos adicionales del request
    request_headers = models.JSONField(default=dict)
    
    class Meta:
        ordering = ['-triggered_at']
    
    def __str__(self):
        return f"Hit on {self.payload.payload_id} from {self.source_ip}"


class SandboxVulnerability(models.Model):
    """Sandbox identificado como vulnerable a XSS"""
    sandbox_name = models.CharField(max_length=255)
    identified_by = models.CharField(max_length=255)  # IP pattern, user-agent, etc
    vulnerable_vectors = models.JSONField(default=list)
    first_detected = models.DateTimeField(auto_now_add=True)
    last_detected = models.DateTimeField(auto_now=True)
    hit_count = models.IntegerField(default=0)
    notes = models.TextField(blank=True)
    
    class Meta:
        ordering = ['-hit_count']
    
    def __str__(self):
        return f"{self.sandbox_name} ({self.hit_count} hits)"
