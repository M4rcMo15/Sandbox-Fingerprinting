"""Debug output XSS injection."""

import ctypes
from .base import BaseInjector
from ..payloads import XSSPayload


class DebugInjector(BaseInjector):
    """Injects XSS payloads into debug output."""
    
    def inject(self, payload: XSSPayload):
        """Send XSS payload to OutputDebugString."""
        try:
            # OutputDebugStringW
            debug_msg = f"[XSS Test {payload.id}] {payload.content}"
            ctypes.windll.kernel32.OutputDebugStringW(debug_msg)
        except Exception:
            pass
