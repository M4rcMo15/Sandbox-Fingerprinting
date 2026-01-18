from django.contrib import admin
from .models import (
    AgentExecution, SandboxInfo, SystemInfo, ProcessInfo,
    NetworkConnection, HookInfo, HookedFunction, CrawlerInfo,
    EDRInfo, EDRProduct, GeoLocation, ToolsInfo
)

class AgentExecutionAdmin(admin.ModelAdmin):
    list_display = ('timestamp', 'hostname', 'target_sandbox', 'public_ip', 'binary_hash')
    list_filter = ('target_sandbox', 'timestamp')
    search_fields = ('hostname', 'target_sandbox', 'public_ip', 'binary_hash')

admin.site.register(AgentExecution, AgentExecutionAdmin)
admin.site.register(SandboxInfo)
admin.site.register(SystemInfo)
admin.site.register(ProcessInfo)
admin.site.register(NetworkConnection)
admin.site.register(HookInfo)
admin.site.register(HookedFunction)
admin.site.register(CrawlerInfo)
admin.site.register(EDRInfo)
admin.site.register(EDRProduct)
admin.site.register(GeoLocation)
admin.site.register(ToolsInfo)
