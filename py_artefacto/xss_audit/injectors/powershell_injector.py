"""PowerShell XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class PowerShellInjector(BaseInjector):
    """Injects XSS payloads via PowerShell commands."""
    
    def inject(self, payload: XSSPayload):
        """Execute PowerShell commands with XSS payloads."""
        commands = [
            f"Write-Host 'XSS Test: {payload.content}'",
            f"Write-Output '{payload.content}'",
            f"Write-Error -Message '{payload.content}' -Category NotSpecified",
            f"New-Item -Path $env:TEMP\\xss_{payload.id}.txt -ItemType File -Value '{payload.content}' -Force",
            f"Get-Process | Select-Object -First 1; Write-Host '{payload.content}'",
        ]
        
        for cmd in commands:
            try:
                subprocess.Popen(
                    ["powershell.exe", "-Command", cmd],
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
            except Exception:
                pass
