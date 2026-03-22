"""Windows Event Log XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class EventLogInjector(BaseInjector):
    """Injects XSS payloads into Windows Event Logs."""
    
    def inject(self, payload: XSSPayload):
        """Create event log entries with XSS payloads."""
        try:
            event_id = str(1000 + hash(payload.id) % 9000)
            
            subprocess.run(
                [
                    "eventcreate",
                    "/ID", event_id,
                    "/L", "APPLICATION",
                    "/T", "INFORMATION",
                    "/SO", "XSSTest",
                    "/D", f"XSS Test {payload.id}: {payload.content}"
                ],
                capture_output=True,
                creationflags=subprocess.CREATE_NO_WINDOW
            )
        except Exception:
            pass
