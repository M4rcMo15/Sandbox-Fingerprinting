"""Environment variable XSS injection."""

import os
from .base import BaseInjector
from ..payloads import XSSPayload


class EnvironmentInjector(BaseInjector):
    """Injects XSS payloads into environment variables."""
    
    def inject(self, payload: XSSPayload):
        """Set environment variables with XSS payloads."""
        env_vars = {
            f"XSS_TEST_{payload.id}": payload.content,
            f"XSS_PAYLOAD_{payload.id}": payload.content,
            "XSS_DATA": payload.content,
            "XSS_CONFIG": payload.content,
            "XSS_REPORT": payload.content,
            "MALWARE_XSS": payload.content,
            "ANALYSIS_XSS": payload.content,
        }
        
        for key, value in env_vars.items():
            try:
                os.environ[key] = value
            except Exception:
                pass
