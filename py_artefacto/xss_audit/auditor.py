"""Main XSS auditor class for injecting payloads into various system vectors."""

import os
import time
from typing import List, Optional

from .payloads import XSSPayload, get_default_payloads, get_contextual_payloads
from .injectors import (
    FileInjector,
    RegistryInjector,
    ProcessInjector,
    PowerShellInjector,
    CMDInjector,
    EventLogInjector,
    EnvironmentInjector,
    NetworkInjector,
    ShortcutInjector,
    ScheduledTaskInjector,
    ServiceInjector,
    ClipboardInjector,
    DebugInjector,
    WMIInjector,
    DNSInjector,
)


class XSSAuditor:
    """Main class for XSS sandbox auditing."""
    
    def __init__(self, payloads: Optional[List[XSSPayload]] = None):
        """
        Initialize the XSS auditor.
        
        Args:
            payloads: Optional list of custom XSS payloads. If None, uses defaults.
        """
        if payloads is None:
            self.payloads = get_default_payloads()
            # Add contextual variations
            primary = self.payloads[0].content
            self.payloads.extend(get_contextual_payloads(primary))
        else:
            self.payloads = payloads
        
        # Initialize all injectors
        self.file_injector = FileInjector()
        self.registry_injector = RegistryInjector()
        self.process_injector = ProcessInjector()
        self.powershell_injector = PowerShellInjector()
        self.cmd_injector = CMDInjector()
        self.eventlog_injector = EventLogInjector()
        self.environment_injector = EnvironmentInjector()
        self.network_injector = NetworkInjector()
        self.shortcut_injector = ShortcutInjector()
        self.task_injector = ScheduledTaskInjector()
        self.service_injector = ServiceInjector()
        self.clipboard_injector = ClipboardInjector()
        self.debug_injector = DebugInjector()
        self.wmi_injector = WMIInjector()
        self.dns_injector = DNSInjector()
    
    def inject_all(self, delay: float = 0.1):
        """
        Inject XSS payloads into all available vectors.
        
        Args:
            delay: Delay in seconds between injection batches
        """
        print("[*] Starting comprehensive XSS audit...")
        
        vectors = [
            ("Files", self.inject_files),
            ("Registry", self.inject_registry),
            ("Processes", self.inject_processes),
            ("PowerShell", self.inject_powershell),
            ("CMD", self.inject_cmd),
            ("Event Logs", self.inject_eventlogs),
            ("Environment", self.inject_environment),
            ("Network", self.inject_network),
            ("Shortcuts", self.inject_shortcuts),
            ("Scheduled Tasks", self.inject_tasks),
            ("Services", self.inject_services),
            ("Clipboard", self.inject_clipboard),
            ("Debug Output", self.inject_debug),
            ("WMI", self.inject_wmi),
            ("DNS", self.inject_dns),
        ]
        
        for name, inject_func in vectors:
            print(f"[+] Injecting into {name}...")
            try:
                inject_func()
                time.sleep(delay)
            except Exception as e:
                print(f"[-] Error injecting into {name}: {e}")
        
        print("[*] XSS audit complete!")
    
    def inject_files(self):
        """Inject XSS payloads into various file types."""
        for payload in self.payloads:
            if payload.vector in ["all", "files"]:
                self.file_injector.inject(payload)
    
    def inject_registry(self):
        """Inject XSS payloads into Windows Registry."""
        for payload in self.payloads:
            if payload.vector in ["all", "registry"]:
                self.registry_injector.inject(payload)
    
    def inject_processes(self):
        """Inject XSS payloads into process arguments."""
        for payload in self.payloads:
            if payload.vector in ["all", "processes"]:
                self.process_injector.inject(payload)
    
    def inject_powershell(self):
        """Inject XSS payloads via PowerShell commands."""
        for payload in self.payloads:
            if payload.vector in ["all", "powershell"]:
                self.powershell_injector.inject(payload)
    
    def inject_cmd(self):
        """Inject XSS payloads via CMD commands."""
        for payload in self.payloads:
            if payload.vector in ["all", "cmd"]:
                self.cmd_injector.inject(payload)
    
    def inject_eventlogs(self):
        """Inject XSS payloads into Windows Event Logs."""
        for payload in self.payloads:
            if payload.vector in ["all", "eventlogs"]:
                self.eventlog_injector.inject(payload)
    
    def inject_environment(self):
        """Inject XSS payloads into environment variables."""
        for payload in self.payloads:
            if payload.vector in ["all", "environment"]:
                self.environment_injector.inject(payload)
    
    def inject_network(self):
        """Inject XSS payloads into network traffic."""
        for payload in self.payloads:
            if payload.vector in ["all", "network"]:
                self.network_injector.inject(payload)
    
    def inject_shortcuts(self):
        """Inject XSS payloads into Windows shortcuts."""
        for payload in self.payloads:
            if payload.vector in ["all", "shortcuts"]:
                self.shortcut_injector.inject(payload)
    
    def inject_tasks(self):
        """Inject XSS payloads into scheduled tasks."""
        for payload in self.payloads:
            if payload.vector in ["all", "tasks"]:
                self.task_injector.inject(payload)
    
    def inject_services(self):
        """Inject XSS payloads into service descriptions."""
        for payload in self.payloads:
            if payload.vector in ["all", "services"]:
                self.service_injector.inject(payload)
    
    def inject_clipboard(self):
        """Inject XSS payloads into clipboard."""
        for payload in self.payloads:
            if payload.vector in ["all", "clipboard"]:
                self.clipboard_injector.inject(payload)
    
    def inject_debug(self):
        """Inject XSS payloads into debug output."""
        for payload in self.payloads:
            if payload.vector in ["all", "debug"]:
                self.debug_injector.inject(payload)
    
    def inject_wmi(self):
        """Inject XSS payloads via WMI queries."""
        for payload in self.payloads:
            if payload.vector in ["all", "wmi"]:
                self.wmi_injector.inject(payload)
    
    def inject_dns(self):
        """Inject XSS payloads via DNS queries."""
        for payload in self.payloads:
            if payload.vector in ["all", "dns"]:
                self.dns_injector.inject(payload)
