from django.contrib import admin
from .models import XSSPayload, XSSHit, SandboxVulnerability


@admin.register(XSSPayload)
class XSSPayloadAdmin(admin.ModelAdmin):
    list_display = ('payload_id', 'vector', 'payload_type', 'status', 'execution', 'created_at')
    list_filter = ('status', 'vector', 'payload_type')
    search_fields = ('payload_id', 'execution__hostname')
    readonly_fields = ('created_at',)


@admin.register(XSSHit)
class XSSHitAdmin(admin.ModelAdmin):
    list_display = ('payload', 'source_ip', 'triggered_at')
    list_filter = ('triggered_at',)
    search_fields = ('payload__payload_id', 'source_ip', 'user_agent')
    readonly_fields = ('triggered_at',)


@admin.register(SandboxVulnerability)
class SandboxVulnerabilityAdmin(admin.ModelAdmin):
    list_display = ('sandbox_name', 'hit_count', 'first_detected', 'last_detected')
    list_filter = ('first_detected', 'last_detected')
    search_fields = ('sandbox_name', 'identified_by')
    readonly_fields = ('first_detected', 'last_detected')
