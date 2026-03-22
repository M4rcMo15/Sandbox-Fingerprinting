"""Clipboard XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class ClipboardInjector(BaseInjector):
    """Injects XSS payloads into clipboard."""
    
    def inject(self, payload: XSSPayload):
        """Copy XSS payload to clipboard."""
        try:
            subprocess.run(
                ["powershell.exe", "-Command", f"Set-Clipboard -Value 'XSS Test: {payload.content}'"],
                capture_output=True,
                creationflags=subprocess.CREATE_NO_WINDOW
            )
        except Exception:
            pass
