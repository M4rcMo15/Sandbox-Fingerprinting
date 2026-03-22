"""Windows Shortcut (LNK) XSS injection."""

import os
import subprocess
import tempfile
from .base import BaseInjector
from ..payloads import XSSPayload


class ShortcutInjector(BaseInjector):
    """Injects XSS payloads into Windows shortcuts."""
    
    def inject(self, payload: XSSPayload):
        """Create LNK files with XSS payloads in descriptions."""
        temp_dir = tempfile.gettempdir()
        lnk_path = os.path.join(temp_dir, f"xss_test_{payload.id}.lnk")
        
        ps_script = f"""
$shell = New-Object -COM WScript.Shell
$shortcut = $shell.CreateShortcut('{lnk_path}')
$shortcut.TargetPath = 'notepad.exe'
$shortcut.Description = 'XSS Test {payload.id}: {payload.content}'
$shortcut.Arguments = '{payload.content}'
$shortcut.Save()
"""
        
        try:
            subprocess.Popen(
                ["powershell.exe", "-Command", ps_script],
                creationflags=subprocess.CREATE_NO_WINDOW
            )
        except Exception:
            pass
