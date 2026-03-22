"""CMD XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class CMDInjector(BaseInjector):
    """Injects XSS payloads via CMD commands."""
    
    def inject(self, payload: XSSPayload):
        """Execute CMD commands with XSS payloads."""
        commands = [
            f"echo {payload.content}",
            f"title {payload.id}_{payload.content[:50]}",
            f"set XSS_{payload.id}={payload.content}",
        ]
        
        for cmd in commands:
            try:
                subprocess.Popen(
                    ["cmd.exe", "/c", cmd],
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
            except Exception:
                pass
