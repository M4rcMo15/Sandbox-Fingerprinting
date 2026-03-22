"""Network traffic XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class NetworkInjector(BaseInjector):
    """Injects XSS payloads into network traffic."""
    
    def inject(self, payload: XSSPayload):
        """Send HTTP requests with XSS payloads in headers."""
        urls = [
            "http://connectivity-check.microsoft.com/connect",
            "http://www.msftncsi.com/ncsi.txt",
        ]
        
        for url in urls:
            try:
                ps_cmd = f"""
$headers = @{{
    'User-Agent' = 'Mozilla/5.0 XSS-Test: {payload.content}'
    'X-XSS-Payload' = '{payload.content}'
    'X-Test-ID' = '{payload.id}'
}}
try {{ Invoke-WebRequest -Uri '{url}' -Headers $headers -TimeoutSec 2 }} catch {{}}
"""
                subprocess.Popen(
                    ["powershell.exe", "-Command", ps_cmd],
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
            except Exception:
                pass
